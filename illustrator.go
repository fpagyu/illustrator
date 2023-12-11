package illustrator

import (
	"bufio"
	"bytes"
	"io"
)

type AIReader struct {
	*bufio.Reader

	lineBuf Buffer
}

func NewAIReader(r io.Reader) (*AIReader, error) {
	reader := &AIReader{
		Reader: bufio.NewReader(r),
	}

	return reader, nil
}

func (r *AIReader) Bytes() []byte {
	return r.lineBuf.Bytes()
}

func (r *AIReader) readLine() bool {
	var n int
	r.lineBuf.Reset()
	for {
		ch, err := r.ReadByte()
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		if ch == '\r' || ch == '\n' {
			if n > 0 {
				break
			}
		} else {
			r.lineBuf.WriteByte(ch)
			n++
		}
	}

	return r.lineBuf.len > 0
}

func (r *AIReader) Draw(drawer Drawer) error {
	header := r.readHeader()
	drawer.SetHeader(header)

	BeginSetup := []byte("%%BeginSetup")
	BeginProlog := []byte("%%BeginProlog")
	BeginLayer := []byte("%AI5_BeginLayer")
	for r.readLine() {
		line := r.Bytes()
		if bytes.HasPrefix(line, BeginSetup) {
			r.readSetup()
		} else if bytes.HasPrefix(line, BeginProlog) {
			r.readProlog()
		} else if bytes.HasPrefix(line, BeginLayer) {
			r.drawLayer(drawer)
		}
	}

	return nil
}

func (r *AIReader) readHeader() *AIHeader {
	var header AIHeader
	Title := []byte("%%Title:")
	EndComment := []byte("%%EndComments")
	BoundingBox := []byte("%%BoundingBox:")
	HiResBoundingBox := []byte("%%HiResBoundingBox:")
	for r.readLine() {
		line := r.Bytes()
		if bytes.HasSuffix(line, EndComment) {
			break
		}

		if bytes.HasPrefix(line, Title) {
			header.SetHeader(bytes.TrimPrefix(line, Title))
			continue
		}

		if bytes.HasPrefix(line, BoundingBox) {
			header.SetBoundingBox(bytes.TrimPrefix(line, BoundingBox))
			continue
		}

		if bytes.HasPrefix(line, HiResBoundingBox) {
			header.SetHiResBoundingBox(bytes.TrimPrefix(line, HiResBoundingBox))
			continue
		}
	}

	return &header
}

func (r *AIReader) readProlog() *AIProlog {
	// %%BeginProlog
	// %%EndProlog
	var prolog AIProlog
	EndProlog := []byte("%%EndProlog")
	for r.readLine() {
		line := r.Bytes()

		if bytes.HasPrefix(line, EndProlog) {
			break
		}
		// todo
	}

	return &prolog
}

func (r *AIReader) readSetup() {
	// %%BeginSetup
	// %%EndSetup
	EndSetup := []byte("%%EndSetup")
	for r.readLine() {
		line := r.Bytes()

		if bytes.HasPrefix(line, EndSetup) {
			break
		}
		// todo
	}
}

func (r *AIReader) drawLayer(d Drawer) {
	var token lineToken
	for r.readLine() {
		line := r.Bytes()
		if line[0] == '%' {
			// skip comments
			continue
		}

		token.parse(line)
		for token.len > 0 {
			op := token.Pop()
			switch op {
			case "A": // locking, 0-unlocking; 1-locking
				token.Pop()
			case "Ap": // show center point
				token.Pop()
			case "Lb":
				r.beginLayer(d, token.PopAll())
			case "LB":
				d.EndLayer()
			case "Ln":
				d.SetLayerName(token.Pop())
			case "O", "R": // fill/stroke overprint
				token.Pop()
			}
		}
	}
}

func (r *AIReader) beginLayer(d Drawer, args []string) {
	layer := AILayer{
		Visible:           args[0] == "1",
		Preview:           args[1] == "1",
		Enabled:           args[2] == "1",
		Printing:          args[3] == "1",
		Dimmed:            args[4] == "1",
		HasMultiLayerMask: args[5] == "1",
		ColorIndex:        toInt8(args[6]),
	}

	colorIndex := 7
	if len(args) > 10 {
		colorIndex = 8
		layer.Name = "Layer " + args[7]
		layer.LayerIndex = toInt(args[7])
	}
	for i := 0; i < 3; i++ {
		layer.RGB[i] = toUint8(args[colorIndex+i])
	}

	d.BeginLayer(&layer)
}
