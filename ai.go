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

	for _, agent := range core.List.Agents() {

		agent.Update()

		var separationTargets []*vector.Vector3
		for _, monster := range monsters {
			if monster.ID() != agent.GameObject().ID() {
				separationTargets = append(separationTargets, monster.Transform().Position())
			}
		}
		separation := steering.NewSeparation(agent.GameObject().Body(), separationTargets, 2).Get()
		agent.GameObject().Body().AddForce(separation.Linear())

		// replan
		lastPlan += elapsed
		if lastPlan > 1 {
			agent.StateMachine.Reset(goap.Idle)
			lastPlan = 0
		}
	}
}

// NewMonsterAgent will return an AI agent with the actions set and goals that a monster have
func NewMonsterAgent() *core.Agent {

	a := []goap.Action{
		actions.NewEat(1),
		actions.NewGetFood(2),
		actions.NewRest(10),
	}

	agent := core.NewAgent(a)

	goal := make(goap.StateList)
	goal.Add(actions.Rested)
	goal.Add(actions.Full)
	agent.SetGoalState(goal)

	initialState := make(goap.StateList)
	initialState.Add(goap.Isnt(actions.Full))

	agent.SetState(initialState)
	return agent
}
