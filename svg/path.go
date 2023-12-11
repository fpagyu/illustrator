package svg

type SvgPath struct {
	id     string
	indent int // 层级
}

func (sp *SvgPath) Id() string {
	return sp.id
}

func (sp *SvgPath) Indent() int {
	return sp.indent
}
