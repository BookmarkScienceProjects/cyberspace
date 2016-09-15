package core

import (
	"github.com/stojg/cyberspace/lib/quadtree"
	"github.com/stojg/vector"
	"math"
)

func NewCollisionRectangle(x, y, z float64) *Collision {
	return &Collision{
		fullWidth: [3]float64{x, y, z},
		halfWidth: [3]float64{x / 2, y / 2, z / 2},
	}
}

type Collision struct {
	Component
	fullWidth [3]float64
	halfWidth [3]float64
}

func (c *Collision) BoundingBox() quadtree.BoundingBox {
	return quadtree.BoundingBox{
		MinX: c.transform.position[0] - c.fullWidth[0],
		MaxX: c.transform.position[0] + c.fullWidth[0],
		MinY: c.transform.position[2] - c.fullWidth[2],
		MaxY: c.transform.position[2] + c.fullWidth[2],
	}
}

// OBB returns the Oriented Bounding Box for this volume
func (c *Collision) OBB() *OBB {
	// @todo cache this so it's not re-calculate for every SAT test
	halfSize := c.halfWidth[0]

	mat := c.gameObject.Body().transformMatrix

	var points [8]*vector.Vector3
	points[0] = mat.TransformVector3(vector.NewVector3(halfSize, halfSize, halfSize))
	points[0] = mat.TransformVector3(vector.NewVector3(halfSize, halfSize, halfSize))
	points[1] = mat.TransformVector3(vector.NewVector3(halfSize, -halfSize, halfSize))
	points[2] = mat.TransformVector3(vector.NewVector3(halfSize, halfSize, -halfSize))
	points[3] = mat.TransformVector3(vector.NewVector3(halfSize, -halfSize, -halfSize))

	points[4] = mat.TransformVector3(vector.NewVector3(-halfSize, halfSize, halfSize))
	points[5] = mat.TransformVector3(vector.NewVector3(-halfSize, -halfSize, halfSize))
	points[6] = mat.TransformVector3(vector.NewVector3(-halfSize, halfSize, -halfSize))
	points[7] = mat.TransformVector3(vector.NewVector3(-halfSize, -halfSize, -halfSize))

	max := vector.NewVector3(-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64)
	min := vector.NewVector3(math.MaxFloat64, math.MaxFloat64, math.MaxFloat64)
	for i := range points {
		for j := range max {
			if points[i][j] > max[j] {
				max[j] = points[i][j]
			}
			if points[i][j] < min[j] {
				min[j] = points[i][j]
			}
		}
	}
	return &OBB{
		centre:   c.gameObject.Transform().Position(),
		MaxPoint: max,
		MinPoint: min,
	}
}

type OBB struct {
	centre *vector.Vector3

	HalfSize *vector.Vector3

	MinPoint *vector.Vector3
	MaxPoint *vector.Vector3
}

func (obb *OBB) CentrePoint() *vector.Vector3 {
	return obb.centre
}
