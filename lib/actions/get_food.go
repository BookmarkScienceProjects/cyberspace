package actions

import (
	"fmt"
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/goap"
	"time"
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
	startTime time.Time
	hasFood   bool
	target    *core.GameObject
}

func (a *getFood) Reset() {
	a.hasFood = false
	a.Action.Reset()
	a.startTime = time.Time{}
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

	target, ok := a.Target().(*core.GameObject)
	if !ok {
		fmt.Printf("in actions.getFood.Perform(): %s is not a *GameObject", a.Target())
		return false
	}

	if a.startTime.IsZero() {
		a.startTime = time.Now()
	}

	if time.Since(a.startTime) > 300*time.Millisecond {
		aObject := agent.(*core.Agent)
		fmt.Println("got food")
		aObject.GameObject().Inventory().Add("food", 1)
		core.List.Remove(target)
		a.hasFood = true
	}

	return true
}

func (a *getFood) IsDone() bool {
	return a.hasFood
}
