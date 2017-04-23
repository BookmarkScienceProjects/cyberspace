package actions

import (
	"fmt"
	"time"

	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/cyberspace/lib/plan"
)

func NewHeal(cost float64, me *core.Agent) *healAction {
	a := &healAction{
		cost:  cost,
		agent: me,
	}
	a.preconditions = make(plan.StateList)
	a.preconditions.Set(ClosestHealing, true)
	a.effect = make(plan.StateList)
	a.effect.Set(Healthy, true)
	return a
}

type healAction struct {
	agent         *core.Agent
	start         time.Time
	cost          float64
	effect        plan.StateList
	preconditions plan.StateList
	target        core.ID
}

func (a *healAction) Reset() {
	a.start = time.Time{}
}

func (a *healAction) Cost() float64 {
	return a.cost
}

func (a healAction) Effects() plan.StateList {
	return a.effect
}

func (a healAction) Preconditions() plan.StateList {
	return a.preconditions
}

func (a *healAction) CheckContextPrecondition(state plan.StateList) bool {
	id, found := state.Get(ClosestHealing).(core.ID)
	if !found {
		return false
	}

	object := core.List.Get(id)
	if object == nil {
		return false
	}
	a.target = object.ID()
	return true
}

func (a *healAction) MoveTo() interface{} {
	me := a.agent.GameObject()
	target := core.List.Get(a.target)
	if target == nil {
		return nil
	}
	sqrDist := me.Transform().Position().NewSub(target.Transform().Position()).SquareLength()
	reach := me.Transform().Scale()[0] + target.Transform().Scale()[0]
	if sqrDist > (reach * reach) {
		return target
	}
	return nil
}

func (a *healAction) Run(agent plan.Agent) (bool, error) {
	if a.start.IsZero() {
		a.start = time.Now()
	}

	target := core.List.Get(a.target)
	if target == nil {
		return false, fmt.Errorf("Cant find target %+v anymore", a.target)
	}

	if time.Since(a.start) > 1*time.Second {
		agent.(*core.Agent).Memory().Internal().Health += 3
		return true, nil
	}

	return false, nil
}
