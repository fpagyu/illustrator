package illustrator

type Drawer interface {
	SetHeader(*AIHeader)

	// Layer
	BeginLayer(*AILayer)
	SetLayerName(name string)
	EndLayer()

	// Group
	Group()
	EndGroup()
	SetGroupAttr()

	// path
	ClosePath()
	PathRender(PathOp)
	Moveto(x, y float64)
	Lineto(x, y float64)
	Curveto1(x0, y0, x1, y1 float64)
	Curveto2(x1, y1, x2, y2 float64)
	Curveto(x0, y0, x1, y1, x2, y2 float64)

	// clip path
	ClipPath()
	ApplyClip()

	// compound path
	CompoundPath()
	EndCompoundPath()

	// set color
	SetRGB(PathOp, [3]uint8)
	SetCMYK(PathOp, [4]float64)
	SetOpacity(opacity string)

	// path attributes
	SetDash()
	SetFlat()
	SetLineCap(v string)
	SetLineJoin(v string)
	SetLineWidth(v string)
	SetMiterLimit(v string)

	// gradient
	DefGradient(g *Gradient) //
	SetGradient(g *Gradient)

	// raster
	SetRaster(obj *Raster)
}
