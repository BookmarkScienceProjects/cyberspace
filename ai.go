package main

import (
	"github.com/stojg/cyberspace/lib/actions"
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/goap"
)

// UpdateAI will run the AI simulation
func UpdateAI(elapsed float64) {

	for _, agent := range core.List.Agents() {
		memory := agent.Memory()
		state := make(goap.StateList)
		goal := make(goap.StateList)

		for _, mem := range memory.Entities() {
			obj := core.List.Get(mem.ID)
			if obj == nil {
				continue
			}
			if obj.CompareTag("food") {
				state.Add(actions.EnemyInSight)
			}
		}

		if memory.Internal().Health < 1 {
			state.Add(goap.Isnt(actions.Healthy))
		} else {
			state.Add(actions.Healthy)
		}
		agent.SetState(state)

		if agent.State().Query(goap.Isnt(actions.Healthy)) {
			goal.Add(actions.Healthy)
		} else if agent.State().Query(actions.EnemyInSight) {
			goal.Add(actions.EnemyKilled)
		} else {
			goal.Add(actions.AreaPatrolled)
		}
		agent.SetGoalState(goal)
		agent.Update()

	}
}

// NewMonsterAgent will return an AI agent with the actions set and goals that a monster have
func NewMonsterAgent() *core.Agent {

	a := []goap.Action{
		actions.NewKillEnemy(2),
		actions.NewHeal(10),
		actions.NewPatrol(4),
	}

	agent := core.NewAgent(a)

	goal := make(goap.StateList)
	agent.SetGoalState(goal)

	core.List.SenseManager().Register(agent)

	initialState := make(goap.StateList)
	initialState.Add(goap.Isnt(actions.Healthy))

	agent.SetState(initialState)
	return agent
}
