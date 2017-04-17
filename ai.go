package main

import (
	"github.com/stojg/cyberspace/lib/actions"
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/cyberspace/lib/percepts"
	"github.com/stojg/goap"
	//"github.com/stojg/steering"
	//"github.com/stojg/vector"
	"math"
)

var lastPlan float64

// UpdateAI will run the AI simulation
func UpdateAI(elapsed float64) {

	//monsters := core.List.FindWithTag("monster")

	for _, agent := range core.List.Agents() {

		lastPlan += elapsed
		if lastPlan > 1.0 {
			agent.Memory().Tick()
			agent.Replan()
			lastPlan = 0
		}

		obj := agent.GameObject()

		// sensor update
		for _, food := range core.List.All() {
			confidence := percepts.Distance(obj, food, 50)
			if confidence < 0.01 {
				continue
			}

			if percepts.InViewCone(obj, food, math.Pi) < 0.01 {
				continue
			}

			if !percepts.CanSeeTarget(obj, food) {
				continue
			}

			agent.Memory().Add(&core.WorkingMemoryFact{
				Type:       core.Item,
				ID:         food.ID(),
				Confidence: confidence,
				Position:   food.Transform().Position(),
				States:     food.Agent().ProvidesGoals,
			})
		}

		agent.Update()

		//var separationTargets []*vector.Vector3
		//for _, monster := range monsters {
		//	if monster.ID() != agent.GameObject().ID() {
		//		separationTargets = append(separationTargets, monster.Transform().Position())
		//	}
		//}
		//separation := steering.NewSeparation(agent.GameObject().Body(), separationTargets, 2).Get()
		//agent.GameObject().Body().AddForce(separation.Linear())

		//// replan
		//lastPlan += elapsed
		//if lastPlan > 1 {
		//	agent.StateMachine.Reset(goap.Idle)
		//	lastPlan = 0
		//}
	}
}

// NewMonsterAgent will return an AI agent with the actions set and goals that a monster have
func NewMonsterAgent() *core.Agent {

	a := []goap.Action{
		actions.NewEat(1),
		actions.NewGetFood(2),
		actions.NewScan(4),
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

func NewNoopAgent() *core.Agent {
	return core.NewAgent([]goap.Action{})
}

func NewFoodAgent() *core.Agent {

	a := []goap.Action{}

	agent := core.NewAgent(a)
	agent.ProvidesGoals = append(agent.ProvidesGoals, actions.HasFood)

	goal := make(goap.StateList)
	agent.SetGoalState(goal)

	initialState := make(goap.StateList)
	agent.SetState(initialState)
	return agent
}
