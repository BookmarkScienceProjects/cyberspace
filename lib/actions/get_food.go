package actions

import (
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/cyberspace/lib/percepts"
	"github.com/stojg/goap"
	"github.com/stojg/vector"
	"math"
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

	foods := core.List.FindWithTag("food")

	if len(foods) < 1 {
		return false
	}

	obj := agent.(*core.Agent).GameObject()

	var target *core.GameObject
	bestConfidence := 0.0

	for _, food := range foods {
		confidence := percepts.Distance(obj, food, 50)
		if confidence < 0.01 {
			continue
		}

		// 269 degrees of view cone
		confidence *= percepts.InViewCone(obj, food, 2*math.Pi)

		if confidence > bestConfidence {
			target = food
			bestConfidence = confidence
		}
	}

	if target != nil {
		a.SetTarget(target)
	}
	return true
}

func (a *getFood) InRange(agent goap.Agent) bool {
	target, ok := a.Target().(*core.GameObject)
	if !ok {
		return false
	}
	me := agent.(*core.Agent).GameObject()
	return percepts.Distance(me, target, me.Transform().Scale()[0]*1.2) > 0
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

	if time.Since(a.startTime) > 200*time.Millisecond {
		core.List.Remove(target)
		agent.AddState(HasFood)
		agent.AddState(goap.Isnt(Rested))
		a.DefaultAction.Done = true
		if a, found := agent.(*core.Agent); found {
			a.Transform().Scale().Add(vector.NewVector3(0.00, 0.125, 0.0))
			b := a.GameObject().Body()
			b.SetMass(b.Mass() + 0.1)
			b.MaxAcceleration().Add(vector.NewVector3(20, 0, 0))
		}
	}
	return true
}
