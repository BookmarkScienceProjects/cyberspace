package actions

import (
	"fmt"
	"time"

	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/cyberspace/lib/percepts"
	"github.com/stojg/goap"
)

func NewHeal(cost float64) *healAction {
	a := &healAction{
		DefaultAction: goap.NewAction("heal_action", cost),
	}
	a.AddEffect(Healthy)
	return a
}

type healAction struct {
	goap.DefaultAction
	startTime time.Time
}

func (a *healAction) Reset() {
	a.DefaultAction.Reset()
	a.startTime = time.Time{}
}

func (a *healAction) CheckContextPrecondition(agent goap.Agent) bool {
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

func (a *healAction) InRange(agent goap.Agent) bool {
	target := a.Target().(*core.GameObject)
	me := agent.(*core.Agent).GameObject()
	return percepts.Distance(me, target, me.Transform().Scale()[0]) > 0
}

func (a *healAction) Perform(agent goap.Agent) bool {
	fmt.Println("performing healing")
	if a.startTime.IsZero() {
		a.startTime = time.Now()
	}
	a.DefaultAction.Done = true
	agent.(*core.Agent).Memory().Internal().Health += 3
	return true
}
