package illustrator

type PathOp int8

const (
	AI_Fill   PathOp = 1 << 0
	AI_Stroke PathOp = 1 << 1
)

const (
	AI_ClipPath     uint8 = 1 << 0
	AI_CompoundPath uint8 = 1 << 1
	AI_RasterImage  uint8 = 1 << 2
)
