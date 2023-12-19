package svg

import "github.com/fpagyu/illustrator"

type SvgGradient struct {
	Defs      []illustrator.Gradient
	Instances []illustrator.Gradient
}

type OffColor struct {
	offset float64
	// opacity float64
	color []float64
}
