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

func (p Path) Chop(step float64) Path {
	var result Path
	for i := 0; i < len(p)-1; i++ {
		a := p[i]
		b := p[i+1]
		v := b.Sub(a)
		l := v.Length()
		if i == 0 {
			result = append(result, a)
		}
		d := step
		for d < l {
			result = append(result, a.Add(v.MulScalar(d/l)))
			d += step
		}
		result = append(result, b)
	}
	return result
}
