package main

import (
	"math"
)

type physicSystem struct{}

func (s *physicSystem) Update(elapsed float64) {

	for i, body := range rigidList.All() {

		if !body.Awake() {
			continue
		}

		model := modelList.Get(i)
		if model == nil {
			panic("Physic system requires that a *Model have been set")
		}

		// Calculate linear acceleration from force inputs.
		body.LastFrameAcceleration = body.Acceleration.Clone()
		body.LastFrameAcceleration.AddScaledVector(body.ForceAccum, body.InvMass)

		// Calculate angular acceleration from torque inputs.
		angularAcceleration := body.InverseInertiaTensorWorld.TransformVector3(body.TorqueAccum)

		// Adjust velocities
		// Update linear velocity from both acceleration and impulse.
		body.Velocity.AddScaledVector(body.LastFrameAcceleration, elapsed)

		// Update angular velocity from both acceleration and impulse.
		body.Rotation.AddScaledVector(angularAcceleration, elapsed)

		// Impose drag
		body.Velocity.Scale(math.Pow(body.LinearDamping, elapsed))
		body.Rotation.Scale(math.Pow(body.AngularDamping, elapsed))

		// fake friction
		body.Velocity.Scale(0.95)
		body.Rotation.Scale(0.95)

		// Adjust positions
		// Update linear position
		model.Position().AddScaledVector(body.Velocity, elapsed)
		// Update angular position
		model.Orientation().AddScaledVector(body.Rotation, elapsed)

		// Normalise the orientation, and update the matrices with the new position and orientation
		body.CalculateDerivedData(model)

		// Clear accumulators.
		body.ClearAccumulators()

		// Update the kinetic energy store, and possibly put the body to sleep.
		if body.CanSleep {
			currentMotion := body.Velocity.ScalarProduct(body.Velocity) + body.Rotation.ScalarProduct(body.Rotation)
			bias := math.Pow(0.5, elapsed)
			motion := bias*body.Motion + (1-bias)*currentMotion
			if motion < body.SleepEpsilon {
				body.SetAwake(false)
			}
		} else if body.Motion > 10*body.SleepEpsilon {
			body.Motion = 10 * body.SleepEpsilon
		}
	}
}
