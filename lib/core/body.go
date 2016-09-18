package core

import (
	"github.com/stojg/vector"
	"sync"
)

// NewBody returns a new rigidbody that is primarily used for simulating physics
func NewBody(invMass float64) *Body {
	body := &Body{
		velocity:                  &vector.Vector3{},
		rotation:                  &vector.Vector3{},
		Forces:                    &vector.Vector3{},
		transformMatrix:           &vector.Matrix4{},
		InverseInertiaTensor:      &vector.Matrix3{},
		InverseInertiaTensorWorld: &vector.Matrix3{},
		ForceAccum:                &vector.Vector3{},
		TorqueAccum:               &vector.Vector3{},
		maxAcceleration:           &vector.Vector3{100, 100, 100},
		Acceleration:              &vector.Vector3{},
		LinearDamping:             0.99,
		AngularDamping:            0.99,
		maxRotation:               3.14 / 1,
		InvMass:                   invMass,
		CanSleep:                  true,
		isAwake:                   true,
		SleepEpsilon:              0.00001,
	}

	inertiaTensor := &vector.Matrix3{}
	// hard coded cube inertia tensor for now
	inertiaTensor.SetBlockInertiaTensor(&vector.Vector3{1, 1, 1}, 1/invMass)
	body.setInertiaTensor(inertiaTensor)
	return body
}

// Body is a struct that contains data for doing rigid body physics simulation
type Body struct {
	Component
	sync.Mutex
	// Holds the linear velocity of the rigid body in world space.
	velocity *vector.Vector3
	// Holds the angular velocity, or rotation for the rigid body in world space.
	rotation *vector.Vector3

	// Holds the inverse of the mass of the rigid body. It is more useful to hold the inverse mass
	// because integration is simpler, and because in real time simulation it is more useful to have
	// bodies with infinite mass (immovable) than zero mass (completely unstable in numerical
	// simulation).
	InvMass float64
	// Holds the inverse of the body's inertia tensor. The inertia tensor provided must not be
	// degenerate (that would mean the body had zero inertia for spinning along one axis). As long
	// as the tensor is finite, it will be invertible. The inverse tensor is used for similar
	// reasons to the use of inverse mass.
	// The inertia tensor, unlike the other variables that define a rigid body, is given in body
	// space.
	InverseInertiaTensor *vector.Matrix3
	// Holds the amount of damping applied to linear motion.  Damping is required to remove energy
	// added through numerical instability in the integrator.
	LinearDamping float64
	// Holds the amount of damping applied to angular motion.  Damping is required to remove energy
	// added through numerical instability in the integrator.
	AngularDamping float64

	/**
	 * Derived Data
	 *
	 * These data members hold information that is derived from the other data in the class.
	 */

	// Holds the inverse inertia tensor of the body in world space. The inverse inertia tensor
	// member is specified in the body's local space. @see inverseInertiaTensor
	InverseInertiaTensorWorld *vector.Matrix3
	// Holds the amount of motion of the body. This is a recency weighted mean that can be used to
	// put a body to sleap.
	Motion float64
	// A body can be put to sleep to avoid it being updated by the integration functions or affected
	// by collisions with the world.
	isAwake bool
	// Some bodies may never be allowed to fall asleep. User controlled bodies, for example, should
	// be always awake.
	CanSleep bool
	// Holds a transform matrix for converting body space into world space and vice versa. This can
	// be achieved by calling the getPointIn*Space functions.
	transformMatrix *vector.Matrix4

	/**
	 * Force and Torque Accumulators
	 *
	 * These data members store the current force, torque and acceleration of the rigid body. Forces
	 * can be added to the rigid body in any order, and the class decomposes them into their
	 * constituents, accumulating them for the next simulation step. At the simulation step, the
	 * accelerations are calculated and stored to be applied to the rigid body.
	 */

	// Holds the accumulated force to be applied at the next integration step.
	ForceAccum *vector.Vector3

	// Holds the accumulated torque to be applied at the next integration step.
	TorqueAccum *vector.Vector3

	// Holds the acceleration of the rigid body.  This value can be used to set acceleration due to
	// gravity (its primary use), or any other constant acceleration.
	Acceleration *vector.Vector3

	maxAcceleration *vector.Vector3

	// limits the linear acceleration
	MaxAngularAcceleration *vector.Vector3
	// limits the angular velocity
	maxRotation float64

	// Holds the linear acceleration of the rigid body, for the previous frame.
	LastFrameAcceleration *vector.Vector3

	SleepEpsilon float64
	Forces       *vector.Vector3
}

