package ps

import (
	"bytes"
	"strconv"
)

type AIHeader struct {
	Title            string
	BoundingBox      [4]int
	HiResBoundingBox [4]float64
}

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
	}
}
