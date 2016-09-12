package actions

import (
	"github.com/stojg/cyberspace/lib/object"
	"github.com/stojg/goap"
)

func NewGetFood(cost float64) *getFood {
	a := &getFood{
		Action: goap.NewAction("getFood", cost),
	}

	a.AddPrecondition("has_food", false)
	a.AddEffect("has_food", true)

	return a
}

type getFood struct {
	goap.Action
	hasFood bool
}

func (a *getFood) Reset() {}

func (a *getFood) CheckProceduralPrecondition(agent goap.Agent) bool {

	if obj, ok := agent.(object.Object); ok {
		id, found := object.List.Nearest(obj, object.Food)
		if found {
			a.SetTarget(id)
			return true
		}
	}
	return false

}

func (a *getFood) RequiresInRange() bool {
	return true
}

func (a *getFood) IsInRange() bool {
	// first time we call this we are not in range, but next time yes
	return true
}

func (a *getFood) Perform(agent goap.Agent) bool {

	return true
}

func (a *getFood) IsDone() bool {
	return a.hasFood
}
