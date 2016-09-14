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

	volume := vector.NewVector3(c.halfWidth[0]*2, c.halfWidth[1]*2, c.halfWidth[2]*2)
	volume.Add(c.transform.position)
	//u := vector.NewVector3(c.transform.orientation.I, c.transform.orientation.J, c.transform.orientation.K)
	//s := c.transform.orientation.R;
	//a := u.Scale(2.0*u.Dot(volume))
	//b := volume.Scale(s*s - u.Dot(u))
	//d := u.NewCross(volume).Scale(2.0 * s)
	//z := a.Add(b).Add(d)

	z := volume.Rotate(c.transform.orientation)

	return &OBB{
		MinPoint: c.transform.position.NewSub(z),
		MaxPoint: c.transform.position.NewAdd(z),
	}
}

type OBB struct {
	HalfSize *vector.Vector3
	MinPoint *vector.Vector3
	MaxPoint *vector.Vector3
}
