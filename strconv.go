package illustrator

import "strconv"

func toInt(s string) int {
	v, _ := strconv.ParseInt(s, 10, 0)
	return int(v)
}

func toInt8(s string) int8 {
	v, _ := strconv.ParseInt(s, 10, 0)
	return int8(v)
}

func toUint8(s string) uint8 {
	v, _ := strconv.ParseUint(s, 10, 0)
	return uint8(v)
}
