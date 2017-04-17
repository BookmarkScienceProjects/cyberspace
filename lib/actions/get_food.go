package actions

import (
	"fmt"
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/cyberspace/lib/percepts"
	"github.com/stojg/goap"
	"github.com/stojg/vector"
	"time"
)

func NewGetFood(cost float64) *getFood {
	a := &getFood{
		DefaultAction: goap.NewAction("getFood", cost),
	}
	a.AddEffect(HasFood)
	return a
}

type getFood struct {
	goap.DefaultAction
	startTime time.Time
}

func (a *getFood) Reset() {
	a.DefaultAction.Reset()
	a.startTime = time.Time{}
}

func (a *getFood) CheckContextPrecondition(agent goap.Agent) bool {

	cAgent := agent.(*core.Agent)

	var target *core.GameObject
	bestConfidence := 0.0

	for _, f := range cAgent.Memory().Data() {
		fObj := core.List.Get(f.ID)
		if fObj == nil {
			continue
		}
		provides := false
		for _, g := range fObj.Agent().ProvidesGoals {
			if g == HasFood {
				provides = true
			}
		}
		if !provides {
			continue
		}

		if f.Confidence > bestConfidence {
			fObj := core.List.Get(f.ID)
			if fObj == nil {
				continue
			}
			target = fObj
			bestConfidence = f.Confidence
		}
	}

	if target == nil {
		return false
	}
	a.SetTarget(target)
	return true
}

func (a *getFood) InRange(agent goap.Agent) bool {
	target, ok := a.Target().(*core.GameObject)
	if !ok {
		return false
	}
	me := agent.(*core.Agent).GameObject()
	return percepts.Distance(me, target, me.Transform().Scale()[0]*1.5) > 0
}

func (a *getFood) Perform(agent goap.Agent) bool {

	target, ok := a.Target().(*core.GameObject)
	if !ok {
		return false
	}
	if core.List.Get(target.ID()) == nil {
		return false
	}

	if a.startTime.IsZero() {
		a.startTime = time.Now()
	}

	//if time.Since(a.startTime) > 200*time.Millisecond {
	core.List.Remove(target)
	agent.AddState(HasFood)
	agent.AddState(goap.Isnt(Rested))
	a.DefaultAction.Done = true
	if a, found := agent.(*core.Agent); found {
		a.Transform().Scale().Add(vector.NewVector3(0.00, 0.125, 0.0))
		body := a.GameObject().Body()
		if body == nil {
			panic(fmt.Sprintf("agent has no body %+v\n", a.GameObject().ID()))
		}
		body.SetMass(body.Mass() + 0.1)
		body.MaxAcceleration().Add(vector.NewVector3(20, 0, 0))
	}
	//}
	return true
}
