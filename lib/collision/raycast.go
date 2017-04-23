package collision

import (
	"math"
	"sort"

	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/vector"
)

type RaycastResult struct {
	// A collision that the ray collided with
	Collision *core.Collision
	// this is the "percent" along the ray for the hit
	Distance float64
}

type Raylist []*RaycastResult

func (r Raylist) Len() int           { return len(r) }
func (r Raylist) Less(a, b int) bool { return r[a].Distance < r[b].Distance }
func (r Raylist) Swap(a, b int) {
	r[a], r[b] = r[b], r[a]
}

func Raycast(start, direction *vector.Vector3, list *core.ObjectList) Raylist {
	var result Raylist
	for _, col := range list.Collisions() {
		obb := col.OBB()
		// the collision doesn't have a body
		if obb == nil {
			panic("RayCast against an object without a Body")
		}
		// this is the 'percent' along the direction that the hit happened
		distance := 0.0
		if rayAABBoxIntersect(start, direction, col.OBB().MinPoint, col.OBB().MaxPoint, &distance) {
			result = append(result, &RaycastResult{
				Collision: col,
				Distance:  distance,
			})
		}
	}
	sort.Sort(result)
	return result
}

func rayAABBoxIntersect(start, direction, min, max *vector.Vector3, t *float64) bool {
	tfirst := 0.0
	tlast := 1.0

	if !raySlabIntersect(start[0], direction[0], min[0], max[0], &tfirst, &tlast) {
		return false
	}
	if !raySlabIntersect(start[1], direction[1], min[1], max[1], &tfirst, &tlast) {
		return false
	}
	if !raySlabIntersect(start[2], direction[2], min[2], max[2], &tfirst, &tlast) {
		return false
	}
	*t = tfirst
	return true
}

// returns the distance between ray_origin and the intersection with the OBB
func raySlabIntersect(start, dir, min, max float64, tfirst, tlast *float64) bool {
	if math.Abs(dir) < 1.0E-8 {
		return start < max && start > min
	}

	tmin := (min - start) / dir
	tmax := (max - start) / dir

	if tmin > tmax {
		tmin, tmax = tmax, tmin
	}

	if tmax < *tfirst || tmin > *tlast {
		return false
	}

	if tmin > *tfirst {
		*tfirst = tmin
	}

	if tmax < *tlast {
		*tlast = tmax
	}
	return true

}
