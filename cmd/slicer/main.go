package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/fogleman/fauxgl"
	"github.com/fogleman/gg"
	"github.com/fogleman/slicer"
)

var (
	quiet     = kingpin.Flag("quiet", "Run in silent mode.").Short('q').Bool()
	directory = kingpin.Flag("directory", "Set the output directory.").Short('d').ExistingDir()
	size      = kingpin.Flag("size", "Set the slice thickness.").Required().Short('s').Float()
	width     = kingpin.Flag("width", "Set the raster width in model units.").Required().Short('w').Float()
	height    = kingpin.Flag("height", "Set the raster height in model units.").Required().Short('h').Float()
	scale     = kingpin.Flag("scale", "Set the raster scale.").Short('x').Required().Float()
	files     = kingpin.Arg("files", "Mesh files to slice.").Required().ExistingFiles()
)

func main() {
	kingpin.Parse()
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

	// slice mesh
	done = timed("slicing mesh")
	layers := slicer.SliceMesh(mesh, *size)
	done()

	// determine output dir
	dir := "."
	if *directory != "" {
		dir = *directory
	}
	_, name := filepath.Split(infile)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	dir = filepath.Join(dir, name)
	os.MkdirAll(dir, os.ModePerm)

	// render pngs
	render(box, layers, dir)
}

type job struct {
	i     int
	layer slicer.Layer
	box   fauxgl.Box
}

func render(box fauxgl.Box, layers []slicer.Layer, dir string) {
	wn := runtime.NumCPU()
	ch := make(chan job, len(layers))
	var wg sync.WaitGroup
	for wi := 0; wi < wn; wi++ {
		wg.Add(1)
		go worker(ch, &wg, dir)
	}
	for i, l := range layers {
		ch <- job{i, l, box}
	}
	close(ch)
	wg.Wait()
}

func fastRGBAToGray(src *image.RGBA) *image.Gray {
	dst := image.NewGray(src.Bounds())
	w := src.Bounds().Size().X
	h := src.Bounds().Size().Y
	for y := 0; y < h; y++ {
		i := src.PixOffset(0, y)
		j := dst.PixOffset(0, y)
		for x := 0; x < w; x++ {
			dst.Pix[j] = src.Pix[i]
			i += 4
			j++
		}
	}
	return dst
}

func savePNG(path string, im image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := png.Encoder{
		CompressionLevel: png.BestSpeed,
	}
	return encoder.Encode(file, im)
}

func worker(ch chan job, wg *sync.WaitGroup, dir string) {
	s := *scale
	w := *width * s
	h := *height * s
	for j := range ch {
		i := j.i
		layer := j.layer
		box := j.box
		center := box.Center()
		dc := gg.NewContext(int(w), int(h))
		dc.InvertY()
		dc.SetRGB(1, 1, 1)
		dc.Clear()
		dc.Translate(w/2, h/2)
		dc.Scale(s, s)
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
		src := dc.Image()
		dst := fastRGBAToGray(src.(*image.RGBA))
		path, _ := filepath.Abs(filepath.Join(dir, fmt.Sprintf("%04d.png", i)))
		savePNG(path, dst)
		fmt.Println(path)
	}
	wg.Done()
}
