package main

import (
	"github.com/stojg/cyberspace/lib/actions"
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/goap"
)

func NewMonsterAgent() *core.Agent {

	actions := []goap.Actionable{
		actions.NewEat(1),
		actions.NewGetFood(2),
		actions.NewRest(10),
	}

	agent := core.NewAgent()
	agent.SetAvailableActions(actions)

	agent.AddGoal("is_hungry", false)

	agent.SetState("is_hungry", true)

	agent.StateMachine().PushState(core.IdleState)

	return agent
}

var lastPlan float64

func UpdateAI(elapsed float64, worldState goap.StateList) {
	for _, obj := range core.List.Agents() {
		obj.SetWorldState(worldState)
		obj.StateMachine().Update(obj)

		lastPlan += elapsed
		// force a replanning every 1 second
		if lastPlan > 1.0 {
			obj.StateMachine().Clear()
			obj.StateMachine().PushState(core.IdleState)
			lastPlan = 0
		}

	}

}
