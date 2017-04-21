package actions

import (
	"time"

	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/cyberspace/lib/percepts"
	"github.com/stojg/cyberspace/lib/planning"
)

func NewHeal(cost float64) *healAction {
	a := &healAction{
		DefaultAction: planning.NewAction("heal_action", cost),
	}
	a.AddEffect(Healthy)
	return a
}

type healAction struct {
	planning.DefaultAction
	start time.Time
}

func (a *healAction) Reset() {
	a.DefaultAction.Reset()
	a.start = time.Time{}
}

func (a *healAction) CheckContextPrecondition(agent planning.Agent) bool {
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

func (a *healAction) InRange(agent planning.Agent) bool {
	target := a.Target().(*core.GameObject)
	me := agent.(*core.Agent).GameObject()
	return percepts.Distance(me, target, me.Transform().Scale()[0]) > 0
}

func (a *healAction) Perform(agent planning.Agent) bool {
	if a.start.IsZero() {
		a.start = time.Now()
	}

	if time.Since(a.start) > 1*time.Second {
		agent.(*core.Agent).Memory().Internal().Health += 4
		a.DefaultAction.Done = true
	}

	return true
}
