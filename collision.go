package main

import (
	"math"
	"sort"

	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/cyberspace/lib/quadtree"
	"github.com/stojg/vector"
)

type byPenetration []*contact

func (a byPenetration) Len() int           { return len(a) }
func (a byPenetration) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPenetration) Less(i, j int) bool { return a[i].penetration < a[j].penetration }

// UpdateCollisions resolves collisions by using a impulse based collision resolution, ie it moves
// entities with impulses (instant velocity change) instead of using a force based collision
// resolution.
func UpdateCollisions(elapsed float64, frame uint64) {

	var potentialCollisions []*contact

	// now we can query the quad tree for each collidable
	for _, a := range core.List.Collisions(frame) {

		// checked is a list of what we already checked so we don't do duplicate checks
		checked := make(map[quadtree.BoundingBoxer]map[quadtree.BoundingBoxer]bool)

		// add only potential collisions that we haven't already checked
		for _, b := range core.List.QuadTree().Query(a.BoundingBox()) {
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
				restitution: 0.99, // hard coded restitution for now
				normal:      &vector.Vector3{},
				toWorld:     &vector.Matrix3{},
			}
			potentialCollisions = append(potentialCollisions, pair)
		}
	}

	// now it's time for the narrow phase of collision detection where we check if
	// objects are colliding and calculate the contact point and interpenetration

	var collisions []*contact
	for _, pair := range potentialCollisions {
		rectangleVsRectangle(pair, frame)
		if pair.IsIntersecting {
			collisions = append(collisions, pair)
		}
	}

	if len(collisions) == 0 {
		return
	}

	sort.Sort(byPenetration(collisions))

	// now it's time to resolve the collision
	for _, contact := range collisions {

		// add a bit leeway so that contacts are separating a bit more and don't constantly
		// get stuck in collision, kinda hacky TBH
		contact.penetration += 0.001

		contact.bodies[0] = contact.a.GameObject().Body()
		if !contact.bodies[0].CanCollide() {
			return
		}
		if contact.b != nil {
			contact.bodies[1] = contact.b.GameObject().Body()
			if !contact.bodies[1].CanCollide() {
				return
			}
		}

		// calculate the total inverse mass
		totalInvMass := contact.bodies[0].InvMass
		if contact.b != nil {
			totalInvMass += contact.bodies[1].InvMass
		}

		// if both object have infinate mass we cannot move them
		if totalInvMass == 0 {
			continue
		}

		// used later
		//contact.calculateInternals(elapsed)

		// 1. first fix the change in velocities that comes from things bouncing together

		// Find the velocity in the direction of the contact normal
		separatingVelocity := contact.separatingVelocity()

		// If the objects are already separating, we don't need to change the velocity
		if separatingVelocity < 0 {

			// Calculate the new separating velocity
			newSepVelocity := -separatingVelocity * contact.restitution

			// Check the velocity build up due to acceleration only
			accCausedVelocity := contact.bodies[0].Forces.Clone()
			if contact.b != nil {
				accCausedVelocity.Sub(contact.bodies[1].Forces)
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
			velocityChangeA := impulsePerIMass.NewScale(contact.bodies[0].InvMass)
			contact.bodies[0].Velocity().Add(velocityChangeA)
			contact.bodies[0].SetAwake(true)
			if contact.b != nil {
				velocityChangeB := impulsePerIMass.NewScale(-contact.bodies[1].InvMass)
				contact.bodies[1].Velocity().Add(velocityChangeB)
				contact.bodies[1].SetAwake(true)
			}
		}

		// 2. now it's time to resolve the interpenetration issue of the colliding objects
		if contact.penetration > 0 {
			movePerIMass := contact.normal.NewScale(contact.penetration / totalInvMass)
			contact.a.Transform().Position().Add(movePerIMass.NewScale(contact.bodies[0].InvMass))
			contact.bodies[0].SetAwake(true)
			if contact.b != nil {
				contact.b.Transform().Position().Add(movePerIMass.NewScale(-contact.bodies[1].InvMass))
				contact.bodies[1].SetAwake(true)
			}
		}
	}
}

