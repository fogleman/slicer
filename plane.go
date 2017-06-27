package slicer

import "github.com/fogleman/fauxgl"

type plane struct {
	Point  fauxgl.Vector
	Normal fauxgl.Vector
}

func (p *plane) intersectSegment(v0, v1 fauxgl.Vector) (fauxgl.Vector, bool) {
	if v0.Z == v1.Z {
		return fauxgl.Vector{}, false
	}
	t := (p.Point.Z - v0.Z) / (v1.Z - v0.Z)
	if t < 0 || t > 1 {
		return fauxgl.Vector{}, false
	}
	v := v0.Add(v1.Sub(v0).MulScalar(t))
	return v, true
}

func (p *plane) intersectTriangle(t *triangle) (fauxgl.Vector, fauxgl.Vector, bool) {
	v1, ok1 := p.intersectSegment(t.V1, t.V2)
	v2, ok2 := p.intersectSegment(t.V2, t.V3)
	v3, ok3 := p.intersectSegment(t.V3, t.V1)
	var p1, p2 fauxgl.Vector
	if ok1 && ok2 {
		p1, p2 = v1, v2
	} else if ok1 && ok3 {
		p1, p2 = v1, v3
	} else if ok2 && ok3 {
		p1, p2 = v2, v3
	} else {
		return fauxgl.Vector{}, fauxgl.Vector{}, false
	}
	p1 = p1.RoundPlaces(8)
	p2 = p2.RoundPlaces(8)
	if p1 == p2 {
		return fauxgl.Vector{}, fauxgl.Vector{}, false
	}
	n := fauxgl.Vector{p1.Y - p2.Y, p2.X - p1.X, 0}
	if n.Dot(t.N) < 0 {
		return p1, p2, true
	} else {
		return p2, p1, true
	}
}
