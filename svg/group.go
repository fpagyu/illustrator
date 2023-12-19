package svg

type SvgGroup struct {
	id     string
	indent int // 层级

	parent *SvgGroup // 父节点
	childs []SvgNode // 子节点
	clips  []SvgPath // clip paths
	attrs  map[string]string
}

func (g *SvgGroup) Id() string {
	return g.id
}

func (g *SvgGroup) Indent() int {
	return g.indent
}

func (g *SvgGroup) Childs() []SvgNode {
	return g.childs
}

func (g *SvgGroup) SetAttr(k, v string) {
	if len(k) == 0 || len(v) == 0 {
		return
	}

	if g.attrs == nil {
		g.attrs = make(map[string]string)
	}
	g.attrs[k] = v
}

func (g *SvgGroup) AddNode(node SvgNode) {
	g.childs = append(g.childs, node)
}

func (g *SvgGroup) Attrs() []string {
	r := make([]string, 0, len(g.attrs)+1)
	if len(g.id) > 0 {
		r = append(r, Attr("id", g.id))
	}

	for k, v := range g.attrs {
		r = append(r, Attr(k, v))
	}

	return r
}
