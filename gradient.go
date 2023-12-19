package illustrator

type Gradient struct {
	Flag         int8 // 1-issue a clip; 2-disable rending
	GradientType int8 // 渐变类型0-linear; 1-radial

	Name    string
	nColors int // number of colors in gradient

	Colors []OffColor
}

func (g *Gradient) AddColor(color *OffColor) {
	g.Colors = append(g.Colors, *color)
}

type OffColor struct {
	colorSpace int8

	offset   float64 // rampPoint
	midPoint float64 // midPoint
	color    []float64
}

func (sc *OffColor) Offset() float64 {
	return sc.offset
}

func (sc *OffColor) Color() string {
	return ""
}

func (sc *OffColor) Opacity() float64 {
	return 1.0
}

func BSArgs(vals []string) *OffColor {
	var stop OffColor
	var l = len(vals)

	stop.offset = toFloat(vals[l-1]) / 100
	stop.midPoint = toFloat(vals[l-2])

	colorSpace := vals[l-3]
	switch colorSpace {
	case "0":
		if l-3 == 1 {
		} else {
			colorSpace = vals[l-5]
		}
	case "1":
		if l-3 == 4 {
		} else {
			colorSpace = vals[l-5]
		}
	case "2":
		if l-3 == 7 {
		} else {
			colorSpace = vals[l-5]
		}
	case "3":
		if l-3 == 6 {
		} else {
			colorSpace = vals[l-5]
		}
	case "4":
		if l-3 == 10 {
		} else {
			colorSpace = vals[l-5]
		}
	default:
		colorSpace = vals[l-5]
	}

	switch colorSpace {
	case "0": // gray
		stop.colorSpace = 0
		gray := toFloat(vals[0])
		stop.color = []float64{gray}
	case "1": // cmyk
		stop.colorSpace = 1
		stop.color = []float64{
			toFloat(vals[0]), toFloat(vals[1]),
			toFloat(vals[2]), toFloat(vals[3]),
		}
	case "2": // rgb
		stop.colorSpace = 2
		stop.color = []float64{
			toFloat(vals[4]),
			toFloat(vals[5]),
			toFloat(vals[6]),
		}
	case "3": // custom cmyk
		stop.colorSpace = 3
		tint := toFloat(vals[5])
		stop.color = []float64{
			toFloat(vals[0])*(1-tint) + tint,
			toFloat(vals[1])*(1-tint) + tint,
			toFloat(vals[2])*(1-tint) + tint,
			toFloat(vals[3])*(1-tint) + tint,
		}

	case "4": // custom rgb
		stop.colorSpace = 4
		tint := toFloat(vals[8])
		stop.color = []float64{
			toFloat(vals[4])*(1-tint) + tint,
			toFloat(vals[5])*(1-tint) + tint,
			toFloat(vals[6])*(1-tint) + tint,
		}
	default:
		return nil
	}

	return &stop
}
