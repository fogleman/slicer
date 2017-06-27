package slicer

import (
	"github.com/fogleman/fauxgl"
	"github.com/paulsmith/gogeos/geos"
)

func pathToGeos(path Path) *geos.Geometry {
	coords := make([]geos.Coord, len(path))
	for i, p := range path {
		coords[i] = geos.Coord{p.X, p.Y, 0}
	}
	return geos.Must(geos.NewLinearRing(coords...))
}

func geosToPath(g *geos.Geometry) Path {
	coords := geos.MustCoords(g.Coords())
	path := make(Path, len(coords))
	for i, c := range coords {
		path[i] = fauxgl.Vector{c.X, c.Y, 0}
	}
	// fmt.Println(path)
	// path.reverse()
	return path
}

func pathsToGeos(paths []Path) *geos.Geometry {
	n := len(paths)
	rings := make([]*geos.Geometry, n)
	polys := make([]*geos.Geometry, n)
	hole := make([]bool, n)
	for i, p := range paths {
		rings[i] = pathToGeos(p)
		polys[i] = geos.Must(geos.PolygonFromGeom(rings[i]))
		hole[i] = p.isClockwise()
	}
	lookup := make(map[int][]int)
	// for each hole, find the smallest non-hole that contains it
	for i := 0; i < n; i++ {
		if !hole[i] {
			continue
		}
		index := -1
		var best float64
		for j := 0; j < n; j++ {
			if hole[j] {
				continue
			}
			contains, _ := polys[j].Contains(polys[i])
			if !contains {
				continue
			}
			area, _ := polys[j].Area()
			if index < 0 || area < best {
				best = area
				index = j
			}
		}
		lookup[index] = append(lookup[index], i)
	}
	// create polygons
	var polygons []*geos.Geometry
	for i, ring := range rings {
		if hole[i] {
			continue
		}
		var holes [][]geos.Coord
		for _, j := range lookup[i] {
			coords := geos.MustCoords(rings[j].Coords())
			holes = append(holes, coords)
		}
		coords := geos.MustCoords(ring.Coords())
		p := geos.Must(geos.NewPolygon(coords, holes...))
		polygons = append(polygons, p)
	}
	return geos.Must(geos.NewCollection(geos.MULTIPOLYGON, polygons...))
}

func geosToPaths(g *geos.Geometry) []Path {
	t, _ := g.Type()
	var paths []Path
	switch t {
	case geos.POLYGON:
		paths = append(paths, geosToPath(geos.Must(g.Shell())))
		holes, _ := g.Holes()
		for _, hole := range holes {
			paths = append(paths, geosToPath(hole))
		}
	case geos.MULTIPOLYGON:
		n, _ := g.NGeometry()
		for i := 0; i < n; i++ {
			paths = append(paths, geosToPaths(geos.Must(g.Geometry(i)))...)
		}
	default:
		panic(t)
	}
	return paths
}
