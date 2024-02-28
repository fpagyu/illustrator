package svg

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/fpagyu/illustrator"

	svg "github.com/ajstarks/svgo"
)

const (
	svghead = `<svg 
	xmlns="http://www.w3.org/2000/svg" 
	xmlns:xlink="http://www.w3.org/1999/xlink" 
	version="1.1" viewBox="0 0 %d %d">
`
)

type SVG struct {
	viewBox  [4]int
	layers   []SvgNode   // 对应ai文件的层
	gradient SvgGradient // 渐变

	group        *SvgGroup  // 当前所在的group
	currentPoint [2]float64 // 当前坐标
	path         PathBuilder
	styles       StyleBuild
	gstyle       StyleBuild // set group attr
}

func (svg *SVG) setStyle(k, v string) {
	svg.styles[k] = v
	if svg.gstyle != nil {
		svg.gstyle[k] = v
	}
}

func (svg *SVG) delStyle(k string) {
	delete(svg.styles, k)
	if svg.gstyle != nil {
		delete(svg.gstyle, k)
	}
}

func (svg *SVG) setCurrentPoint(x, y float64) {
	svg.currentPoint[0] = x
	svg.currentPoint[1] = y
}

func (svg *SVG) SetHeader(header *illustrator.AIHeader) {
	svg.viewBox = header.BoundingBox
}

func (svg *SVG) BeginLayer(layer *illustrator.AILayer) {
	if svg.styles == nil {
		svg.styles = make(StyleBuild)
	}

	svg.group = &SvgGroup{
		id:     layer.Name,
		parent: svg.group,
	}

	if parent := svg.group.parent; parent != nil {
		svg.group.indent = parent.indent + 1
	}
	svg.layers = append(svg.layers, svg.group)
}

func (svg *SVG) EndLayer() {
	svg.group = svg.group.parent
	// clear(svg.styles)
}

func (svg *SVG) SetLayerName(name string) {
	if name[0] == '(' && name[len(name)-1] == ')' {
		name = name[1 : len(name)-1]
	}
	svg.group.id = name
}

func (svg *SVG) Group() {
	if svg.group == nil {
		svg.group = &SvgGroup{
			id: "group",
		}
		svg.layers = append(svg.layers, svg.group)
	}

	svg.group = &SvgGroup{
		id:     "group",
		parent: svg.group,
		indent: svg.group.indent + 1,
	}

	parent := svg.group.parent
	parent.childs = append(parent.childs, svg.group)
}

func (svg *SVG) EndGroup() {
	svg.group = svg.group.parent
	svg.gstyle = make(StyleBuild)
}

func (svg *SVG) SetGroupAttr() {
	style := svg.gstyle.nofillstroke()
	if l := len(svg.group.childs); l > 0 {
		g, _ := svg.group.childs[l-1].(*SvgGroup)
		if g != nil {
			g.SetAttr("style", style)
		}
	}
	svg.gstyle = nil
}

func (svg *SVG) CompoundPath() {
	svg.path.Reset()
	svg.path.SetAsCompound()
}

func (svg *SVG) EndCompoundPath() {
	if svg.path.IsClip() {
		svg.group.clips = append(svg.group.clips, *svg.path.compoundPath)
	} else {
		// svg.group.childs = append(svg.group.childs, svg.path.compoundPath)
		svg.group.childs = append(svg.group.childs, &SvgCompoundPath{
			SvgPath: svg.path.compoundPath,
		})
	}
	svg.path.UnsetCompound()
}

func (svg *SVG) Moveto(x, y float64) {
	if !svg.path.IsCompound() {
		svg.path.Reset()
	}

	svg.setCurrentPoint(x, y)
	x = x - float64(svg.viewBox[0])
	y = float64(svg.viewBox[3]) - y

	d := fmt.Sprintf("M%s,%s", Float(x), Float(y))
	svg.path.WriteString(d)
}

func (svg *SVG) Lineto(x, y float64) {
	// current point
	px := svg.currentPoint[0]
	py := svg.currentPoint[1]
	svg.setCurrentPoint(x, y)

	d := fmt.Sprintf("l%s,%s", Float(x-px), Float(py-y))
	svg.path.WriteString(d)
}