// Position returns the position of the transform for this body
func (rb *Body) Position() *vector.Vector3 {
	return rb.transform.position
}

// Orientation returns the current orientation of the transform for this body
func (rb *Body) Orientation() *vector.Quaternion {
	return rb.transform.orientation
}

// Rotation returns the current rotation (the angular version of velocity) for this body
func (rb *Body) Rotation() *vector.Vector3 {
	return rb.rotation
}

// MaxRotation returns the maximum rotation speed for this body in radians
func (rb *Body) MaxRotation() float64 {
	return rb.maxRotation
}

// MaxAcceleration returns the maximum acceleration this body that do
func (rb *Body) MaxAcceleration() *vector.Vector3 {
	return rb.maxAcceleration
}

// Velocity returns the current velocity for this body
func (rb *Body) Velocity() *vector.Vector3 {
	return rb.velocity
}

// Mass returns the mass for this body
func (rb *Body) Mass() float64 {
	return 1 / rb.InvMass
}

// AddForce adds a force vector that originating from the center of this body
func (rb *Body) AddForce(force *vector.Vector3) {
	rb.ForceAccum.Add(force)
	rb.SetAwake(true)
}

// AddForceAtBodyPoint adds a force that is located at a specific point on this body. The point
// should be relative to the body, not world space
func (rb *Body) AddForceAtBodyPoint(transform *Transform, force, point *vector.Vector3) {
	// convert to coordinates relative to center of mass
	pt := rb.PointInWorldSpace(point)
	rb.AddForceAtPoint(transform, force, pt)
	rb.SetAwake(true)
}

// AddForceAtPoint adds a force that is located at a specific point on this body. The point
// should be in world space.
func (rb *Body) AddForceAtPoint(body *Transform, force, point *vector.Vector3) {
	// convert to coordinates relative to center of mass
	pt := point.NewSub(body.position)
	rb.ForceAccum.Add(force)
	rb.TorqueAccum.Add(pt.NewCross(force))
	rb.SetAwake(true)
}

// AddTorque as a torque (spin force) that has it's centre at the body
func (rb *Body) AddTorque(torque *vector.Vector3) {
	rb.TorqueAccum.Add(torque)
	rb.SetAwake(true)
}

// ClearAccumulators will clear all currently added forces and torque. Should be called after the
// physics integration is done
func (rb *Body) ClearAccumulators() {
	rb.Forces.Clear()
	rb.ForceAccum.Clear()
	rb.TorqueAccum.Clear()
}

// CalculateDerivedData recalculates the bodies internal data such as transformation matrices and
// inertia tensors, should be called once per physics simulation step to keep the data up to date.
func (rb *Body) CalculateDerivedData(transform *Transform) {
	transform.Orientation().Normalize()
	rb.calculateTransformMatrix(rb.transformMatrix, transform.position, transform.Orientation())
	rb.transformInertiaTensor(rb.InverseInertiaTensorWorld, transform.Orientation(), rb.InverseInertiaTensor, rb.transformMatrix)
}

// PointInWorldSpace returns the point of the body in world space coordinates
func (rb *Body) PointInWorldSpace(bodyRelativePoint *vector.Vector3) *vector.Vector3 {
	return rb.transformMatrix.TransformVector3(bodyRelativePoint)
}

// Awake reports if this body is current awake or asleep. If it's asleep it means that it's not
// moving and it can be skipped during physics simulation to decrease cpu usage
func (rb *Body) Awake() bool {
	rb.Lock()
	defer rb.Unlock()
	return rb.isAwake
}

// SetAwake sets this body in an awake or asleep state.
func (rb *Body) SetAwake(t bool) {
	rb.Lock()
	defer rb.Unlock()
	rb.isAwake = t
}

