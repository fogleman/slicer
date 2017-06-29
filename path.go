package slicer

import "github.com/fogleman/fauxgl"

type Path []fauxgl.Vector

func joinPaths(paths []Path) []Path {
	lookup := make(map[fauxgl.Vector]Path, len(paths))
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
