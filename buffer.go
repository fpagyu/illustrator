package illustrator

type Buffer struct {
	buf []byte
	len int
}

func (w *Buffer) grow(n int) {
	n = len(w.buf) + n
	if n > 64*1024 {
		panic("too large to grow")
	}

	buf := make([]byte, n)
	copy(buf, w.buf)
	w.buf = buf
}

func (w *Buffer) WriteByte(b byte) error {
	if w.len >= len(w.buf) {
		w.grow(512)
	}

	w.buf[w.len] = b
	w.len++

	return nil
}

func (w *Buffer) Reset() {
	w.len = 0
}

func (w *Buffer) Len() int {
	return w.len
}

func (w *Buffer) Bytes() []byte {
	return w.buf[:w.len]
}

func (w *Buffer) Text() string {
	return string(w.buf[:w.len])
}
