package illustrator

type AILayer struct {
	Name       string
	LayerIndex int // layer index

	Visible           bool
	Preview           bool
	Enabled           bool
	Printing          bool
	Dimmed            bool
	HasMultiLayerMask bool
	ColorIndex        int8 // between -1 and 26
	RGB               [3]uint8
}
