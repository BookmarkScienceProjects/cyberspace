package actions

import (
	"math"
	"time"

	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/cyberspace/lib/percepts"
	"github.com/stojg/goap"
)

func NewKillEnemy(cost float64) *killEnemyAction {
	a := &killEnemyAction{
		DefaultAction: goap.NewAction("kill_enemy", cost),
	}
	a.AddPrecondition(EnemyInSight)
	a.AddEffect(EnemyKilled)
	return a
}

type killEnemyAction struct {
	goap.DefaultAction
	startTime time.Time
}

func (a *killEnemyAction) Reset() {
	a.DefaultAction.Reset()
	a.startTime = time.Time{}
}

func (a *killEnemyAction) CheckContextPrecondition(agent goap.Agent) bool {

	cAgent := agent.(*core.Agent)

	var target *core.GameObject
	closestDistance := math.MaxFloat64

	for _, f := range cAgent.Memory().Entities() {
		if f.Distance < closestDistance {
			obj := core.List.Get(f.ID)
			if obj == nil {
				continue
			}
			if obj.CompareTag("food") {
				closestDistance = f.Distance
				target = obj
			}
		}
	}

	if target == nil {
		return false
	}
	a.SetTarget(target)
	return true
}

func (a *killEnemyAction) InRange(agent goap.Agent) bool {
	target, ok := a.Target().(*core.GameObject)
	if !ok {
		return false
	}
	me := agent.(*core.Agent).GameObject()
	return percepts.Distance(me, target, me.Transform().Scale()[0]*2) > 0
}

func (a *killEnemyAction) Perform(agent goap.Agent) bool {
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

	core.List.Remove(target)
	agent.(*core.Agent).Memory().Internal().Health -= 1
	a.DefaultAction.Done = true
	return true
}
