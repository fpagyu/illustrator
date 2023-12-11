package illustrator

import (
	"bytes"
	"compress/zlib"
	"io"

	"github.com/klauspost/compress/zstd"
)

type CompressHandle interface {
	Decompress() ([]byte, error)

	Write(stream []byte) (n int, err error)
}

type ZlibCompress struct {
	buf *bytes.Buffer
}

func (c *ZlibCompress) Write(stream []byte) (int, error) {
	if c.buf == nil {
		c.buf = bytes.NewBuffer(stream)
		return len(stream), nil
	}

	return c.buf.Write(stream)
}

func (c *ZlibCompress) Decompress() ([]byte, error) {
	rd, err := zlib.NewReader(c.buf)
	if err != nil {
		return nil, err
	}
	defer rd.Close()

	return io.ReadAll(rd)
}

// ZStdCompress
type ZStdCompress struct {
	buf *bytes.Buffer
}

func (c *ZStdCompress) Write(stream []byte) (int, error) {
	if c.buf == nil {
		c.buf = bytes.NewBuffer(stream)
		return len(stream), nil
	}

	return c.buf.Write(stream)
}

func (c *ZStdCompress) Decompress() ([]byte, error) {
	rd, err := zstd.NewReader(nil)
	if err != nil {
		return nil, err
	}
	defer rd.Close()

	return rd.DecodeAll(c.buf.Bytes(), nil)
}
