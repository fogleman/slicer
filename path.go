package slicer

import (
	"math"

	"github.com/fogleman/fauxgl"
)

type Path []fauxgl.Vector

// func (path Path) reverse() {
// 	for i := len(path)/2 - 1; i >= 0; i-- {
// 		opp := len(path) - 1 - i
// 		path[i], path[opp] = path[opp], path[i]
// 	}
// }

func (path Path) isClockwise() bool {
	maxy := path[0].Y
	maxi := 0
	for i, p := range path {
		if p.Y > maxy {
			maxy = p.Y
			maxi = i
		}
	}
	p1 := path[(maxi+len(path)-1)%len(path)]
	p2 := path[maxi]
	p3 := path[(maxi+1)%len(path)]
	if p1.X <= p2.X && p2.X <= p3.X {
		return true
	}
	if p1.X >= p2.X && p2.X >= p3.X {
		return false
	}
	dx1 := p2.X - p1.X
	dy1 := p2.Y - p1.Y
	dx2 := p3.X - p2.X
	dy2 := p3.Y - p2.Y
	m := fauxgl.Matrix{dx1, dx2, 0, 0, dy1, dy2, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
	d := m.Determinant()
	if math.Abs(d) < 1e-9 {
		return p1.X < p3.X
	}
	return d < 0
}

func joinPaths(paths []Path) []Path {
	lookup := make(map[fauxgl.Vector]Path)
	for _, path := range paths {
		lookup[path[0]] = path
	}
	var result []Path
	for len(lookup) > 0 {
		var v fauxgl.Vector
		for v = range lookup {
			break
		}
		var path Path
		for {
			path = append(path, v)
			if p, ok := lookup[v]; ok {
				delete(lookup, v)
				v = p[len(p)-1]
			} else {
				break
			}
		}
		result = append(result, path)
	}
	return result
}
