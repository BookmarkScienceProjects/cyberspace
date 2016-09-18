package main

import (
	"github.com/stojg/cyberspace/lib/actions"
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/goap"
	"github.com/stojg/steering"
	"github.com/stojg/vector"
)

var lastPlan float64

// UpdateAI will run the AI simulation
func UpdateAI(elapsed float64, worldState goap.StateList) {

	monsters := core.List.FindWithTag("monster")

	for _, obj := range core.List.Agents() {

		obj.SetWorldState(worldState)
		obj.StateMachine().Update(obj)

		var separationTargets []*vector.Vector3
		for _, monster := range monsters {
			if monster.ID() != obj.GameObject().ID() {
				separationTargets = append(separationTargets, monster.Transform().Position())
			}
		}
		separation := steering.NewSeparation(obj.GameObject().Body(), separationTargets, 4).Get()
		obj.GameObject().Body().AddForce(separation.Linear())

		lastPlan += elapsed
		// force a replanning every 1 second
		if lastPlan > 1.0 {
			obj.StateMachine().Clear()
			obj.StateMachine().PushState(core.IdleState)
			lastPlan = 0
		}

	}

}

// NewMonsterAgent will return an AI agent with the actions set and goals that a monster have
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
