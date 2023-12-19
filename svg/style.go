package svg

import (
	"strings"
)

type StyleBuild map[string]string

func (b StyleBuild) fill() string {
	var w strings.Builder
	for k, v := range b {
		if strings.HasPrefix(k, "stroke") {
			continue
		}
		w.WriteString(k)
		w.WriteByte(':')
		w.WriteString(v)
		w.WriteByte(';')
	}

	return w.String()
}

func (b StyleBuild) stroke() string {
	var w strings.Builder
	w.WriteString("fill:none;")
	for k, v := range b {
		if strings.HasPrefix(k, "fill") {
			continue
		}
		w.WriteString(k)
		w.WriteByte(':')
		w.WriteString(v)
		w.WriteByte(';')
	}

	return w.String()
}

func (b StyleBuild) styles() string {
	var w strings.Builder
	for k, v := range b {
		w.WriteString(k)
		w.WriteByte(':')
		w.WriteString(v)
		w.WriteByte(';')
	}

	return w.String()
}

func (b StyleBuild) nofillstroke() string {
	var w strings.Builder

	for k, v := range b {
		if strings.HasPrefix(k, "fill") {
			continue
		}

		if strings.HasPrefix(k, "stroke") {
			continue
		}

		w.WriteString(k)
		w.WriteByte(':')
		w.WriteString(v)
		w.WriteByte(';')
	}

	return w.String()
}