func (svg *SVG) Curveto(x0, y0, x1, y1, x2, y2 float64) {
	px := svg.currentPoint[0]
	py := svg.currentPoint[1]

	svg.path.WriteString(fmt.Sprintf("c%s,%s", Float(x0-px), Float(py-y0)))
	svg.path.WriteString(fmt.Sprintf(" %s,%s", Float(x1-px), Float(py-y1)))
	svg.path.WriteString(fmt.Sprintf(" %s,%s", Float(x2-px), Float(py-y2)))

	svg.setCurrentPoint(x2, y2)
}

func (svg *SVG) Curveto1(x0, y0, x1, y1 float64) {
	svg.Curveto(x0, y0, x1, y1, x1, y1)
}

func (svg *SVG) Curveto2(x1, y1, x2, y2 float64) {
	x0 := svg.currentPoint[0]
	y0 := svg.currentPoint[1]
	svg.Curveto(x0, y0, x1, y1, x2, y2)
}

func (svg *SVG) ClosePath() {
	if svg.path.Len() > 0 {
		svg.path.WriteByte('z')
	}
}

func (svg *SVG) PathRender(t illustrator.PathOp) {
	if svg.path.Len() == 0 || t == 0 {
		// no fill no stroke
		return
	}

	var path *SvgPath
	if svg.path.IsCompound() {
		path = svg.path.compoundPath
	} else {
		if svg.path.IsClip() {
			l := len(svg.group.clips)
			path = &svg.group.clips[l-1]
		} else {
			path = &SvgPath{id: "path"}
			svg.group.childs = append(svg.group.childs, path)
		}
	}

	path.pathOp |= t
	path.d = svg.path.String()
	isFill := (t & illustrator.AI_Fill) > 0
	isStroke := (t & illustrator.AI_Stroke) > 0
	if isFill && isStroke {
		path.SetStyle(svg.styles.styles())
	} else if isFill {
		path.SetStyle(svg.styles.fill())
	} else if isStroke {
		path.SetStyle(svg.styles.stroke())
	}
}

func (svg *SVG) ClipPath() {
	if svg.path.Len() == 0 || svg.group == nil {
		return
	}

	svg.path.SetAsClip()
}

func (svg *SVG) ApplyClip() {
	if svg.group == nil {
		return
	}

	if svg.path.IsCompound() {
		return
	}

	if svg.path.IsClip() {
		svg.group.clips = append(svg.group.clips, SvgPath{
			id: "clippath",
			d:  svg.path.String(),
		})
	}
}

func (svg *SVG) SetRGB(t illustrator.PathOp, rgb [3]uint8) {
	if (t & illustrator.AI_Fill) == illustrator.AI_Fill {
		// set fill
		svg.setStyle("fill", fmt.Sprintf("#%02X%02X%02X", rgb[0], rgb[1], rgb[2]))
		return
	}

	if (t & illustrator.AI_Stroke) == illustrator.AI_Stroke {
		// set stroke
		svg.setStyle("stroke", fmt.Sprintf("#%02X%02X%02X", rgb[0], rgb[1], rgb[2]))
		return
	}
}

func (svg *SVG) SetCMYK(t illustrator.PathOp, cmyk [4]float64) {
	c, m := cmyk[0], cmyk[1]
	y, k := cmyk[2], cmyk[3]
	r := (1 - c) * (1 - k)
	g := (1 - m) * (1 - k)
	b := (1 - y) * (1 - k)

	if (t & illustrator.AI_Fill) == illustrator.AI_Fill {
		// set fill
		svg.setStyle("fill", fmt.Sprintf("#%02X%02X%02X", r, g, b))
		return
	}

	if (t & illustrator.AI_Stroke) == illustrator.AI_Stroke {
		// set stroke
		svg.setStyle("stroke", fmt.Sprintf("#%02X%02X%02X", r, g, b))
		return
	}
}

func (svg *SVG) SetOpacity(opacity string) {
	if opacity == "1" {
		svg.delStyle("opacity")
	} else {
		svg.setStyle("opacity", opacity)
	}
}

