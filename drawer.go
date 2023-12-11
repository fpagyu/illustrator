package illustrator

import "github.com/fpagyu/illustrator/ps"

type Drawer interface {
	SetHeader(*ps.AIHeader)

	// Layer
	BeginLayer(*ps.AILayer)
	SetLayerName(name string)
	EndLayer()
}
