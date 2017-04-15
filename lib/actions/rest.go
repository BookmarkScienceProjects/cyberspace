package actions

import (
	"github.com/stojg/cyberspace/lib/core"
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

	gameObject, found := agent.(*core.Agent)
	if !found {
		return false
	}
	t := gameObject.Transform()
	if t == nil {
		panic("wtf?")
	}
	agentPos := t.Position().Clone()

	nearest := beds[0]
	nearestDistance := agentPos.NewSub(beds[0].Transform().Position()).Length()

	for _, food := range beds {
		dist := agentPos.NewSub(food.Transform().Position()).Length()
		if dist < nearestDistance {
			nearest = food
			nearestDistance = dist
		}
	}
	a.SetTarget(nearest)
	return true
}

func (a *rest) InRange(agent goap.Agent) bool {
	target := a.Target().(*core.GameObject)
	gameObject := agent.(*core.Agent)
	agentTransform := gameObject.Transform()
	agentPos := agentTransform.Position().Clone()
	dist := agentPos.NewSub(target.Transform().Position()).Length()
	return dist < agentTransform.Scale()[0]*1.5
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
			b.MaxAcceleration().Sub(vector.NewVector3(20, 20, 20))
		}
	}

	return true
}