// setInertiaTensor is a utility function to set the inertia tensor correct
func (rb *Body) setInertiaTensor(inertiaTensor *vector.Matrix3) {
	rb.InverseInertiaTensor.SetInverse(inertiaTensor)
}

// transformInertiaTensor does an inertia tensor transform by a vector.Quaternion. The result of
// this transform will be set on the iitWorld matrix
func (rb *Body) transformInertiaTensor(iitWorld *vector.Matrix3, q *vector.Quaternion, iitBody *vector.Matrix3, rotMat *vector.Matrix4) {
	t4 := rotMat[0]*iitBody[0] + rotMat[1]*iitBody[3] + rotMat[2]*iitBody[6]
	t9 := rotMat[0]*iitBody[1] + rotMat[1]*iitBody[4] + rotMat[2]*iitBody[7]
	t14 := rotMat[0]*iitBody[2] + rotMat[1]*iitBody[5] + rotMat[2]*iitBody[8]
	t28 := rotMat[4]*iitBody[0] + rotMat[5]*iitBody[3] + rotMat[6]*iitBody[6]
	t33 := rotMat[4]*iitBody[1] + rotMat[5]*iitBody[4] + rotMat[6]*iitBody[7]
	t38 := rotMat[4]*iitBody[2] + rotMat[5]*iitBody[5] + rotMat[6]*iitBody[8]
	t52 := rotMat[8]*iitBody[0] + rotMat[9]*iitBody[3] + rotMat[10]*iitBody[6]
	t57 := rotMat[8]*iitBody[1] + rotMat[9]*iitBody[4] + rotMat[10]*iitBody[7]
	t62 := rotMat[8]*iitBody[2] + rotMat[9]*iitBody[5] + rotMat[10]*iitBody[8]

	iitWorld[0] = t4*rotMat[0] + t9*rotMat[1] + t14*rotMat[2]
	iitWorld[1] = t4*rotMat[4] + t9*rotMat[5] + t14*rotMat[6]
	iitWorld[2] = t4*rotMat[8] + t9*rotMat[9] + t14*rotMat[10]
	iitWorld[3] = t28*rotMat[0] + t33*rotMat[1] + t38*rotMat[2]
	iitWorld[4] = t28*rotMat[4] + t33*rotMat[5] + t38*rotMat[6]
	iitWorld[5] = t28*rotMat[8] + t33*rotMat[9] + t38*rotMat[10]
	iitWorld[6] = t52*rotMat[0] + t57*rotMat[1] + t62*rotMat[2]
	iitWorld[7] = t52*rotMat[4] + t57*rotMat[5] + t62*rotMat[6]
	iitWorld[8] = t52*rotMat[8] + t57*rotMat[9] + t62*rotMat[10]
}

// calculateTransformMatrix creates a transform matrix from a position and orientation. The result
// is be set on the rotMat
func (rb *Body) calculateTransformMatrix(rotMat *vector.Matrix4, pos *vector.Vector3, orientation *vector.Quaternion) {

	rotMat[0] = 1 - 2*orientation.J*orientation.J - 2*orientation.K*orientation.K
	rotMat[1] = 2*orientation.I*orientation.J - 2*orientation.R*orientation.K
	rotMat[2] = 2*orientation.I*orientation.K + 2*orientation.R*orientation.J
	rotMat[3] = pos[0]

	rotMat[4] = 2*orientation.I*orientation.J + 2*orientation.R*orientation.K
	rotMat[5] = 1 - 2*orientation.I*orientation.I - 2*orientation.K*orientation.K
	rotMat[6] = 2*orientation.J*orientation.K - 2*orientation.R*orientation.I
	rotMat[7] = pos[1]

	rotMat[8] = 2*orientation.I*orientation.K - 2*orientation.R*orientation.J
	rotMat[9] = 2*orientation.J*orientation.K + 2*orientation.R*orientation.I
	rotMat[10] = 1 - 2*orientation.I*orientation.I - 2*orientation.J*orientation.J
	rotMat[11] = pos[2]
}
