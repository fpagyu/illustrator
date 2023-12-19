package svg

import (
	"strconv"
)

func Float(v float64) string {
	s := strconv.FormatFloat(v, 'f', 2, 64)

	i := len(s) - 1
	for s[i] == '0' {
		i--
	}
	if s[i] == '.' {
		return s[0:i]
	} else {
		return s[0 : i+1]
	}
}
