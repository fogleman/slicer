package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/fogleman/fauxgl"
	embree "github.com/fogleman/go-embree"
)

func timed(name string) func() {
	fmt.Printf("%s... ", name)
	start := time.Now()
	return func() {
		fmt.Println(time.Since(start))
	}
}

func fauxglToEmbree(mesh *fauxgl.Mesh) *embree.Mesh {
	triangles := make([]embree.Triangle, len(mesh.Triangles))
	for i, t := range mesh.Triangles {
		triangles[i] = embree.Triangle{
			embree.Vector{t.V1.Position.X, t.V1.Position.Y, t.V1.Position.Z},
			embree.Vector{t.V2.Position.X, t.V2.Position.Y, t.V2.Position.Z},
			embree.Vector{t.V3.Position.X, t.V3.Position.Y, t.V3.Position.Z},
		}
	}
	return embree.NewMesh(triangles)
}

func main() {
	var done func()

	done = timed("creating sphere")
	sphere := fauxgl.NewSphere(6)
	hitSphere := fauxgl.NewSphere(3)
	embreeSphere := fauxglToEmbree(hitSphere)
	spherePoints := make(map[fauxgl.Vector]bool)
	for _, t := range sphere.Triangles {
		spherePoints[t.V1.Position] = true
		spherePoints[t.V2.Position] = true
		spherePoints[t.V3.Position] = true
	}
	done()

	done = timed("loading mesh")
	mesh, err := fauxgl.LoadMesh(os.Args[1])
	if err != nil {
		panic(err)
	}
	done()

	done = timed("first pass")
	lookup1 := make(map[fauxgl.Vector]float64)
	for _, t := range mesh.Triangles {
		n := t.Normal()
		a := t.Area()
		if math.IsNaN(n.Length()) {
			continue
		}
		ray := embree.Ray{embree.Vector{}, embree.Vector{n.X, n.Y, n.Z}}
		hit := embreeSphere.Intersect(ray)
		p := n.MulScalar(hit.T)
		st := hitSphere.Triangles[hit.Index]
		p1 := st.V1.Position
		p2 := st.V2.Position
		p3 := st.V3.Position
		b := fauxgl.Barycentric(p1, p2, p3, p)
		lookup1[p1] += a * b.X
		lookup1[p2] += a * b.Y
		lookup1[p3] += a * b.Z
	}
	done()

	done = timed("second pass")
	lookup2 := make(map[fauxgl.Vector]float64)
	for p1, a := range lookup1 {
		for p2 := range spherePoints {
			p := p1.X*p2.X + p1.Y*p2.Y + p1.Z*p2.Z
			if p < 0.5 {
				continue
			}
			if p >= 1 {
				p = 1
			} else {
				p = math.Pow(p, 32)
			}
			lookup2[p2] += a * p
		}
	}
	done()

	done = timed("creating oriented.stl")
	var bestVector fauxgl.Vector
	bestScore := math.Inf(1)
	for k, v := range lookup2 {
		if v < bestScore {
			bestScore = v
			bestVector = k
		}
	}
	mesh.Transform(fauxgl.RotateTo(bestVector, fauxgl.Vector{0, 0, 1}))
	mesh.SaveSTL("oriented.stl")
	done()

	done = timed("creating normals.stl")
	for _, t := range sphere.Triangles {
		t.V1.Position = t.V1.Position.MulScalar(lookup2[t.V1.Position])
		t.V2.Position = t.V2.Position.MulScalar(lookup2[t.V2.Position])
		t.V3.Position = t.V3.Position.MulScalar(lookup2[t.V3.Position])
	}
	sphere.SaveSTL("normals.stl")
	done()
}
