package actions

import (
	"fmt"

	"time"

	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/cyberspace/lib/plan"
	"github.com/stojg/steering"
	"github.com/stojg/vector"
)

func NewPatrol(cost float64, me *core.Agent) *patrol {
	a := &patrol{
		cost:  cost,
		agent: me,
	}
	a.preconditions = make(plan.StateList)
	a.effect = make(plan.StateList)
	a.effect.Set(AreaPatrolled, true)
	return a
}

type patrol struct {
	agent         *core.Agent
	steer         *steering.Face
	start         time.Time
	cost          float64
	effect        plan.StateList
	preconditions plan.StateList
	target        *core.GameObject
}

func (a *patrol) Cost() float64 {
	return a.cost
}

func (a *patrol) Effects() plan.StateList {
	return a.effect
}

func (a *patrol) Preconditions() plan.StateList {
	return a.preconditions
}

func (a *patrol) Reset() {
	a.steer = nil
	a.start = time.Time{}
}

func (a *patrol) CheckContextPrecondition(state plan.StateList) bool {
	obj := a.agent.GameObject()
	q := obj.Transform().Orientation()
	test := vector.X().Rotate(q).Inverse()
	a.steer = steering.NewFace(obj.Body(), test.Add(obj.Transform().Position()))
	return true
}

func (a *patrol) MoveTo() interface{} {
	return nil
}

func (a *patrol) Run(agent plan.Agent) (bool, error) {
	if a.start.IsZero() {
		a.start = time.Now()
	}
	if time.Since(a.start) > 1000*time.Millisecond {
		return false, fmt.Errorf("Patrolling took to long, aborting")
	}
	obj := agent.(*core.Agent).GameObject()
	steer := a.steer.Get()
	if steer.Angular().Length() < 1 {
		return true, nil
	}
	obj.Body().AddTorque(steer.Angular())
	return false, nil
}
