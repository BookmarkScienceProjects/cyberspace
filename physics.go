package main

import (
	"github.com/stojg/cyberspace/lib/core"
	"math"
)

type physicSystem struct{}

func UpdatePhysics(elapsed float64) {

	for _, body := range core.List.Bodies() {

		if !body.Awake() {
			continue
		}

		// Calculate linear acceleration from force inputs.
		body.LastFrameAcceleration = body.Acceleration.Clone()
		body.LastFrameAcceleration.AddScaledVector(body.ForceAccum, body.InvMass)

		// Calculate angular acceleration from torque inputs.
		angularAcceleration := body.InverseInertiaTensorWorld.Transform(body.TorqueAccum)

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
		body.Transform().Position().AddScaledVector(body.Velocity, elapsed)
		// Update angular position
		body.Transform().Orientation().AddScaledVector(body.Rotation, elapsed)

		// Normalise the orientation, and update the matrices with the new position and orientation
		body.CalculateDerivedData(body.Transform())

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
