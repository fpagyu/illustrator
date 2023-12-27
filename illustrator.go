package illustrator

import (
	"bufio"
	"bytes"
	"io"
	"log"
)

var (
	AI5_EndRaster   = []byte("%AI5_EndRaster")
	AI5_BeginRaster = []byte("%AI5_BeginRaster")
)

type AIReader struct {
	*bufio.Reader

	lineBuf Buffer
}

func NewAIReader(r io.Reader) (*AIReader, error) {
	reader := &AIReader{
		Reader: bufio.NewReader(r),
	}

	return reader, nil
}

func (r *AIReader) Bytes() []byte {
	return r.lineBuf.Bytes()
}

func (r *AIReader) readLine() bool {
	var n int
	r.lineBuf.Reset()
	for {
		ch, err := r.ReadByte()
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		if ch == '\r' || ch == '\n' {
			if n > 0 {
				break
			}
		} else {
			r.lineBuf.WriteByte(ch)
			n++
		}
	}

	return r.lineBuf.len > 0
}

func (r *AIReader) Draw(drawer Drawer) error {
	header := r.readHeader()
	drawer.SetHeader(header)

	BeginSetup := []byte("%%BeginSetup")
	BeginProlog := []byte("%%BeginProlog")
	BeginLayer := []byte("%AI5_BeginLayer")
	for r.readLine() {
		line := r.Bytes()
		if bytes.HasPrefix(line, BeginSetup) {
			r.readSetup(drawer)
		} else if bytes.HasPrefix(line, BeginProlog) {
			r.readProlog()
		} else if bytes.HasPrefix(line, BeginLayer) {
			r.drawLayer(drawer)
		}
	}

	return nil
}

func (r *AIReader) readHeader() *AIHeader {
	var header AIHeader
	Title := []byte("%%Title:")
	EndComment := []byte("%%EndComments")
	BoundingBox := []byte("%%BoundingBox:")
	HiResBoundingBox := []byte("%%HiResBoundingBox:")
	for r.readLine() {
		line := r.Bytes()
		if bytes.HasSuffix(line, EndComment) {
			break
		}

		if bytes.HasPrefix(line, Title) {
			header.SetHeader(bytes.TrimPrefix(line, Title))
			continue
		}

		if bytes.HasPrefix(line, BoundingBox) {
			header.SetBoundingBox(bytes.TrimPrefix(line, BoundingBox))
			continue
		}

		if bytes.HasPrefix(line, HiResBoundingBox) {
			header.SetHiResBoundingBox(bytes.TrimPrefix(line, HiResBoundingBox))
			continue
		}
	}

	return &header
}

func (r *AIReader) readProlog() *AIProlog {
	// %%BeginProlog
	// %%EndProlog
	var prolog AIProlog
	EndProlog := []byte("%%EndProlog")
	for r.readLine() {
		line := r.Bytes()

		if bytes.HasPrefix(line, EndProlog) {
			break
		}
		// todo
	}

	return &prolog
}

func (r *AIReader) readSetup(drawer Drawer) {
	// %%BeginSetup
	// %%EndSetup
	XI := []byte("XI")
	Bd := []byte(" Bd")
	EndSetup := []byte("%%EndSetup")
	for r.readLine() {
		line := r.Bytes()

		if bytes.HasPrefix(line, EndSetup) {
			break
		}

		// todo
		if bytes.HasSuffix(line, XI) {
			r.readRasterData()
		}

		if bytes.HasSuffix(line, Bd) {
			r.defGradient(drawer)
		}
	}
}

