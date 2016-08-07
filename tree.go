package main

import (
	"github.com/stojg/vivere/lib/vector"
	"math"
	"math/rand"
)

func BuildTree(c Collidable) {
	for _, c := range c.Children() {
		BuildTree(c)
	}

	// 1. place all instances (leaf nodes)
	for _, leaf := range c.Leaves() {
		leaf.SetPosition(vector.NewVector3(
			(rand.Float64()-0.5)*0,
			leaf.Instance().Scale[1]/2,
			(rand.Float64()-0.5)*0,
		))
	}

	//2. resolve all collisions among leaf nodes
	if len(c.Leaves()) > 1 {
		num := 1
		for num != 0 {
			//Println("Resolve leaf nodes", c.Name())
			num = IntersectionResolving(c.Leaves())
		}
	}

	for _, leaf := range c.Children() {
		//Println("Placing", leaf.Name())
		leaf.AddPosition(vector.NewVector3(
			(rand.Float64()-0.5)*300,
			//leaf.Instance().Scale[1]/2,
			0,
			(rand.Float64()-0.5)*300,
		))
	}

	if len(c.Children()) > 1 {
		num := 1
		for num != 0 {
			num = IntersectionResolving(c.Children())
		}
	}
}

func SetPositionOfInstances(c Collidable) {
	for _, c := range c.Children() {
		SetPositionOfInstances(c)
	}

	for _, leaf := range c.Leaves() {
		i := leaf.Instance()
		if i != nil {
			body := modelList.Get(i.ID)
			body.Position = leaf.Position()
			body.Position[1] = body.Scale[1] / 2
		}
	}
}

func IntersectionResolving(leaves []Collidable) int {

	checked := make(map[int]map[int]bool)
	collisions := make([]*rectCol, 0)

	for i, a := range leaves {
		for j, b := range leaves {
			if i == j {
				continue
			}
			if _, ok := checked[i][j]; ok {
				continue
			}
			if _, ok := checked[j][i]; ok {
				continue
			}
			if _, ok := checked[i]; !ok {
				checked[i] = make(map[int]bool)
			}

			if _, ok := checked[j]; !ok {
				checked[j] = make(map[int]bool)
			}
			checked[i][j], checked[j][i] = true, true

			pair := &rectCol{
				A: a,
				B: b,
			}

			pair.intersects()
			//pair.circleVsCircle()

			if pair.IsIntersecting {
				collisions = append(collisions, pair)
			}
		}
	}

	num := len(collisions)
	var biggest *rectCol
	pen := -math.MaxFloat64
	for _, pair := range collisions {
		if pair.penetration > pen {
			pen = pair.penetration
			biggest = pair
		}

	}
	if pen > 0 {
		biggest.Resolve()
	}

	return num
}

type rectCol struct {
	A, B           Collidable
	penetration    float64
	normal         *vector.Vector3
	IsIntersecting bool
}

func (contact *rectCol) circleVsCircle() {

	var d [3]float64
	for i := range d {
		d[i] = (contact.A.Position()[i] - contact.B.Position()[i])
	}

	sqrLength := d[0]*d[0] + d[1]*d[1] + d[2]*d[2]
	if sqrLength < 1 {
		return
	}

	// Early out to avoid expensive sqrt
	if sqrLength > (contact.A.Radius()+contact.B.Radius())*(contact.A.Radius()+contact.B.Radius()) {
		//if sqrLength > (cA.Radius+cB.Radius)*(cA.Radius+cB.Radius) {
		return
	}

	length := math.Sqrt(sqrLength)

	for i := range d {
		d[i] *= 1 / length
	}

	contact.penetration = contact.A.Radius() + contact.B.Radius() - length + 1
	contact.normal = &vector.Vector3{d[0], d[1], d[2]}
	contact.IsIntersecting = true
}

func (contact *rectCol) intersects() {
	mtvDistance := math.MaxFloat32 // Set current minimum distance (max float value so next value is always less)
	mtvAxis := &vector.Vector3{}   // Axis along which to travel with the minimum distance

	// [Axes of potential separation]
	// [X Axis]
	if !testAxisSeparation(vector.UnitX, contact.A.MinPoint(0), contact.A.MaxPoint(0), contact.B.MinPoint(0), contact.B.MaxPoint(0), mtvAxis, &mtvDistance) {
		return
	}

	// [Y Axis]
	if !testAxisSeparation(vector.UnitY, contact.A.MinPoint(1), contact.A.MaxPoint(1), contact.B.MinPoint(1), contact.B.MaxPoint(1), mtvAxis, &mtvDistance) {
		return
	}

	// [Z Axis]
	if !testAxisSeparation(vector.UnitZ, contact.A.MinPoint(2), contact.A.MaxPoint(2), contact.B.MinPoint(2), contact.B.MaxPoint(2), mtvAxis, &mtvDistance) {
		return
	}

	contact.penetration = mtvDistance + 1
	contact.normal = mtvAxis.Normalize()
	contact.IsIntersecting = true
}

func (contact *rectCol) Resolve() {
	if contact.penetration <= 0 {
		return
	}
	movePerIMass := contact.normal.NewScale(contact.penetration / 2)
	contact.A.AddPosition(movePerIMass.NewScale(1))
	contact.B.AddPosition(movePerIMass.NewScale(-1))
}
