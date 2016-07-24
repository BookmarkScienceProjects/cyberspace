package main

import (
	. "github.com/stojg/vivere/lib/components"
	"github.com/stojg/vivere/lib/vector"
)

func NewAI(ent *Entity) *AI {
	ai := &AI{
		entity: ent,
	}
	//entity := modelList.Rand()
	//ai.states[ent] = NewSeek(modelList.Get(ent), rigidList.Get(ent), entity)
	return ai
}

type AI struct {
	entity *Entity
}

func (s *AI) Update(elapsed float64) {
	//steering := ent.GetSteering()

	//body.AddForce(steering.linear)
	//model := modelList.Get(id)
	//ste := NewLookWhereYoureGoing(body, model).GetSteering()
	//body.AddTorque(ste.angular)
	body := rigidList.Get(s.entity)

	inst := monitor.FindByEntityID(*s.entity)
	if inst == nil {
		return
	}
	body.SetAwake(true)

	rotSpeed := body.Rotation.Length()
	if rotSpeed < inst.CPU/100 {
		body.AddTorque(&vector.Vector3{0.1, 0.1, 0.1})
	} else if rotSpeed > inst.CPU/100 {
		body.AddTorque(body.Rotation.NewInverse())
	}


}
