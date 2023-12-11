package illustrator

type lineToken struct {
	stack []string
	len   int
}

func (tk *lineToken) parse(line []byte) bool {
	var left int
	for i := range line {
		if line[left] == '(' {
			if line[i] == ')' && line[i-1] != '\\' {
				tk.Push(string(line[left : i+1]))
				left = i + 1
			}
		} else {
			if line[i] == ' ' {
				if left < i {
					tk.Push(string(line[left:i]))
				}
				left = i + 1
			}
		}
	}

	if left < len(line) {
		tk.Push(string(line[left:]))
	}

	return tk.len > 0
}

func (tk *lineToken) Push(v string) {
	if len(tk.stack) == 0 {
		tk.len = 0
		tk.stack = make([]string, 32)
	}

	if tk.len >= len(tk.stack) {
		panic("stack overflow")
	}

	tk.stack[tk.len] = v
	tk.len++
}

func (tk *lineToken) Pop() (v string) {
	if tk.len > 0 {
		tk.len--
		v = tk.stack[tk.len]
	}
	return
}

func (tk *lineToken) PopN(n int) (vals []string) {
	if l := tk.len; l >= n {
		tk.len -= n
		vals = tk.stack[tk.len:l]
	}

	return vals
}

func (tk *lineToken) PopAll() (vals []string) {
	if tk.len > 0 {
		vals = tk.stack[0:tk.len]
		tk.len = 0
	}

	return
}

func (tk *lineToken) Top() (v string) {
	if tk.len > 0 {
		v = tk.stack[tk.len-1]
	}
	return
}
