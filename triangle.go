package slicer

import "github.com/fogleman/fauxgl"

type triangle struct {
	N, V1, V2, V3 fauxgl.Vector
	MinZ, MaxZ    float64
}

func newTriangle(ft *fauxgl.Triangle) *triangle {
	t := triangle{}
	// micro-adjust the z coordinates to avoid intersecting plane at exact tip
	t.V1 = ft.V1.Position.RoundPlaces(9)
	t.V2 = ft.V2.Position.RoundPlaces(9)
	t.V3 = ft.V3.Position.RoundPlaces(9)
	t.V1.Z += 5e-10
	t.V2.Z += 5e-10
	t.V3.Z += 5e-10
	// compute triangle normal
	e1 := t.V2.Sub(t.V1)
	e2 := t.V3.Sub(t.V1)
	t.N = e1.Cross(e2).Normalize()
	// compute min and max z value
	t.MinZ = t.V1.Z
	if t.V2.Z < t.MinZ {
		t.MinZ = t.V2.Z
	}
	if t.V3.Z < t.MinZ {
		t.MinZ = t.V3.Z
	}
	t.MaxZ = t.V1.Z
	if t.V2.Z > t.MaxZ {
		t.MaxZ = t.V2.Z
	}
	if t.V3.Z > t.MaxZ {
		t.MaxZ = t.V3.Z
	}
	return &t
}
