package actions

import (
	"github.com/stojg/goap"
	"time"
)

func NewEat(cost float64) *eat {
	a := &eat{
		Action: goap.NewAction("eat", cost),
	}
	a.AddPrecondition(HasFood)
	a.AddPrecondition(goap.Isnt(Full))
	a.AddEffect(Full)
	return a
}

type eat struct {
	goap.Action
	startTime time.Time
}

func (a *eat) Reset() {
	a.Action.Reset()
	a.startTime = time.Time{}
}

func (a *eat) CheckProceduralPrecondition(agent goap.Agent) bool {
	return true
}

func (a *eat) RequiresInRange() bool {
	return false
}

func (a *eat) Perform(agent goap.Agent) bool {
	if a.startTime.IsZero() {
		a.startTime = time.Now()
	}

	if time.Since(a.startTime) > 10*time.Millisecond {
		agent.AddState(goap.Dont(HasFood))
		agent.AddState(Full)
		a.SetIsDone()
	}

	return true
}