func (svg *SVG) SetDash() {}
func (svg *SVG) SetFlat() {}
func (svg *SVG) SetLineCap(v string) {
	switch v {
	case "0":
		svg.delStyle("stroke-linecap")
	case "1":
		svg.setStyle("stroke-linecap", "round")
	case "2":
		svg.setStyle("stroke-linecap", "square")
	}
}

func (svg *SVG) SetLineJoin(v string) {
	switch v {
	case "0": // miter
		svg.delStyle("stroke-linejoin")
	case "1":
		svg.setStyle("stroke-linejoin", "round")
	case "2":
		svg.setStyle("stroke-linejoin", "bevel")
	}
}

func (svg *SVG) SetLineWidth(v string) {
	// svg.styles["stroke-width"] = v
	svg.setStyle("stroke-width", v)
}

func (svg *SVG) SetMiterLimit(v string) {
	// svg.styles["stroke-miterlimit"] = v
	svg.setStyle("stroke-miterlimit", v)
}

func (svg *SVG) SetGradient(g *illustrator.Gradient) {
	if g.Flag == 2 { // disable rending
		return
	}

	g.Name = "gradient" + strconv.Itoa(len(svg.gradient.Instances))
	svg.gradient.Instances = append(svg.gradient.Instances, *g)
}

func (svg *SVG) DefGradient(g *illustrator.Gradient) {
	svg.gradient.Defs = append(svg.gradient.Defs, *g)
}

func (svg *SVG) SetRaster(raster *illustrator.Raster) {
	if len(raster.RawData) == 0 {
		return
	}

	if !svg.path.IsCompound() {
		svg.path.Reset()
	}

	image := SvgImage{
		id:     "<Image>",
		parent: svg.group,
		width:  int(raster.Width),
		height: int(raster.Height),
		matrix: raster.Matrix,
		styles: svg.styles.nofillstroke(),
	}

	image.matrix[4] = image.matrix[4] - float64(svg.viewBox[0])
	image.matrix[5] = float64(svg.viewBox[3]) - image.matrix[5]

	if svg.group != nil {
		image.indent = svg.group.indent + 1
		svg.group.childs = append(svg.group.childs, &image)
	}

	imgData, err := raster.B64Data()
	if err != nil {
		log.Println("Set Raster Image error:", err)
	}

	image.SetImage(imgData)
}

func (_svg *SVG) writeNodes(canvas *Canvas, nodes []SvgNode) {
	for _, e := range nodes {
		switch node := e.(type) {
		case *SvgPath:
			node.id = canvas.nextPathId()
			canvas.Path(node.d, node.Attrs()...)
		case *SvgCompoundPath:
			node.id = canvas.nextPathId()
			canvas.Path(node.d, node.Attrs()...)
		case *SvgGroup:
			node.id = canvas.nextGroupId()
			canvas.Group(node.Attrs()...)
			_svg.writeClips(canvas, node)
			_svg.writeNodes(canvas, node.childs)
			canvas.Gend()
		case *SvgImage:
			node.id = canvas.nextImageId()
			canvas.writeImage(node)
			// canvas.Image()
		}
	}
}

func (_Svg *SVG) writeClips(canvas *Canvas, group *SvgGroup) {
	if len(group.clips) == 0 {
		return
	}

	for i := range group.clips {
		clip := &group.clips[i]
		if clip.pathOp == 0 {
			continue
		}
		clip.id = canvas.nextPathId()
		canvas.Path(clip.d, clip.Attrs()...)
	}

	clipid := canvas.nextClipId()
	canvas.ClipPath(Attr("id", clipid))
	for i := range group.clips {
		canvas.Path(group.clips[i].d)
	}
	canvas.ClipEnd()

	for _, n := range group.childs {
		n.SetAttr("clip-path", fmt.Sprintf("url(#%s)", clipid))
	}
}

func (_svg *SVG) writeLayers(canvas *Canvas) {
	for _, e := range _svg.layers {
		node, _ := e.(*SvgGroup)
		canvas.Group(Attr("id", node.id))
		_svg.writeNodes(canvas, node.childs)
		canvas.Gend()
	}
}

