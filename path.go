package slicer

import "github.com/fogleman/fauxgl"

type Path []fauxgl.Vector

// func (path Path) ToGEOS() *geos.Geometry {
// 	coords := make([]geos.Coord, len(path))
// 	for i, p := range path {
// 		coords[i] = geos.Coord{p.X, p.Y, 0}
// 	}
// 	return geos.Must(geos.NewPolygon(coords))
// }

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
