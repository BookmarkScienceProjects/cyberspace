package actions

import (
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/goap"
	"github.com/stojg/steering"
	"github.com/stojg/vector"
)

func NewScan(cost float64) *scan {
	a := &scan{
		DefaultAction: goap.NewAction("scan", cost),
	}
	a.AddEffect(HasFood)
	return a
}

type scan struct {
	goap.DefaultAction
	steer *steering.Face
}

func (a *scan) Reset() {
	a.DefaultAction.Reset()
	a.steer = nil
}

func (a *scan) CheckContextPrecondition(agent goap.Agent) bool {
	obj := agent.(*core.Agent).GameObject()
	q := obj.Transform().Orientation()
	test := vector.X().Rotate(q).Inverse()
	a.steer = steering.NewFace(obj.Body(), test.Add(obj.Transform().Position()))
	return true
}

func (a *scan) InRange(agent goap.Agent) bool {
	return true
}

func (a *scan) Perform(agent goap.Agent) bool {
	obj := agent.(*core.Agent).GameObject()

	steer := a.steer.Get()

	if steer.Angular().Length() < 1 {
		a.DefaultAction.Done = true
		return true
	}

	//obj.Body().AddForce(steer.Linear())
	obj.Body().AddTorque(steer.Angular())

	return true

}
