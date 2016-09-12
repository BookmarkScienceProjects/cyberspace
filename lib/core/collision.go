package core

import (
	"github.com/stojg/vector"
	"github.com/stojg/cyberspace/lib/quadtree"
)

func NewCollisionRectangle(x, y, z float64) *Collision {
	return &Collision{
		fullWidth: [3]float64{x, y, z},
		halfWidth: [3]float64{x/2, y/2, z/2},
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
	r := &OBB{
		MinPoint: vector.Zero(),
		MaxPoint: vector.Zero(),
	}
	r.MinPoint[0] = c.transform.position[0] - c.halfWidth[0]
	r.MaxPoint[0] = c.transform.position[0] + c.halfWidth[0]
	r.MinPoint[1] = c.transform.position[1] - c.halfWidth[1]
	r.MaxPoint[1] = c.transform.position[1] + c.halfWidth[1]
	r.MinPoint[2] = c.transform.position[2] - c.halfWidth[2]
	r.MaxPoint[2] = c.transform.position[2] + c.halfWidth[2]
	return r
}

type OBB struct {
	HalfSize *vector.Vector3
	MinPoint *vector.Vector3
	MaxPoint *vector.Vector3
}
