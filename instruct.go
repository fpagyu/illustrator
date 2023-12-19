package illustrator

import (
	"math"
)

type ColorArgs struct {
	cmyk [4]float64
	rgb  [3]float64

	gray float64 //

	tint       float64 // between 0-1
	colorSpace int8    // 0-CMYK, 1-RGB
}

func (args *ColorArgs) SetTint(tint float64) {
	if tint < 1e-6 || tint > 1.0 {
		return
	}

	if args.colorSpace == 1 {
		args.rgb[0] = args.rgb[0]*(1-tint) + tint // r
		args.rgb[1] = args.rgb[1]*(1-tint) + tint // g
		args.rgb[2] = args.rgb[2]*(1-tint) + tint // b
	} else {
		c, m := args.cmyk[0], args.cmyk[1]
		y, k := args.cmyk[2], args.cmyk[3]

		args.cmyk[0] = c*(1-tint) + tint
		args.cmyk[1] = m*(1-tint) + tint
		args.cmyk[2] = y*(1-tint) + tint
		args.cmyk[3] = k*(1-tint) + tint
	}
}

func (args *ColorArgs) RGB() [3]uint8 {
	return [3]uint8{
		uint8(math.Round(args.rgb[0] * 255)),
		uint8(math.Round(args.rgb[1] * 255)),
		uint8(math.Round(args.rgb[2] * 255)),
	}
}

func (args *ColorArgs) CMYK() [4]float64 {
	return args.cmyk
}

func XAArgs(vals []string) *ColorArgs {
	var args ColorArgs
	if len(vals) == 3 {
		args.rgb[0] = toFloat(vals[0])
		args.rgb[1] = toFloat(vals[1])
		args.rgb[2] = toFloat(vals[2])
	} else if len(vals) == 7 {
		args.cmyk[0] = toFloat(vals[0])
		args.cmyk[1] = toFloat(vals[1])
		args.cmyk[2] = toFloat(vals[2])
		args.cmyk[3] = toFloat(vals[3])

		args.rgb[0] = toFloat(vals[4])
		args.rgb[1] = toFloat(vals[5])
		args.rgb[2] = toFloat(vals[6])
	}

	return &args
}

func XKArgs(vals []string) *ColorArgs {
	var args ColorArgs
	if l := len(vals); l == 10 || l == 7 {
		switch vals[l-1] {
		case "1": // rgb
			args.colorSpace = 1
			args.rgb[0] = toFloat(vals[4])
			args.rgb[1] = toFloat(vals[5])
			args.rgb[2] = toFloat(vals[6])
		case "0": // cmyk
			args.colorSpace = 0
			args.cmyk[0] = toFloat(vals[0])
			args.cmyk[1] = toFloat(vals[1])
			args.cmyk[2] = toFloat(vals[2])
			args.cmyk[3] = toFloat(vals[3])
		}
	}

	args.SetTint(toFloat(vals[len(vals)-2]))
	return &args
}

func XYArgs(vals []string) (opacity string) {
	if len(vals) == 5 {
		return vals[1]
	}

	return "1"
}

func KArgs(vals []string) *ColorArgs {
	// cyan magenta yellow black K
	var args ColorArgs
	args.colorSpace = 0 // set cmyk color space
	if len(vals) == 4 {
		args.cmyk[0] = toFloat(vals[0])
		args.cmyk[1] = toFloat(vals[1])
		args.cmyk[2] = toFloat(vals[2])
		args.cmyk[3] = toFloat(vals[3])
	} else {
		panic("invalid k arguments")
	}

	return &args
}

func XArgs(vals []string) *ColorArgs {
	// cyan magenta yellow black (name) gray X
	var args ColorArgs
	args.colorSpace = 0 // set cmyk color space
	if len(vals) == 6 {
		args.cmyk[0] = toFloat(vals[0])
		args.cmyk[1] = toFloat(vals[1])
		args.cmyk[2] = toFloat(vals[2])
		args.cmyk[3] = toFloat(vals[3])
		args.SetTint(toFloat(vals[5]))
	} else {
		panic("invalid x arguments")
	}

	return &args
}

func XXArgs(vals []string) *ColorArgs {
	var args ColorArgs
	args.colorSpace = toInt8(vals[len(vals)-1])
	if args.colorSpace == 0 {
		args.rgb[0] = toFloat(vals[0])
		args.rgb[1] = toFloat(vals[1])
		args.rgb[2] = toFloat(vals[2])
	} else {
		args.cmyk[0] = toFloat(vals[0])
		args.cmyk[1] = toFloat(vals[1])
		args.cmyk[2] = toFloat(vals[2])
		args.cmyk[3] = toFloat(vals[3])
	}

	// set tint
	args.SetTint(toFloat(vals[len(vals)-2]))
	return &args
}
