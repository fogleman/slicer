package slicer

import (
	"bytes"
	"fmt"
	"strings"
)

type Layer struct {
	Z     float64
	Paths []Path
}

func (layer Layer) SVG() string {
	var buf bytes.Buffer
	for _, path := range layer.Paths {
		for i, point := range path {
			if i == 0 {
				buf.WriteString("M ")
			} else {
				buf.WriteString("L ")
			}
			buf.WriteString(fmt.Sprintf("%g %g ", point.X, point.Y))
		}
		buf.WriteString("Z ")
	}
	return strings.TrimSpace(buf.String())
}