func (r *AIReader) drawLayer(d Drawer) {
	var token lineToken
	for r.readLine() {
		line := r.Bytes()

		if line[0] == '%' {
			if bytes.Equal(AI5_BeginRaster, line) {
				r.beginRaster(d)
			}
			// skip comments
			continue
		}

		token.parse(line)
		for token.len > 0 {
			op := token.Pop()
			switch op {
			case "A": // locking, 0-unlocking; 1-locking
				token.Pop()
			case "Ap": // show center point
				token.Pop()
			case "Lb":
				r.beginLayer(d, token.PopAll())
			case "LB":
				d.EndLayer()
			case "Ln":
				d.SetLayerName(token.Pop())
			case "O", "R": // fill/stroke overprint
				token.Pop()
			case "d": // setdash
				token.Pop()
			case "D": //
			case "i": // setflat
				token.Pop()
			case "j": // linejoin
				d.SetLineJoin(token.Pop())
			case "J": // linecap
				d.SetLineCap(token.Pop())
			case "w": // linewidth
				d.SetLineWidth(token.Pop())
			case "M": // setmiterlimit
				d.SetMiterLimit(token.Pop())
			case "f": // fill
				d.ClosePath()
				d.PathRender(AI_Fill)
			case "F":
				d.PathRender(AI_Fill)
			case "s": // stroke
				d.ClosePath()
				d.PathRender(AI_Stroke)
			case "S":
				d.PathRender(AI_Stroke)
			case "b": // fill and stroke
				d.ClosePath()
				d.PathRender(AI_Fill | AI_Stroke)
			case "B":
				d.PathRender(AI_Fill | AI_Stroke)
			case "h": // close path
				d.ClosePath()
				d.ClipPath()
			case "H": // close path
				d.ClipPath()
			case "W": // clip
				d.ApplyClip()
			case "n": // no fill no stroke
				d.PathRender(0)
			case "N":
				d.ClosePath()
				d.PathRender(0)
			case "u": // begin group
				d.Group()
			case "U": // end group
				d.EndGroup()
			case "q": // begin clip group
				d.Group()
			case "Q": // end  group
				d.EndGroup()
			case "*u": // begin compound path
				d.CompoundPath()
			case "*U": // end compound path
				d.EndCompoundPath()
			case "m":
				args := token.PopN(2)
				x := toFloat(args[0])
				y := toFloat(args[1])
				d.Moveto(x, y)
			case "l", "L":
				args := token.PopN(2)
				x := toFloat(args[0])
				y := toFloat(args[1])
				d.Lineto(x, y)
			case "y", "Y":
				args := toFloatSlice(token.PopN(4))
				d.Curveto1(args[0], args[1], args[2], args[3])
			case "v", "V":
				args := toFloatSlice(token.PopN(4))
				d.Curveto2(args[0], args[1], args[2], args[3])
			case "c", "C":
				args := toFloatSlice(token.PopN(6))
				d.Curveto(args[0], args[1], args[2], args[3], args[4], args[5])
			case "g": // set fill tint
				if tint := token.Pop(); tint != "0" {
					log.Println("todo: NotImplement: g")
				}
			case "G": // set stroke tint
				if tint := token.Pop(); tint != "0" {
					log.Println("todo: NotImplement: g")
				}
			case "k": // fill setcmykcolor
				args := KArgs(token.PopAll())
				d.SetCMYK(AI_Fill, args.CMYK())
			case "K": // stroke setcmykcolor
				args := KArgs(token.PopAll())
				d.SetCMYK(AI_Stroke, args.CMYK())
			case "x": // custom fill
				args := XArgs(token.PopAll())
				d.SetCMYK(AI_Fill, args.CMYK())
			case "X":
				args := XArgs(token.PopAll())
				d.SetCMYK(AI_Stroke, args.CMYK())
			case "Xy": // set opacity
				args := XYArgs(token.PopN(5))
				d.SetOpacity(args)
			case "Xa":
				args := XAArgs(token.PopAll())
				d.SetRGB(AI_Fill, args.RGB())
			case "XA":
				args := XAArgs(token.PopAll())
				d.SetRGB(AI_Stroke, args.RGB())
			case "Xk":
				args := XKArgs(token.PopAll())
				if args.colorSpace == 1 {
					d.SetRGB(AI_Fill, args.RGB())
				} else {
					d.SetCMYK(AI_Fill, args.CMYK())
				}
			case "XK":
				args := XKArgs(token.PopAll())
				if args.colorSpace == 1 {
					d.SetRGB(AI_Stroke, args.RGB())
				} else {
					d.SetCMYK(AI_Stroke, args.CMYK())
				}
			case "Xx": // custom fill color
				args := XXArgs(token.PopAll())
				if args.colorSpace == 1 {
					d.SetRGB(AI_Fill, args.RGB())
				} else {
					d.SetCMYK(AI_Fill, args.CMYK())
				}
			case "XX": // custom stroke color
				args := XXArgs(token.PopAll())
				if args.colorSpace == 1 {
					d.SetRGB(AI_Stroke, args.RGB())
				} else {
					d.SetCMYK(AI_Stroke, args.CMYK())
				}
			case "XR": // fill rule
			case "Xw": // 0--visible; 1--invisible
				// args := token.Pop()
			case "XW": // 6 () XW; 9 () XW;
				args := token.PopN(2)
				if args[0] == "6" {
					d.SetGroupAttr()
				}
			case "XG":
				args := token.PopN(2)
				if args[0] != "()" {
					log.Println("XG NotImplement:", args)
				}
			case "Bb": // begin gradient instance
				r.beginGradient(d)
			case "XI":
				r.readRasterData()
			}
		}
	}
}