func rectangleVsRectangle(contact *contact, frame uint64) {
	rA := contact.a.OBB(frame)
	rB := contact.b.OBB(frame)

	if rA == nil || rB == nil {
		return
	}

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

	// used later in the future https://github.com/idmillington/cyclone-physics/blob/fd0cf4956fd83ebf9e2e75421dfbf9f5cdac49fa/src/collide_fine.cpp#L422
	//toCentre := rB.CentrePoint().NewSub(rB.CentrePoint())
	//bestSingleAxis := mtvAxis.Clone()

	contact.penetration = mtvDistance
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
	IsIntersecting bool
	bodies         [2]*core.Body
	a              *core.Collision
	b              *core.Collision
	restitution    float64
	penetration    float64

	// Holds the position of the contact in world coordinates.
	contactPoint *vector.Vector3

	// the contact normal in world space
	normal *vector.Vector3

	// A transform matrix that converts coordinates in the contact's frame of reference to world
	// coordinates. The columns of this matrix form an orthonormal set of vectors
	toWorld *vector.Matrix3

	// Holds the closing velocity at the point of contact. This is set when the calculateInternals
	// function is run.
	contactVelocity *vector.Vector3

	// Holds the world space position of the contact point relative to centre of each body. This is
	// set when the calculateInternals function is run.
	relativeContactPosition [2]*vector.Vector3

	// Holds the required change in velocity for this contact to be
	// resolved.
	desiredDeltaVelocity float64
}

func (contact *contact) calculateInternals(duration float64) {
	// Check if the first object is nil, and swap if it is.
	if contact.bodies[0] == nil {
		contact.swapBodies()
	}

	// Calculate an set of axis at the contact point.
	contact.calculateContactBasis()

	// Store the relative position of the contact relative to each body
	contact.relativeContactPosition[0] = contact.contactPoint.NewSub(contact.a.Transform().Position())
	if contact.bodies[1] != nil {
		contact.relativeContactPosition[1] = contact.contactPoint.NewSub(contact.b.Transform().Position())
	}

	// Find the relative velocity of the bodies at the contact point.
	contact.contactVelocity = contact.calculateLocalVelocity(0, duration)
	if contact.bodies[1] != nil {
		contact.contactVelocity.Sub(contact.calculateLocalVelocity(1, duration))
	}

	// Calculate the desired change in velocity for resolution
	contact.calculateDesiredDeltaVelocity(duration)
}

func (contact *contact) swapBodies() {
	contact.bodies[0], contact.bodies[1] = contact.bodies[1], contact.bodies[0]
	contact.a, contact.b = contact.b, contact.a
}

// calculateContactBasis calculates an orthonormal basis for the contact point, based on the
// primary friction direction (for anisotropic friction) or a random orientation (for isotropic
// friction)
//
// Constructs an arbitrary orthonormal basis for the contact.  This is stored as a 3x3 matrix, where
// each vector is a column (in other words the matrix transforms contact space into world space).
// The x // direction is generated from the contact normal, and the y and z directionss are set so
// they are at right angles to it.
func (contact *contact) calculateContactBasis() {

	var contactTangent [2]*vector.Vector3
	contactTangent[0] = vector.NewVector3(0, 0, 0)
	contactTangent[1] = vector.NewVector3(0, 0, 0)

	// Check whether the Z-axis is nearer to the X or Y axis
	if math.Abs(contact.normal[0]) > math.Abs(contact.normal[1]) {
		// Scaling factor to ensure the results are normalised
		s := 1.0 / math.Sqrt(contact.normal[2]*contact.normal[2]+contact.normal[0]*contact.normal[0])

		// The new X-axis is at right angles to the world Y-axis
		contactTangent[0][0] = contact.normal[2] * s
		contactTangent[0][1] = 0
		contactTangent[0][2] = -contact.normal[0] * s

		// The new Y-axis is at right angles to the new X- and Z- axes
		contactTangent[1][0] = contact.normal[1] * contactTangent[0][0]
		contactTangent[1][1] = contact.normal[2]*contactTangent[0][0] - contact.normal[0]*contactTangent[0][2]
		contactTangent[1][2] = -contact.normal[1] * contactTangent[0][0]
	} else {
		// Scaling factor to ensure the results are normalised
		s := 1.0 / math.Sqrt(contact.normal[2]*contact.normal[2]+contact.normal[1]*contact.normal[1])

		// The new X-axis is at right angles to the world X-axis
		contactTangent[0][0] = 0
		contactTangent[0][1] = -contact.normal[2] * s
		contactTangent[0][2] = contact.normal[1] * s

		// The new Y-axis is at right angles to the new X- and Z- axes
		contactTangent[1][0] = contact.normal[1]*contactTangent[0][2] - contact.normal[2]*contactTangent[0][1]
		contactTangent[1][1] = -contact.normal[0] * contactTangent[0][2]
		contactTangent[1][2] = contact.normal[0] * contactTangent[0][1]
	}
	contact.toWorld.SetFromComponents(contact.normal, contactTangent[0], contactTangent[1])
}

