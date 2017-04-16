package actions

import (
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/cyberspace/lib/percepts"
	"github.com/stojg/goap"
	"github.com/stojg/vector"
	"time"
)

func NewRest(cost float64) *rest {
	a := &rest{
		DefaultAction: goap.NewAction("rest", cost),
	}
	a.AddPrecondition(Full)
	a.AddEffect(Rested)
	return a
}

type rest struct {
	goap.DefaultAction
	startTime time.Time
}

func (a *rest) Reset() {
	a.DefaultAction.Reset()
	a.startTime = time.Time{}
}

func (a *rest) CheckContextPrecondition(agent goap.Agent) bool {
	beds := core.List.FindWithTag("bed")

	if len(beds) < 1 {
		return false
	}
	obj := agent.(*core.Agent).GameObject()

	var target *core.GameObject
	bestConfidence := 0.0

	for _, bed := range beds {
		confidence := percepts.Distance(obj, bed, 50)
		if confidence < 0.01 {
			continue
		}

		if confidence > bestConfidence {
			target = bed
			bestConfidence = confidence
		}
	}
	if target == nil {
		return false
	}
	a.SetTarget(target)
	return true
}

func (a *rest) InRange(agent goap.Agent) bool {
	target := a.Target().(*core.GameObject)
	me := agent.(*core.Agent).GameObject()
	return percepts.Distance(me, target, me.Transform().Scale()[0]*1.5) > 0
}

func (a *rest) Perform(agent goap.Agent) bool {
	if a.startTime.IsZero() {
		a.startTime = time.Now()
	}

	if time.Since(a.startTime) > 10*time.Millisecond {
		agent.AddState(goap.Isnt(Full))
		agent.AddState(Rested)
		a.DefaultAction.Done = true
		if a, found := agent.(*core.Agent); found {
			a.Transform().Scale().Sub(vector.NewVector3(0.00, 0.125, 0.0))
			b := a.GameObject().Body()
			b.SetMass(b.Mass() - 0.1)
			b.MaxAcceleration().Sub(vector.NewVector3(20, 0, 0))
		}
	}

	return true
}
