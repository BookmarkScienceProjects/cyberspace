package main

import (
	. "github.com/stojg/vivere/lib/components"
	"github.com/stojg/vivere/lib/vector"
	"math/rand"
)

func NewAI(ent *Entity) *AI {
	rotation := &vector.Vector3{
		(rand.Float64() - 0.5),
		(rand.Float64() - 0.5),
		(rand.Float64() - 0.5),
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
}

func (s *AI) Update(elapsed float64) {
	body := rigidList.Get(s.entity)

	inst := monitor.FindByEntityID(*s.entity)
	if inst == nil {
		return
	}


	rotSpeed := body.Rotation.Length()
	if rotSpeed < inst.CPUUtilization/40 {
		body.AddTorque(s.spin)
	} else {
		body.SetAwake(true)
	}

}
