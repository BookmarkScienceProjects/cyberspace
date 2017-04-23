package actions

import (
	"fmt"
	"time"

	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/cyberspace/lib/plan"
)

func NewKillEnemy(cost float64, me *core.Agent) *killEnemyAction {
	a := &killEnemyAction{
		cost:  cost,
		agent: me,
	}
	a.preconditions = make(plan.StateList)
	a.preconditions.Set(EnemyInSight, true)
	a.preconditions.Set(Healthy, true)
	a.effect = make(plan.StateList)
	a.effect.Set(EnemyKilled, true)
	return a
}

type killEnemyAction struct {
	agent         *core.Agent
	start         time.Time
	cost          float64
	effect        plan.StateList
	preconditions plan.StateList
	target        core.ID
}

func (a *killEnemyAction) Cost() float64 {
	return a.cost
}

func (a *killEnemyAction) Effects() plan.StateList {
	return a.effect
}

func (a *killEnemyAction) Preconditions() plan.StateList {
	return a.preconditions
}

func (a *killEnemyAction) Reset() {
	a.target = 0
	a.start = time.Time{}
}

func (a *killEnemyAction) CheckContextPrecondition(state plan.StateList) bool {
	id, found := state.Get(EnemyInSight).(core.ID)
	if !found {
		return false
	}

	enemy := core.List.Get(id)
	if enemy == nil {
		return false
	}

	a.target = enemy.ID()

	return true
}

func (a *killEnemyAction) MoveTo() interface{} {
	me := a.agent.GameObject()
	target := core.List.Get(a.target)
	if target == nil {
		return nil
	}
	dist := me.Transform().Position().NewSub(target.Transform().Position()).SquareLength()
	reach := (me.Transform().Scale()[0] + target.Transform().Scale()[0])
	if dist > (reach * reach) {
		return target
	}
	return nil
}

func (a *killEnemyAction) Run(agent plan.Agent) (bool, error) {
	//if a.start.IsZero() {
	//	a.start = time.Now()
	//}
	target := core.List.Get(a.target)
	if target == nil {
		return false, fmt.Errorf("Cant find target %+v anymore", a.target)
	}
	if time.Since(a.start) > 500*time.Millisecond {
		agent.(*core.Agent).Memory().Internal().Health -= 1
		core.List.Remove(target)
	}
	return false, nil
}
