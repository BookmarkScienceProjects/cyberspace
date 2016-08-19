package main

import (
	"github.com/stojg/vivere/lib/components"
	. "github.com/stojg/vivere/lib/components"
	"github.com/stojg/vivere/lib/vector"
	"math/rand"
)

func NewAI(ent *Entity) *AI {
	rotation := &vector.Vector3{
		0,
		rand.Float64() - 0.5,
		0,
	}
	rotation.Normalize()
	ai := &AI{
		entity: ent,
		spin:   rotation,
	}
	return ai
}

type AI struct {
	entity *Entity
	spin   *vector.Vector3
	target *vector.Vector3
}

func (s *AI) Update(elapsed float64) {
	body := rigidList.Get(s.entity)
	body.SetAwake(true)

	inst := monitor.FindByEntityID(*s.entity)
	if inst == nil {
		return
	}

	ids := inst.tree.Siblings(inst.Name)

	model := modelList.Get(s.entity)
	model.Position[1] = model.Scale[1] / 2

	if len(ids) > 1 {
		midpoint := vector.NewVector3(0, 0, 0)
		for i := range ids {
			model := modelList.Get(ids[i])
			midpoint.Add(model.Position)
		}
		midpoint.Scale(1 / float64(len(ids)))

		steering := s.getSteering(model, midpoint, body)
		body.AddForce(steering.linear)

	}

	rotSpeed := body.Rotation.Length()
	if rotSpeed < inst.CPUUtilization/80 {
		body.AddTorque(s.spin)
	}

}

func (s *AI) getSteering(model *components.Model, midpoint *vector.Vector3, body *components.RigidBody) *SteeringOutput {

	targetRadius := model.Scale[1] * 2
	slowRadius := targetRadius * 5
	maxSpeed := 150.0
	timeToTarget := 0.1

	// Get a new steering output
	steering := NewSteeringOutput()
	// Get the direction to the target
	direction := midpoint.NewSub(model.Position)
	distance := direction.Length()
	// We have arrived, no output
	if distance < targetRadius {
		return steering
	}

	// We are outside the slow radius, so full speed ahead
	var targetSpeed float64
	if distance > slowRadius {
		targetSpeed = maxSpeed
	} else {
		targetSpeed = maxSpeed * distance / slowRadius
	}

	// The target velocity combines speed and direction
	targetVelocity := direction
	targetVelocity.Normalize()
	targetVelocity.Scale(targetSpeed)
	// Acceleration tries to get to the target velocity
	steering.linear = targetVelocity.NewSub(body.Velocity)
	steering.linear.Scale(1 / timeToTarget)

	return steering
}
