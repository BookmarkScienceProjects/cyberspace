package actions

import (
	"fmt"

	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/goap"
	"github.com/stojg/steering"
	"github.com/stojg/vector"
)

func NewPatrol(cost float64) *patrol {
	a := &patrol{
		DefaultAction: goap.NewAction("patrol", cost),
	}
	a.AddEffect(AreaPatrolled)
	return a
}

type patrol struct {
	goap.DefaultAction
	steer *steering.Face
}

func (a *patrol) Reset() {
	a.DefaultAction.Reset()
	a.steer = nil
}

func (a *patrol) CheckContextPrecondition(agent goap.Agent) bool {
	obj := agent.(*core.Agent).GameObject()
	q := obj.Transform().Orientation()
	test := vector.X().Rotate(q).Inverse()
	a.steer = steering.NewFace(obj.Body(), test.Add(obj.Transform().Position()))
	return true
}

func (a *patrol) InRange(agent goap.Agent) bool {
	return true
}

func (a *patrol) Perform(agent goap.Agent) bool {
	fmt.Println("performing patrolling")
	obj := agent.(*core.Agent).GameObject()
	steer := a.steer.Get()
	if steer.Angular().Length() < 1 {
		a.DefaultAction.Done = true
		return true
	}
	obj.Body().AddTorque(steer.Angular())
	return true
}
