package svg

import (
	"github.com/fpagyu/illustrator/ps"
)

type SVG struct {
	viewBox [4]int
	layers  []SvgNode // 对应ai文件的层

	group *SvgGroup // 当前所在的group
}

func (svg *SVG) SetHeader(header *ps.AIHeader) {
	svg.viewBox = header.BoundingBox
}

func (svg *SVG) BeginLayer(layer *ps.AILayer) {
	svg.group = &SvgGroup{
		id: layer.Name,
	}
	svg.layers = append(svg.layers, svg.group)
}

func (svg *SVG) EndLayer() {
	svg.group = nil
}

func (svg *SVG) SetLayerName(name string) {
	svg.group.id = name
}
