package actions

import (
	"fmt"
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/goap"
)

func NewGetFood(cost float64) *getFood {
	a := &getFood{
		Action: goap.NewAction("getFood", cost),
	}
	a.AddPrecondition("food_exists", true)
	a.AddEffect("has_food", true)
	return a
}

type getFood struct {
	goap.Action
	hasFood bool
	target  *core.GameObject
}

func (a *getFood) Reset() {
	a.hasFood = false
	a.Action.Reset()
}

func (a *getFood) CheckProceduralPrecondition(agent goap.Agent) bool {

	foods := core.List.FindWithTag("food")

	if len(foods) < 1 {
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

	target, found := a.Target().(*core.GameObject)
	if !found {
		fmt.Printf("in actions.getFood.Perform(): %s is not a *GameObject", a.Target())
		return false
	}
	core.List.Remove(target)
	a.hasFood = true
	return true
}

func (a *getFood) IsDone() bool {
	return a.hasFood
}
