package illustrator

import (
	"bytes"
	"strconv"
)

type AIHeader struct {
	Title            string
	BoundingBox      [4]int
	HiResBoundingBox [4]float64
}

type AIProlog struct{}

func (h *AIHeader) SetHeader(title []byte) {
	var start, end int
	for i := range title {
		if title[i] == '(' {
			start = i + 1
		} else if title[i] == ')' {
			end = i
		}
	}

	if end > start {
		h.Title = string(title[start:end])
	}
}

func (h *AIHeader) SetBoundingBox(line []byte) {
	vals := bytes.Split(line, []byte{' '})

	var i int
	for _, v := range vals {
		if len(v) == 0 {
			continue
		}

		h.BoundingBox[i], _ = strconv.Atoi(string(v))
		i++
	}
}

func (h *AIHeader) SetHiResBoundingBox(line []byte) {
	vals := bytes.Split(line, []byte{' '})

	var i int
	for _, v := range vals {
		if len(v) == 0 {
			continue
		}

		h.HiResBoundingBox[i], _ = strconv.ParseFloat(string(v), 0)
		i++
	}
}

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
