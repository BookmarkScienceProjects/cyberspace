package actions

import (
	"fmt"
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/goap"
	"time"
)

func NewEat(cost float64) *eat {
	a := &eat{
		Action: goap.NewAction("eat", cost),
	}
	a.AddPrecondition("has_food", true)
	a.AddPrecondition("is_hungry", true)
	a.AddEffect("is_hungry", false)
	return a
}

type eat struct {
	goap.Action
	startTime time.Time
	ate       bool
}

func (a *eat) Reset() {
	a.Action.Reset()
	a.startTime = time.Time{}
	a.ate = false
}

func (a *eat) CheckProceduralPrecondition(agent goap.Agent) bool {
	return true
}

func (a *eat) RequiresInRange() bool {
	return false
}

func (a *eat) Perform(agent goap.Agent) bool {
	if a.startTime.IsZero() {
		a.startTime = time.Now()
	}

	if time.Since(a.startTime) > 500*time.Millisecond {
		aObject := agent.(*core.Agent)
		fmt.Println("ate")
		aObject.GameObject().Inventory().Remove("food", 1)
		a.ate = true
	}

	return true
}

func (a *eat) IsDone() bool {
	return a.ate
}