func (r *AIReader) beginLayer(d Drawer, args []string) {
	layer := AILayer{
		Visible:           args[0] == "1",
		Preview:           args[1] == "1",
		Enabled:           args[2] == "1",
		Printing:          args[3] == "1",
		Dimmed:            args[4] == "1",
		HasMultiLayerMask: args[5] == "1",
		ColorIndex:        toInt8(args[6]),
	}

	colorIndex := 7
	if len(args) > 10 {
		colorIndex = 8
		layer.Name = "Layer " + args[7]
		layer.LayerIndex = toInt(args[7])
	}
	for i := 0; i < 3; i++ {
		layer.RGB[i] = toUint8(args[colorIndex+i])
	}

	d.BeginLayer(&layer)
}

func (r *AIReader) readRasterData() []byte {
	buf := bytes.NewBuffer(nil)

	var enddata1 = []byte("%%EndData")
	var enddata2 = []byte("%_%%EndData")
	for {
		ch, err := r.ReadByte()
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		buf.WriteByte(ch)

		if b := buf.Bytes(); bytes.HasSuffix(b, enddata1) {
			if bytes.HasSuffix(b, enddata2) {
				return b[:len(b)-len(enddata2)]
			} else {
				return b[:len(b)-len(enddata1)]
			}
		}
	}

	return nil
}

func (r *AIReader) defGradient(d Drawer) {
	var token lineToken
	token.parse(r.Bytes())
	args := token.PopAll()
	if !(len(args) == 4 && args[3] == "Bd") {
		log.Println("invalid Bd arguments:", args)
		return
	}
	gradient := &Gradient{
		Name:         args[0],
		GradientType: toInt8(args[1]),
		nColors:      toInt(args[2]),
	}

	for r.readLine() {
		line := r.Bytes()

		if line[0] == '%' {
			// skip comment
			continue
		}

		token.parse(line)
		for token.len > 0 {
			op := token.Pop()
			switch op {
			case "BD":
				d.DefGradient(gradient)
				return
			case "%_Bs", "%_BS":
				if args := BSArgs(token.PopAll()); args != nil {
					gradient.AddColor(args)
				}
			}
		}
	}
}

func (r *AIReader) beginGradient(d Drawer) {
	var gradient Gradient

	var token lineToken
	for r.readLine() {
		line := r.Bytes()

		if line[0] == '%' {
			// skip comment
			continue
		}

		token.parse(line)
		for token.len > 0 {
			op := token.Pop()
			switch op {
			case "f":
				// fill path
				d.ClosePath()
				d.PathRender(AI_Fill)
			case "Bc": // define gradient instance cap
			case "Bg":
				args := token.PopAll()
				gradient.Flag = toInt8(args[0])
			case "Bh": // xHilight yHilight angle length Bh
			case "Bm": // set gradient matrix
			case "Xm": // set linear gradient matrix
			case "BB":
				d.SetGradient(&gradient)
				args := token.Pop()
				if args == "0" {
					// no action
				} else if args == "1" {
					// stroke path
					d.PathRender(AI_Stroke)
				} else {
					// close and stroke path
					d.ClosePath()
					d.PathRender(AI_Stroke)
				}
				return
			}
		}
	}
}

func (r *AIReader) beginRaster(d Drawer) {
	var token lineToken

	XI := []byte("XI")
	var xiargs []byte
	for r.readLine() {
		line := r.Bytes()

		if bytes.Equal(AI5_EndRaster, line) {
			break
		}

		if line[0] == '[' && line[len(line)-1] != 'h' {
			xiargs = make([]byte, len(line))
			copy(xiargs, line)
		}

		if bytes.Equal(line, XI) && len(xiargs) > 0 {
			line = append(xiargs, ' ', 'X', 'I')
		}

		token.parse(line)
		for token.len > 0 {
			op := token.Pop()
			switch op {
			case "Xh":
			case "XF":
			case "XG":
			case "XI":
				data := r.readRasterData()
				if obj := XIArgs(token.PopAll()); obj != nil {
					obj.RawData = data
					d.SetRaster(obj)
				}
			}
		}
	}
}