func (_svg *SVG) writeGradients(canvas *Canvas) {
	for i := range _svg.gradient.Instances {
		instance := &_svg.gradient.Instances[i]
		canvas.wirteGradient(instance)
	}
}

func (_svg *SVG) writeDefs(canvas *Canvas) {
	canvas.Def()
	// _svg.writeGradients(canvas)
	canvas.DefEnd()
}

func (_svg *SVG) writeTo(w io.Writer, writeOption *SvgWriteOption) error {
	canvas := &Canvas{
		SVG:         svg.New(w),
		writeOption: writeOption,
	}

	ux := _svg.viewBox[2] - _svg.viewBox[0]
	uy := _svg.viewBox[3] - _svg.viewBox[1]
	fmt.Fprintf(canvas.Writer, svghead, ux, uy)
	fmt.Fprintln(canvas.Writer, `<!-- Generated by LEMI -->`)
	_svg.writeDefs(canvas)
	_svg.writeLayers(canvas)
	canvas.End()
	return nil
}

func (svg *SVG) Save(path string, options ...func(*SvgWriteOption)) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var writeOption SvgWriteOption
	for _, opt := range options {
		opt(&writeOption)
	}

	return svg.writeTo(file, &writeOption)
}

func (svg *SVG) Nodes(depth int) []SvgNode {
	// 以层遍历的方式输出svg的节点列表
	if len(svg.layers) == 0 {
		return nil
	}
	nodes := make([]SvgNode, 0, 100)

	for _, n := range svg.layers {
		nodes = append(nodes, n)
	}

	indent := 0
	start, end := 0, len(nodes)
	for start < end {
		indent++
		if depth >= 0 && indent > depth {
			break
		}
		for i := start; i < end; i++ {
			g, ok := nodes[i].(*SvgGroup)
			if ok {
				for _, n := range g.childs {
					nodes = append(nodes, n)
				}
			}
		}
		start, end = end, len(nodes)
	}

	return nodes
}

type Canvas struct {
	*svg.SVG

	clipid  int
	pathid  int
	groupid int
	imageid int

	writeOption *SvgWriteOption
}

func (c *Canvas) nextClipId() string {
	c.clipid++
	return "clip" + strconv.Itoa(c.clipid)
}

func (c *Canvas) nextPathId() string {
	c.pathid++
	return "path" + strconv.Itoa(c.pathid)
}

func (c *Canvas) nextGroupId() string {
	c.groupid++
	return "group" + strconv.Itoa(c.groupid)
}

func (c *Canvas) nextImageId() string {
	c.imageid++
	return "img" + strconv.Itoa(c.imageid)
}

func (c *Canvas) writeGradientColors(g *illustrator.Gradient) {
	// todo
}

func (c *Canvas) wirteGradient(g *illustrator.Gradient) {
	if g.GradientType == 0 {
		// linear
		fmt.Fprintf(c.Writer, "<linearGradient id=\"%s\" gradientUnits=\"userSpaceOnUse\">",
			g.Name,
		)
		c.writeGradientColors(g)
		fmt.Fprintf(c.Writer, "</linearGradient>")
	} else {
		// radial
		fmt.Fprintf(c.Writer, "<radialGradient id=\"%s\" gradientUnits=\"userSpaceOnUse\">",
			g.Name,
		)
		c.writeGradientColors(g)
		fmt.Fprintf(c.Writer, "</radialGradient>")
	}
}

func (c *Canvas) writeImage(img *SvgImage) {
	if c.writeOption.IgnoreImage {
		return // skip to write image node
	}
	styles := "overflow:visible;" + img.styles
	transform := fmt.Sprintf("matrix(%s,%s,%s,%s,%s,%s)",
		Float(img.matrix[0]), Float(img.matrix[1]), Float(img.matrix[2]),
		Float(img.matrix[3]), Float(img.matrix[4]), Float(img.matrix[5]),
	)
	fmt.Fprintf(c.Writer, `<image id="%s" width="%d" height="%d" transform="%s" style="%s" href="%s"></image>`,
		img.id, img.width, img.height, transform, styles, img.b64Img,
	)
	fmt.Fprintln(c.Writer)
}
