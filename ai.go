package main

import (
	"github.com/stojg/cyberspace/lib/actions"
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/cyberspace/lib/planning"
)

// UpdateAI will run the AI simulation
func UpdateAI(elapsed float64) {

	for _, agent := range core.List.Agents() {
		memory := agent.Memory()
		state := make(planning.StateList)
		goal := make(planning.StateList)

		for _, mem := range memory.Entities() {
			obj := core.List.Get(mem.ID)
			if obj == nil {
				continue
			}
			if obj.CompareTag("food") {
				state.Add(actions.EnemyInSight)
			}
		}

		agent.SetState(state)

		if agent.State().Query(actions.EnemyInSight) {
			goal.Add(actions.EnemyKilled)
		}
		agent.SetGoalState(goal)
		agent.Update()

	}
}

// NewMonsterAgent will return an AI agent with the actions set and goals that a monster have
func NewMonsterAgent() *core.Agent {

	a := []planning.Action{
		actions.NewKillEnemy(2),
		actions.NewHeal(10),
		actions.NewPatrol(4),
	}

	agent := core.NewAgent(a)

	goal := make(planning.StateList)
	agent.SetGoalState(goal)

	core.List.SenseManager().Register(agent)

	initialState := make(planning.StateList)
	initialState.Add(planning.Isnt(actions.Healthy))

	agent.SetState(initialState)
	return agent
}
