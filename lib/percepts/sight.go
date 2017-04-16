package percepts

import (
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/vector"
	"math"
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
func InViewCone(me, other *core.GameObject, viewCone float64) float64 {
	direction := other.Transform().Position().NewSub(me.Transform().Position())
	// two entities in the same position can see each other
	if direction.SquareLength() == 0 {
		return 1
	}
	// normalise after the above check for efficiency
	direction.Normalize()

	guardFacing := vector.X().Rotate(me.Transform().Orientation())

	dot := guardFacing.Dot(direction)
	angle := math.Acos(dot)

	halfViewCone := viewCone / 2
	confidence := (halfViewCone - angle) / halfViewCone
	if confidence < 0 {
		return 0
	}
	return confidence
}
