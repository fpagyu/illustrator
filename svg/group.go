package svg

type SvgGroup struct {
	id      string
	indent  int // 层级
	visible bool

	childs []SvgNode // 子节点
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

func (g *SvgGroup) AddNode(node SvgNode) {
	g.childs = append(g.childs, node)
}
