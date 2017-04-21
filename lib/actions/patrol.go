package actions

import (
	"fmt"

	"time"

	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/cyberspace/lib/planning"
	"github.com/stojg/steering"
	"github.com/stojg/vector"
)

func NewPatrol(cost float64) *patrol {
	a := &patrol{
		DefaultAction: planning.NewAction("patrol", cost),
	}
	a.AddEffect(AreaPatrolled)
	return a
}

type patrol struct {
	planning.DefaultAction
	steer *steering.Face
	start time.Time
}

func (a *patrol) Reset() {
	a.DefaultAction.Reset()
	a.steer = nil
	a.start = time.Time{}
}

func (a *patrol) CheckContextPrecondition(agent planning.Agent) bool {
	obj := agent.(*core.Agent).GameObject()
	q := obj.Transform().Orientation()
	test := vector.X().Rotate(q).Inverse()
	a.steer = steering.NewFace(obj.Body(), test.Add(obj.Transform().Position()))
	return true
}

func (a *patrol) InRange(agent planning.Agent) bool {
	return true
}

func (a *patrol) Perform(agent planning.Agent) bool {
	if a.start.IsZero() {
		a.start = time.Now()
	}
	if time.Since(a.start) > 2*time.Second {
		fmt.Println("Patrolling took to long time, aborting")
		return false
	}
	obj := agent.(*core.Agent).GameObject()
	steer := a.steer.Get()
	if steer.Angular().Length() < 1 {
		a.DefaultAction.Done = true
		return true
	}
	obj.Body().AddTorque(steer.Angular())
	return true
}
