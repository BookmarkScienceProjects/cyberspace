package main

import (
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/cyberspace/lib/quadtree"
	"github.com/stojg/vector"
	"math"
)

func UpdateCollisions(elapsed float64) {
	var collisions []*contact

	tree := quadtree.NewQuadTree(
		quadtree.BoundingBox{
			MinX: -3200 / 2,
			MaxX: 3200 / 2,
			MinY: -3200 / 2,
			MaxY: 3200 / 2,
		},
	)

	// build the quad tree for broad phase collision detection
	for _, collision := range core.List.Collisions() {
		tree.Add(collision)
	}

	var potentialCollisions []*contact

	// now we can query the quad tree for each collidable
	for _, a := range core.List.Collisions() {

		// checked is a list of what we already checked so we don't do duplicate checks
		checked := make(map[quadtree.BoundingBoxer]map[quadtree.BoundingBoxer]bool)

		// add only potential collisions that we haven't already checked
		for _, b := range tree.Query(a.BoundingBox()) {
			// it cannot collide with it self
			if a == b {
				continue
			}
			// we already checked this combination a vs b
			if _, ok := checked[a][b]; ok {
				continue
			}
			// we already checked the "inverse" combination b vs a
			if _, ok := checked[b][a]; ok {
				continue
			}
			// we need to allocate space for the inner map [a][b] if it doesn't exists
			if _, ok := checked[a]; !ok {
				checked[a] = make(map[quadtree.BoundingBoxer]bool)
			}
			// we need to allocate space for the inner map [b][a] if it doesn't exists
			if _, ok := checked[b]; !ok {
				checked[b] = make(map[quadtree.BoundingBoxer]bool)
			}

			// mark this combination of collidables as checked
			checked[a][b], checked[b][a] = true, true

			// this the pair of contacts that we would like to check
			pair := &contact{
				a:           a,
				b:           b.(*core.Collision),
				restitution: 0.99, // hard coded restitution for nowx
				normal:      &vector.Vector3{},
			}
			potentialCollisions = append(potentialCollisions, pair)
		}
	}

	// now it's time for the narrow phase of collision detection where we check if
	// objects are colliding and calculate the contact point and interpenetration
	for _, pair := range potentialCollisions {
		rectangleVsRectangle(pair)
		if pair.IsIntersecting {
			collisions = append(collisions, pair)
		}
	}

	// now it's time to resolve the collision
	for _, contact := range collisions {

		contact.aBody = contact.a.GameObject().Body()
		if contact.b != nil {
			contact.bBody = contact.b.GameObject().Body()
		}

		// calculate the total inverse mass
		totalInvMass := contact.aBody.InvMass
		if contact.b != nil {
			totalInvMass += contact.bBody.InvMass
		}

		// if both object have infinate mass we cannot move them
		if totalInvMass == 0 {
			continue
		}

		// first fix the change in velocities that comes from things bouncing together

		// Find the velocity in the direction of the contact normal
		separatingVelocity := contact.separatingVelocity()

		// If the objects are already separating, we don't need to change the velocity
		if separatingVelocity < 0 {

			// Calculate the new separating velocity
			newSepVelocity := -separatingVelocity * contact.restitution

			// Check the velocity build up due to acceleration only
			accCausedVelocity := contact.aBody.Forces.Clone()
			if contact.b != nil {
				accCausedVelocity.Sub(contact.bBody.Forces)
			}

			// If we have closing velocity due to acceleration buildup,
			// remove it from the new separating velocity
			accCausedSepVelocity := accCausedVelocity.Dot(contact.normal) * elapsed
			if accCausedSepVelocity < 0 {
				newSepVelocity += contact.restitution * accCausedSepVelocity

				// make sure that we haven't removed more than was there to begin with
				if newSepVelocity < 0 {
					newSepVelocity = 0
				}
			}

			deltaVelocity := newSepVelocity - separatingVelocity

			// Only if either of the objects have mass that isn't infinite can we change the velocities
			impulsePerIMass := contact.normal.NewScale(deltaVelocity / totalInvMass)
			velocityChangeA := impulsePerIMass.NewScale(contact.aBody.InvMass)
			contact.aBody.Velocity.Add(velocityChangeA)

			if contact.b != nil {
				velocityChangeB := impulsePerIMass.NewScale(-contact.bBody.InvMass)
				contact.bBody.Velocity.Add(velocityChangeB)
			}

			contact.aBody = contact.a.GameObject().Body()
			if contact.b != nil {
				contact.bBody = contact.b.GameObject().Body()
			}
		}

		// now it's time to resolve the interpenetration issue of the colliding objects
		if contact.penetration > 0 {
			movePerIMass := contact.normal.NewScale(contact.penetration / totalInvMass)
			contact.a.Transform().Position().Add(movePerIMass.NewScale(contact.aBody.InvMass))
			if contact.b != nil {
				contact.b.Transform().Position().Add(movePerIMass.NewScale(-contact.bBody.InvMass))
			}
		}
	}
}

