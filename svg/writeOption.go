package svg

type SvgWriteOption struct {
	IgnoreImage bool // 忽略位图数据, 保存为svg的时候, image数据不会写入
}

func SetIgnoreImage(v bool) func(*SvgWriteOption) {
	return func(swo *SvgWriteOption) {
		swo.IgnoreImage = v
	}
}
