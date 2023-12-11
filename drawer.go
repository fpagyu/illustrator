package illustrator

type Drawer interface {
	SetHeader(*AIHeader)

	// Layer
	BeginLayer(*AILayer)
	SetLayerName(name string)
	EndLayer()
}