func rectangleVsRectangle(contact *contact) {
	rA := contact.a.OBB()
	rB := contact.b.OBB()

	// [Minimum Translation Vector]
	mtvDistance := math.MaxFloat32 // Set current minimum distance (max float value so next value is always less)
	mtvAxis := &vector.Vector3{}   // Axis along which to travel with the minimum distance

	// [Axes of potential separation]
	// [X Axis]
	if !testAxisSeparation(vector.UnitX, rA.MinPoint[0], rA.MaxPoint[0], rB.MinPoint[0], rB.MaxPoint[0], mtvAxis, &mtvDistance) {
		return
	}

	// [Y Axis]
	if !testAxisSeparation(vector.UnitY, rA.MinPoint[1], rA.MaxPoint[1], rB.MinPoint[1], rB.MaxPoint[1], mtvAxis, &mtvDistance) {
		return
	}

	// [Z Axis]
	if !testAxisSeparation(vector.UnitZ, rA.MinPoint[2], rA.MaxPoint[2], rB.MinPoint[2], rB.MaxPoint[2], mtvAxis, &mtvDistance) {
		return
	}

	contact.penetration = mtvDistance * 1.001
	contact.normal = mtvAxis.Normalize()
	contact.IsIntersecting = true
}

// TestAxisStatic checks if two axis overlaps and in that case calculates how much
// * Two convex shapes only overlap if they overlap on all axes of separation
// * In order to create accurate responses we need to find the
//    collision vector (Minimum Translation Vector)
// * Find if the two boxes intersect along a single axis
// * Compute the intersection interval for that axis
// * Keep the smallest intersection/penetration value
func testAxisSeparation(axis vector.Vector3, minA, maxA, minB, maxB float64, mtvAxis *vector.Vector3, mtvDistance *float64) bool {

	//	axisLengthSquared := axis.Dot(&axis)
	axisLengthSquared := axis[0]*axis[0] + axis[1]*axis[1] + axis[2]*axis[2]

	// If the axis is degenerate then ignore
	if axisLengthSquared < 1.0e-8 {
		return false
	}

	// Calculate the two possible overlap ranges
	// Either we overlap on the left or the right sides
	d0 := maxB - minA // 'Left' side
	d1 := maxA - minB // 'Right' side

	// Intervals do not overlap, so no intersection
	if d0 <= 0.0 || d1 <= 0.0 {
		return false
	}

	var overlap float64
	// Find out if we overlap on the 'right' or 'left' of the object.
	if d0 < d1 {
		overlap = d0
	} else {
		overlap = -d1
	}

	// The mtd vector for that axis
	var sep [3]float64
	sep[0] = axis[0] * (overlap / axisLengthSquared)
	sep[1] = axis[1] * (overlap / axisLengthSquared)
	sep[2] = axis[2] * (overlap / axisLengthSquared)

	// The mtd vector length squared
	sepLengthSquared := sep[0]*sep[0] + sep[1]*sep[1] + sep[2]*sep[2]

	// If that vector is smaller than our computed Minimum Translation
	// Distance use that vector as our current MTV distance
	if sepLengthSquared < *mtvDistance {
		*mtvDistance = math.Sqrt(sepLengthSquared) / 2
		mtvAxis.Set(sep[0], sep[1], sep[2])
	}
	return true
}

type contact struct {
	a              *core.Collision
	b              *core.Collision
	restitution    float64
	penetration    float64
	normal         *vector.Vector3
	IsIntersecting bool
	aBody          *core.Body
	bBody          *core.Body
}

func (contact *contact) separatingVelocity() float64 {
	relativeVel := contact.aBody.Velocity.Clone()
	if contact.b != nil {
		relativeVel.Sub(contact.bBody.Velocity)
	}
	return relativeVel.Dot(contact.normal)
}
