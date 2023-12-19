package svg

type SvgImage struct {
	id     string
	indent int
	parent *SvgGroup

	width  int
	height int
	matrix [6]float64

	styles string

	// data []byte
	b64Img string
}

func (si *SvgImage) Id() string {
	return si.id
}

func (si *SvgImage) Indent() int {
	return si.indent
}

func (si *SvgImage) SetAttr(key, val string) {
}

func (si *SvgImage) SetImage(b64data string) {
	si.b64Img = b64data
}
