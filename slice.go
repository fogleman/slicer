package slicer

import (
	"math"
	"runtime"
	"sort"
	"sync"

	"github.com/fogleman/fauxgl"
)

func SliceMesh(m *fauxgl.Mesh, step float64) []Layer {
	wn := runtime.NumCPU()
	minz := m.BoundingBox().Min.Z
	maxz := m.BoundingBox().Max.Z

	// copy triangles
	triangles := make([]*triangle, len(m.Triangles))
	var wg sync.WaitGroup
	for wi := 0; wi < wn; wi++ {
		wg.Add(1)
		go func(wi int) {
			for i := wi; i < len(m.Triangles); i += wn {
				triangles[i] = newTriangle(m.Triangles[i])
			}
			wg.Done()
		}(wi)
	}
	wg.Wait()

	// sort triangles
	sort.Slice(triangles, func(i, j int) bool {
		return triangles[i].MinZ < triangles[j].MinZ
	})

	// create jobs for workers
	n := int(math.Ceil((maxz - minz) / step))
	in := make(chan job, n)
	out := make(chan Layer, n)
	for wi := 0; wi < wn; wi++ {
		go worker(in, out)
	}
	index := 0
	var active []*triangle
	for i := 0; i < n; i++ {
		z := fauxgl.RoundPlaces(minz+step*float64(i), 8)
		// remove triangles below plane
		newActive := active[:0]
		for _, t := range active {
			if t.MaxZ >= z {
				newActive = append(newActive, t)
			}
		}
		active = newActive
		// add triangles above plane
		for index < len(triangles) && triangles[index].MinZ <= z {
			active = append(active, triangles[index])
			index++
		}
		// copy triangles for worker job
		activeCopy := make([]*triangle, len(active))
		copy(activeCopy, active)
		in <- job{z, activeCopy}
	}
	close(in)

	// read results from workers
	layers := make([]Layer, n)
	for i := 0; i < n; i++ {
		layers[i] = <-out
	}

	// sort layers
	sort.Slice(layers, func(i, j int) bool {
		return layers[i].Z < layers[j].Z
	})
	return layers
}

type job struct {
	Z         float64
	Triangles []*triangle
}

func worker(in chan job, out chan Layer) {
	var paths []Path
	for j := range in {
		paths = paths[:0]
		for _, t := range j.Triangles {
			if v1, v2, ok := intersectTriangle(j.Z, t); ok {
				paths = append(paths, Path{v1, v2})
			}
		}
		out <- Layer{j.Z, joinPaths(paths)}
	}
}

func intersectSegment(z float64, v0, v1 fauxgl.Vector) (fauxgl.Vector, bool) {
	if v0.Z == v1.Z {
		return fauxgl.Vector{}, false
	}
	t := (z - v0.Z) / (v1.Z - v0.Z)
	if t < 0 || t > 1 {
		return fauxgl.Vector{}, false
	}
	v := v0.Add(v1.Sub(v0).MulScalar(t))
	return v, true
}

func intersectTriangle(z float64, t *triangle) (fauxgl.Vector, fauxgl.Vector, bool) {
	v1, ok1 := intersectSegment(z, t.V1, t.V2)
	v2, ok2 := intersectSegment(z, t.V2, t.V3)
	v3, ok3 := intersectSegment(z, t.V3, t.V1)
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
