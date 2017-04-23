package percepts

import (
	"math"

	"github.com/stojg/cyberspace/lib/collision"
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/vector"
)

// Distance checks distance between me and other with a maxDistance and returns the confidence between 0 and 1
func Distance(me, other *core.GameObject, maxDistance float64) float64 {
	direction := me.Transform().Position().NewSub(other.Transform().Position())
	if direction.SquareLength() > maxDistance*maxDistance {
		return 0
	}
	return (maxDistance - direction.Length()) / maxDistance
}

// InViewCone checks if other is in me view cone and returns a confidence between 0 and 1.
// viewCone should be in radians
func InSight(from *vector.Vector3, orientation *vector.Quaternion, to *vector.Vector3, viewCone float64) bool {
	direction := to.NewSub(from)
	// two entities in the same position can see each other
	if direction.SquareLength() == 0 {
		return false
	}
	// normalise after the above check for efficiency
	direction.Normalize()
	meFacing := vector.X().Rotate(orientation)

	dot := meFacing.Dot(direction)
	angle := math.Acos(dot)

	return angle < (viewCone / 2)
}

func InLineOfSight(from, to *vector.Vector3, frame uint64) collision.Raylist {
	direction := to.NewSub(from)
	return collision.Raycast(from, direction, core.List, frame)
}

// InViewCone checks if other is in me view cone and returns a confidence between 0 and 1.
// viewCone should be in radians
func InViewCone(me, other *core.GameObject, viewCone float64) float64 {
	direction := other.Transform().Position().NewSub(me.Transform().Position())
	// two entities in the same position can see each other
	if direction.SquareLength() == 0 {
		return 1
	}
	// normalise after the above check for efficiency
	direction.Normalize()

	meFacing := vector.X().Rotate(me.Transform().Orientation())

	dot := meFacing.Dot(direction)
	angle := math.Acos(dot)

	halfViewCone := viewCone / 2
	confidence := (halfViewCone - angle) / halfViewCone
	if confidence < 0 {
		return 0
	}
	return confidence
}

func CanSeeTarget(me, other *core.GameObject, frame uint64) bool {
	direction := other.Transform().Position().NewSub(me.Transform().Position())
	res := collision.Raycast(me.Transform().Position(), direction, core.List, frame)
	if len(res) == 0 {
		return false
	}

	for _, rr := range res {
		if rr.Collision == me.Collision() {
			continue
		}
		if rr.Collision != other.Collision() {
			return false
		}
		return true
	}
	return true
}
