package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/fogleman/fauxgl"
	"github.com/fogleman/slicer"
)

var (
	quiet     = kingpin.Flag("quiet", "Run in silent mode.").Short('q').Bool()
	size      = kingpin.Flag("size", "Set the slice thickness.").Short('s').Float()
	number    = kingpin.Flag("number", "Set the number of slices.").Short('n').Int()
	directory = kingpin.Flag("directory", "Set the output directory.").Short('d').ExistingDir()
	files     = kingpin.Arg("files", "Mesh files to slice.").Required().ExistingFiles()
)

func main() {
	kingpin.Parse()
	if *size == 0 && *number == 0 {
		kingpin.Fatalf("must specify slice thickness or count, try --help")
	}
	for _, filename := range *files {
		process(filename)
	}
}

func timed(name string) func() {
	if *quiet {
		return func() {}
	}
	fmt.Printf("%s... ", name)
	start := time.Now()
	return func() {
		fmt.Println(time.Since(start))
	}
}

func log(text string) {
	if !*quiet {
		fmt.Println(text)
	}
}

func process(infile string) {
	var done func()

	log(fmt.Sprintf("input: %s", infile))

	// load mesh
	done = timed("loading mesh")
	mesh, err := fauxgl.LoadMesh(infile)
	if err != nil {
		panic(err)
	}
	box := mesh.BoundingBox()
	done()

	// determine slice thickness
	var step float64
	if *size != 0 {
		step = *size
	}
	if *number != 0 {
		step = box.Size().Z / float64(*number)
	}

	// slice mesh
	done = timed("slicing mesh")
	layers := slicer.SliceMesh(mesh, step)
	done()

	// determine output filename
	dir, name := filepath.Split(infile)
	if *directory != "" {
		dir = *directory
	}
	outfile, _ := filepath.Abs(filepath.Join(dir, name) + ".svg")

	// write output
	done = timed("creating svg")
	svg := createSVG(box, layers)
	ioutil.WriteFile(outfile, []byte(svg), 0644)
	done()

	log(fmt.Sprintf("output: %s", outfile))
	log("")
}

func createSVG(box fauxgl.Box, layers []slicer.Layer) string {
	const S = 1600
	const P = 32

	center := box.Center()
	size := box.Size()
	sx := (S - P*2) / size.X
	sy := (S - P*2) / size.Y
	scale := math.Min(sx, sy)
	width := 0.25 / scale
	transform := fmt.Sprintf(
		"translate(%d %d) scale(%g) translate(%g %g)",
		S/2, S/2, scale, -center.X, -center.Y)

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("<svg version=\"1.1\" width=\"%d\" height=\"%d\" xmlns=\"http://www.w3.org/2000/svg\">\n", S, S))
	buf.WriteString(fmt.Sprintf("<g transform=\"%s\">\n", transform))
	for _, layer := range layers {
		buf.WriteString(fmt.Sprintf("<path slicer-z=\"%g\" stroke=\"#000000\" stroke-width=\"%g\" fill=\"#000000\" fill-rule=\"evenodd\" fill-opacity=\"0.01\" d=\"%s\"></path>\n", layer.Z, width, layer.SVG()))
	}
	buf.WriteString("</g>\n")
	buf.WriteString("</svg>\n")
	return buf.String()
}

// type job struct {
// 	i     int
// 	layer slicer.Layer
// 	box   fauxgl.Box
// }

// func render(box fauxgl.Box, layers []slicer.Layer) {
// 	wn := runtime.NumCPU()
// 	ch := make(chan job, len(layers))
// 	var wg sync.WaitGroup
// 	for wi := 0; wi < wn; wi++ {
// 		wg.Add(1)
// 		go worker(ch, &wg)
// 	}
// 	for i, l := range layers {
// 		ch <- job{i, l, box}
// 	}
// 	close(ch)
// 	wg.Wait()
// }

// func worker(ch chan job, wg *sync.WaitGroup) {
// 	const S = 1600
// 	const P = 50
// 	for j := range ch {
// 		i := j.i
// 		layer := j.layer
// 		box := j.box
// 		center := box.Center()
// 		size := box.Size()
// 		sx := (S - P*2) / size.X
// 		sy := (S - P*2) / size.Y
// 		scale := math.Min(sx, sy)
// 		dc := gg.NewContext(S, S)
// 		dc.InvertY()
// 		dc.SetRGB(1, 1, 1)
// 		dc.Clear()
// 		dc.Translate(S/2, S/2)
// 		dc.Scale(scale, scale)
// 		dc.Translate(-center.X, -center.Y)
// 		dc.SetFillRuleWinding()
// 		for _, path := range layer.Paths {
// 			dc.NewSubPath()
// 			for _, point := range path {
// 				dc.LineTo(point.X, point.Y)
// 			}
// 			dc.ClosePath()
// 		}
// 		dc.SetRGB(0, 0, 0)
// 		dc.Fill()
// 		// dc.SetRGB(0.6, 0.6, 0.6)
// 		// dc.FillPreserve()
// 		// dc.SetRGB(0, 0, 0)
// 		// dc.SetLineWidth(3)
// 		// dc.Stroke()
// 		dc.SavePNG(fmt.Sprintf("out%04d.png", i))
// 	}
// 	wg.Done()
// }
