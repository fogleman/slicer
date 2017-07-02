package main

import (
	"fmt"
	"math"
	"os"

	"github.com/fogleman/fauxgl"
)

const D = 5

type Key struct {
	Theta, Phi int
}

func MakeKey(v fauxgl.Vector) Key {
	theta := fauxgl.Round(fauxgl.Degrees(math.Acos(v.Z))/D) * D
	phi := fauxgl.Round(fauxgl.Degrees(math.Atan2(v.Y, v.X))/D) * D
	phi = (phi + 360) % 360
	return Key{theta, phi}
}

func (key Key) Opposite() Key {
	theta := 180 - key.Theta
	phi := (key.Phi + 180) % 360
	return Key{theta, phi}
}

func (key Key) Vector() fauxgl.Vector {
	theta := fauxgl.Radians(float64(key.Theta))
	phi := fauxgl.Radians(float64(key.Phi))
	x := math.Sin(theta) * math.Cos(phi)
	y := math.Sin(theta) * math.Sin(phi)
	z := math.Cos(theta)
	return fauxgl.Vector{x, y, z}
}

func main() {
	mesh, err := fauxgl.LoadMesh(os.Args[1])
	if err != nil {
		panic(err)
	}
	fmt.Println(1)

	lookup1 := make(map[Key]float64)
	for _, t := range mesh.Triangles {
		n := t.Normal()
		if math.IsNaN(n.Length()) {
			continue
		}
		k := MakeKey(n)
		a := t.Area()
		lookup1[k] += a
	}

	lookup2 := make(map[Key]float64)
	for key1, a := range lookup1 {
		for theta := 0; theta <= 180; theta += D {
			for phi := 0; phi < 360; phi += D {
				key2 := Key{theta, phi}
				dot := key1.Vector().Dot(key2.Vector())
				if dot < -1 {
					dot = -1
				}
				if dot > 1 {
					dot = 1
				}
				p := 1 - math.Acos(dot)/math.Pi
				p = math.Pow(p, 16)
				lookup2[key2] += a * p
				// lookup2[key2.Opposite()] += a * p
			}
		}
	}

	// var bestKey Key
	// bestScore := math.Inf(1)
	// for k, v := range lookup2 {
	// 	if v < bestScore {
	// 		bestScore = v
	// 		bestKey = k
	// 	}
	// }
	// fmt.Println(bestKey)

	// mesh.Transform(fauxgl.RotateTo(bestKey.Vector(), fauxgl.Vector{0, 0, 1}))
	// mesh.SaveSTL("out.stl")

	sphere := fauxgl.NewSphere2(8)
	for _, t := range sphere.Triangles {
		t.V1.Position = t.V1.Position.MulScalar(lookup2[MakeKey(t.V1.Position)])
		t.V2.Position = t.V2.Position.MulScalar(lookup2[MakeKey(t.V2.Position)])
		t.V3.Position = t.V3.Position.MulScalar(lookup2[MakeKey(t.V3.Position)])
	}
	sphere.SaveSTL("sphere.stl")
}
