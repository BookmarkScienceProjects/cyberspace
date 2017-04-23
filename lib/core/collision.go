package core

import (
	"math"

	"github.com/stojg/cyberspace/lib/quadtree"
	"github.com/stojg/vector"
)

// NewCollisionRectangle returns a new Collision struct
func NewCollisionRectangle(x, y, z float64) *Collision {
	col := &Collision{
		fullWidth: [3]float64{x, y, z},
		halfWidth: [3]float64{x / 2, y / 2, z / 2},
	}
	col.cachedOBB = &OBB{
		//centre:   vector.Zero(),
		MaxPoint: vector.Zero(),
		MinPoint: vector.Zero(),
	}
	return col
}

// Collision is a struct that keeps track of the collision geometry for a GameObject
type Collision struct {
	Component
	fullWidth          [3]float64
	halfWidth          [3]float64
	cachedOBB          *OBB
	cachedFrameCounter uint64
}

// BoundingBox return the AABB bounding box the GameObject it is connected to
func (c *Collision) BoundingBox() quadtree.BoundingBox {
	return quadtree.BoundingBox{
		MinX: c.transform.position[0] - c.fullWidth[0],
		MaxX: c.transform.position[0] + c.fullWidth[0],
		MinY: c.transform.position[2] - c.fullWidth[2],
		MaxY: c.transform.position[2] + c.fullWidth[2],
	}
}

// OBB returns the Oriented Bounding Box for this volume
func (c *Collision) OBB(frame uint64) *OBB {
	// @todo cache with the frame counter this so it's not re-calculate for every SAT test

	if frame == c.cachedFrameCounter {
		return c.cachedOBB
	}
	c.cachedFrameCounter = frame

	mat := c.gameObject.Body().transformMatrix
	var points [8]*vector.Vector3
	vec := vector.NewVector3(c.halfWidth[0], c.halfWidth[1], c.halfWidth[2])
	points[0] = mat.TransformVector3(vec)
	vec.Set(c.halfWidth[0], -c.halfWidth[1], c.halfWidth[2])
	points[1] = mat.TransformVector3(vec)
	vec.Set(c.halfWidth[0], c.halfWidth[1], -c.halfWidth[2])
	points[2] = mat.TransformVector3(vec)
	vec.Set(c.halfWidth[0], -c.halfWidth[1], -c.halfWidth[2])
	points[3] = mat.TransformVector3(vec)
	vec.Set(-c.halfWidth[0], c.halfWidth[1], c.halfWidth[2])
	points[4] = mat.TransformVector3(vec)
	vec.Set(-c.halfWidth[0], -c.halfWidth[1], c.halfWidth[2])
	points[5] = mat.TransformVector3(vec)
	vec.Set(-c.halfWidth[0], c.halfWidth[1], -c.halfWidth[2])
	points[6] = mat.TransformVector3(vec)
	vec.Set(-c.halfWidth[0], -c.halfWidth[1], -c.halfWidth[2])
	points[7] = mat.TransformVector3(vec)
	c.cachedOBB.MaxPoint.Set(-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64)
	c.cachedOBB.MinPoint.Set(math.MaxFloat64, math.MaxFloat64, math.MaxFloat64)
	for j := 0; j < 3; j++ {
		for i := 0; i < 8; i++ {
			if points[i][j] > c.cachedOBB.MaxPoint[j] {
				c.cachedOBB.MaxPoint[j] = points[i][j]
			}
			if points[i][j] < c.cachedOBB.MinPoint[j] {
				c.cachedOBB.MinPoint[j] = points[i][j]
			}
		}
	}
	return c.cachedOBB
}

// OBB is a struct that represents the Oriented Bounded Box, i.e. a rotated AABB
type OBB struct {
	MinPoint *vector.Vector3
	MaxPoint *vector.Vector3
}
