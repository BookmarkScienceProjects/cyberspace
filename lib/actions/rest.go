package actions

import (
	"github.com/stojg/goap"
)

func NewRest(cost float64) *rest {
	a := &rest{
		Action: goap.NewAction("rest", cost),
	}
	a.AddPrecondition("is_tired", true)
	a.AddEffect("is_tired", false)
	return a
}

type rest struct {
	goap.Action
	isRested bool
}

func (a *rest) Reset() {}

func (a *rest) CheckProceduralPrecondition(agent goap.Agent) bool {
	return true
}

func (a *rest) RequiresInRange() bool {
	return false
}

func (a *rest) IsInRange() bool {
	return true
}

func (a *rest) Perform(agent goap.Agent) bool {
	return true
}

func (a *rest) IsDone() bool {
	return a.isRested
}
