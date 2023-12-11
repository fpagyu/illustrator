package svg

type SvgNode interface {
	Id() string
	Indent() int // 元素层级/缩进
}
