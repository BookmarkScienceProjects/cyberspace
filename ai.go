package main

import (
	. "github.com/stojg/vivere/lib/components"
	"github.com/stojg/vivere/lib/vector"
)

func NewAI(ent *Entity) *AI {
	rotation := &vector.Vector3{
		0, //rand.Float64() - 0.5,
		1,
		0, //rand.Float64() - 0.5,
	}
	//rotation.Normalize()
	ai := &AI{
		entity: ent,
		spin:   rotation,
	}
	return ai
}

type AI struct {
	entity *Entity
	spin   *vector.Vector3
}

func (s *AI) Update(elapsed float64) {
	body := rigidList.Get(s.entity)
	body.SetAwake(true)

	inst := monitor.FindByEntityID(*s.entity)
	if inst == nil {
		return
	}

	rotSpeed := body.Rotation.Length()
	if rotSpeed < inst.CPUUtilization/80 {
		body.AddTorque(s.spin)
	}

}
