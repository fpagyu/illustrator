package svg

import "fmt"

type SvgNode interface {
	Id() string
	Indent() int // 元素层级/缩进
	SetAttr(k, v string)
}

func Attr(k, v string) string {
	return fmt.Sprintf(`%s="%s"`, k, v)
}