func (contact *contact) separatingVelocity() float64 {
	relativeVel := contact.bodies[0].Velocity().Clone()
	if contact.b != nil {
		relativeVel.Sub(contact.bodies[1].Velocity())
	}
	return relativeVel.Dot(contact.normal)
}

func (contact *contact) calculateLocalVelocity(bodyIndex int, duration float64) *vector.Vector3 {
	thisBody := contact.bodies[bodyIndex]

	// Work out the velocity of the contact point.
	velocity := thisBody.Rotation().NewCross(contact.relativeContactPosition[bodyIndex])
	velocity.Add(thisBody.Velocity())

	// Turn the velocity into contact-coordinates.
	contactVelocity := contact.toWorld.TransformTranspose(velocity)

	// Calculate the amount of velocity that is due to forces without reactions.
	accVelocity := thisBody.LastFrameAcceleration.NewScale(duration)

	// Calculate the velocity in contact-coordinates.
	accVelocity = contact.toWorld.TransformTranspose(accVelocity)

	// We ignore any component of acceleration in the contact normal
	// direction, we are only interested in planar acceleration
	accVelocity[0] = 0

	// Add the planar velocities - if there's enough friction they will
	// be removed during velocity resolution
	contactVelocity.Add(accVelocity)

	// And return it
	return contactVelocity
}

func (contact *contact) calculateDesiredDeltaVelocity(duration float64) {
	const velocityLimit = 0.25

	// Calculate the acceleration induced velocity accumulated this frame
	velocityFromAcc := 0.0

	if contact.bodies[0].Awake() {
		velocityFromAcc += contact.bodies[0].LastFrameAcceleration.NewScale(duration).Dot(contact.normal)
	}

	if contact.bodies[1] != nil && contact.bodies[1].Awake() {
		velocityFromAcc -= contact.bodies[1].LastFrameAcceleration.NewScale(duration).Dot(contact.normal)
	}

	// If the velocity is very slow, limit the restitution
	thisRestitution := contact.restitution
	if math.Abs(contact.contactVelocity[0]) < velocityLimit {
		thisRestitution = 0.0
	}

	// Combine the bounce velocity with the removed
	// acceleration velocity.
	contact.desiredDeltaVelocity = -contact.contactVelocity[0] - thisRestitution*(contact.contactVelocity[0]-velocityFromAcc)
}

func (contact *contact) calculateFrictionlessImpulse(inverseInertiaTensor [2]*vector.Matrix3) *vector.Vector3 {

	var impulseContact *vector.Vector3

	// Build a vector that shows the change in velocity in world space for a unit impulse in the
	// direction of the contact normal.
	deltaVelWorld := contact.relativeContactPosition[0].NewCross(contact.normal)
	deltaVelWorld = inverseInertiaTensor[0].Transform(deltaVelWorld)
	deltaVelWorld = deltaVelWorld.NewCross(contact.relativeContactPosition[0])
	// Work out the change in velocity in contact coordinates
	deltaVelocity := deltaVelWorld.Dot(contact.normal)
	// Add the linear component of velocity change
	deltaVelocity += contact.bodies[0].InvMass

	// Check if we need to the second body's data
	if contact.bodies[1] != nil {
		// Go through the same transformation sequence again
		deltaVelWorld = contact.relativeContactPosition[1].NewCross(contact.normal)
		deltaVelWorld = inverseInertiaTensor[1].Transform(deltaVelWorld)
		deltaVelWorld = deltaVelWorld.NewCross(contact.relativeContactPosition[1])
		// Add the change in velocity due to rotation
		deltaVelocity += deltaVelWorld.Dot(contact.normal)
		// Add the change in velocity due to linear motion
		deltaVelocity += contact.bodies[1].InvMass
	}

	// Calculate the required size of the impulse
	impulseContact[0] = contact.desiredDeltaVelocity / deltaVelocity
	impulseContact[1] = 0
	impulseContact[2] = 0
	return impulseContact
}
