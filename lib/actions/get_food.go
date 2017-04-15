package actions

import (
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/goap"
	"github.com/stojg/vector"
	"time"
)

func NewGetFood(cost float64) *getFood {
	a := &getFood{
		Action: goap.NewAction("getFood", cost),
	}
	a.AddEffect(HasFood)
	return a
}

type getFood struct {
	goap.Action
	startTime time.Time
}

func (a *getFood) Reset() {
	a.Action.Reset()
	a.startTime = time.Time{}
}

func (a *getFood) CheckContextPrecondition(agent goap.Agent) bool {

	foods := core.List.FindWithTag("food")

	if len(foods) < 1 {
		return false
	}

	gameObject, found := agent.(*core.Agent)
	if !found {
		return false
	}
	t := gameObject.Transform()
	agentPos := t.Position().Clone()

	nearest := foods[0]
	nearestDistance := agentPos.NewSub(foods[0].Transform().Position()).Length()

	for _, food := range foods {
		dist := agentPos.NewSub(food.Transform().Position()).Length()
		if dist < nearestDistance {
			nearest = food
			nearestDistance = dist
		}
	}

	a.SetTarget(nearest)
	return true
}

func (a *getFood) RequiresInRange() bool {
	return true
}

func (a *getFood) Perform(agent goap.Agent) bool {

	target, ok := a.Target().(*core.GameObject)
	if !ok {
		return false
	}

	t := core.List.Get(target.ID())
	if t == nil {
		return false
	}

	if a.startTime.IsZero() {
		a.startTime = time.Now()
	}

	if time.Since(a.startTime) > 200*time.Millisecond {
		core.List.Remove(target)
		agent.AddState(HasFood)
		agent.AddState(goap.Isnt(Rested))
		a.SetIsDone()
		if a, found := agent.(*core.Agent); found {
			a.Transform().Scale().Add(vector.NewVector3(0.00, 0.125, 0.0))
			b := a.GameObject().Body()
			b.SetMass(b.Mass() + 0.1)
			b.MaxAcceleration().Add(vector.NewVector3(20, 20, 20))
		}
	}

	return true
}
