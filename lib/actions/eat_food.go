package actions

import (
	"github.com/stojg/goap"
	"time"
)

func NewEat(cost float64) *eat {
	a := &eat{
		DefaultAction: goap.NewAction("eat", cost),
	}
	a.AddPrecondition(HasFood)
	a.AddPrecondition(goap.Isnt(Full))
	a.AddEffect(Full)
	return a
}

type eat struct {
	goap.DefaultAction
	startTime time.Time
}

func (a *eat) Reset() {
	a.DefaultAction.Reset()
	a.startTime = time.Time{}
}

func (a *eat) CheckContextPrecondition(agent goap.Agent) bool {
	return true
}

func (a *eat) InRange(agent goap.Agent) bool {
	return true
}

func (a *eat) Perform(agent goap.Agent) bool {
	if a.startTime.IsZero() {
		a.startTime = time.Now()
	}

	if time.Since(a.startTime) > 10*time.Millisecond {
		agent.AddState(goap.Dont(HasFood))
		agent.AddState(Full)
		a.DefaultAction.Done = true
	}

	return true
}
