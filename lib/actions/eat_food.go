package actions

import (
	"github.com/stojg/goap"
)

func NewEat(cost float64) *eat {
	a := &eat{
		Action: goap.NewAction("eat", cost),
	}
	a.AddPrecondition("has_food", true)
	a.AddPrecondition("is_hungry", true)
	a.AddEffect("is_hungry", false)
	return a
}

type eat struct {
	goap.Action
}

func (a *eat) Reset() {
	a.Action.Reset()
}

func (a *eat) CheckProceduralPrecondition(agent goap.Agent) bool {
	return true
}

func (a *eat) RequiresInRange() bool {
	return false
}

func (a *eat) Perform(agent goap.Agent) bool {
	return true
}

func (a *eat) IsDone() bool {
	return true
}
