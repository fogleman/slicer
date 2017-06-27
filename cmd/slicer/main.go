package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/fogleman/fauxgl"
	"github.com/fogleman/gg"
	"github.com/fogleman/slicer"
)

func timed(name string) func() {
	if len(name) > 0 {
		fmt.Printf("%s... ", name)
	}
	start := time.Now()
	return func() {
		fmt.Println(time.Since(start))
	}
}

func main() {
	var done func()

	done = timed("loading mesh")
	mesh, err := fauxgl.LoadMesh(os.Args[1])
	if err != nil {
		panic(err)
	}
	box := mesh.BoundingBox()
	step := 0.1
	done()

	done = timed("slicing mesh")
	layers := slicer.SliceMesh(mesh, step)
	done()

	// return

	done = timed("rendering layers")
	wn := runtime.NumCPU()
	ch := make(chan job, len(layers))
	var wg sync.WaitGroup
	for wi := 0; wi < wn; wi++ {
		wg.Add(1)
		go worker(ch, &wg)
	}
	for i, l := range layers {
		ch <- job{i, l, box}
	}
	close(ch)
	wg.Wait()
	done()
}

type job struct {
	i     int
	layer slicer.Layer
	box   fauxgl.Box
}

func worker(ch chan job, wg *sync.WaitGroup) {
	const S = 1600
	const P = 50
	for j := range ch {
		i := j.i
		layer := j.layer
		box := j.box
		center := box.Center()
		size := box.Size()
		sx := (S - P*2) / size.X
		sy := (S - P*2) / size.Y
		scale := math.Min(sx, sy)
		dc := gg.NewContext(S, S)
		dc.InvertY()
		dc.SetRGB(1, 1, 1)
		dc.Clear()
		dc.Translate(S/2, S/2)
		dc.Scale(scale, scale)
		dc.Translate(-center.X, -center.Y)
		dc.SetFillRuleWinding()
		for _, path := range layer.Paths {
			dc.NewSubPath()
			for _, point := range path {
				dc.LineTo(point.X, point.Y)
			}
			dc.ClosePath()
		}
		dc.SetRGB(0, 0, 0)
		dc.Fill()
		// dc.SetRGB(0.6, 0.6, 0.6)
		// dc.FillPreserve()
		// dc.SetRGB(0, 0, 0)
		// dc.SetLineWidth(3)
		// dc.Stroke()
		dc.SavePNG(fmt.Sprintf("out%04d.png", i))
	}
	wg.Done()
}
