package main

import (
	"github.com/stojg/cyberspace/lib/actions"
	"github.com/stojg/cyberspace/lib/collision"
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

		if !CanSeeTarget(agent) {
			agent.Replan()
		}

		var separationTargets []*vector.Vector3
		for _, monster := range monsters {
			if monster.ID() != agent.GameObject().ID() {
				separationTargets = append(separationTargets, monster.Transform().Position())
			}
		}
		separation := steering.NewSeparation(agent.GameObject().Body(), separationTargets, 2).Get()
		agent.GameObject().Body().AddForce(separation.Linear())

		//// replan
		//lastPlan += elapsed
		//if lastPlan > 1 {
		//	agent.StateMachine.Reset(goap.Idle)
		//	lastPlan = 0
		//}
	}
}
func CanSeeTarget(agent *core.Agent) bool {
	if len(agent.CurrentActions()) == 0 {
		return true
	}
	currAction := agent.CurrentActions()[0]
	target := currAction.Target()
	if target == nil {
		return true
	}

	other := target.(*core.GameObject)
	direction := other.Transform().Position().NewSub(agent.Transform().Position())
	res := collision.Raycast(agent.Transform().Position(), direction, core.List)
	if len(res) == 0 {
		return false
	}

	for _, rr := range res {
		if rr.Distance == 0 {
			continue
		}
		if rr.Collision != other.Collision() {
			return rr.Distance < 0.2
		}
	}
	return true
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
