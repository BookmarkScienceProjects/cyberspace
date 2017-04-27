package main

import (
	"math"

	"github.com/stojg/cyberspace/lib/actions"
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/cyberspace/lib/plan"
)

// UpdateAI will run the AI simulation
func UpdateAI(elapsed float64) {

	for _, agent := range core.List.Agents() {
		obj := agent.GameObject()
		if obj.CompareTag("monster") {
			HandleMonsterAi(agent, elapsed)
		} else if obj.CompareTag("food") {
			HandleFoodAI(agent, elapsed)
		}
	}
}

func HandleMonsterAi(ag *core.Agent, elapsed float64) {
	memory := ag.Memory()
	me := ag.GameObject()
	state := make(plan.StateList)

	if ag.Memory().Internal().Health > 1 {
		state.Set(actions.Healthy, true)
	}

	closestDistance := math.MaxFloat64
	for _, mem := range memory.Entities() {
		obj := core.List.Get(mem.ID)
		if obj == nil {
			continue
		}
		if obj.CompareTag("food") {
			dist := me.SqrDistance(obj)
			if dist < closestDistance {
				closestDistance = dist
				state.Set(actions.EnemyInSight, obj.ID())
			}
		}
	}

	//closestDistance = math.MaxFloat64
	//for _, mem := range memory.Entities() {
	//	obj := core.List.Get(mem.ID)
	//	if obj == nil {
	//		continue
	//	}
	//	if obj.CompareTag("healing_station") {
	//		dist := me.SqrDistance(obj)
	//		if dist < closestDistance {
	//			closestDistance = dist
	//			state.Set(actions.ClosestHealing, obj.ID())
	//		}
	//	}
	//}

	state.Set(actions.Healthy, true)

	me.Transform().Scale().Set(1, 1-ag.Memory().Internal().Health*0.1, 1)
	me.Body().SetMass(1 - ag.Memory().Internal().Health*0.15)
	needsReplanning := !ag.State().Compare(state)
	ag.SetState(state)
	if needsReplanning {
		ag.Replan()
	}
	ag.Update(elapsed)
}

func HandleFoodAI(ag *core.Agent, elapsed float64) {
	memory := ag.Memory()
	me := ag.GameObject()
	state := make(plan.StateList)

	if ag.Memory().Internal().Health > 1 {
		state.Set(actions.Healthy, true)
	}

	closestDistance := math.MaxFloat64
	for _, mem := range memory.Entities() {
		obj := core.List.Get(mem.ID)
		if obj == nil {
			continue
		}
		if obj.CompareTag("grass") {
			dist := me.SqrDistance(obj)
			if dist < closestDistance {
				closestDistance = dist
				state.Set(actions.EnemyInSight, obj.ID())
			}
		}
	}
	state.Set(actions.Healthy, true)

	needsReplanning := !ag.State().Compare(state)
	ag.SetState(state)
	if needsReplanning {
		ag.Replan()
	}
	ag.Update(elapsed)
}

// NewMonsterAgent will return an AI agent with the actions set and goals that a monster have
func NewMonsterAgent() *core.Agent {
	agent := core.NewAgent()
	agent.AddAction(actions.NewKillEnemy(2, agent))
	//agent.AddAction(actions.NewHeal(1, agent))
	agent.AddAction(actions.NewPatrol(4, agent))
	goalSet := make([]plan.StateList, 0)

	firstGoal := make(plan.StateList)
	firstGoal.Set(actions.EnemyKilled, true)
	goalSet = append(goalSet, firstGoal)

	secondGoal := make(plan.StateList)
	secondGoal.Set(actions.AreaPatrolled, true)
	goalSet = append(goalSet, secondGoal)

	agent.SetGoals(goalSet)

	core.List.SenseManager().Register(agent)
	return agent
}

func NewFoodAgent() *core.Agent {
	agent := core.NewAgent()
	agent.AddAction(actions.NewPatrol(4, agent))
	agent.AddAction(actions.NewKillEnemy(2, agent))
	goalSet := make([]plan.StateList, 0)

	firstGoal := make(plan.StateList)
	firstGoal.Set(actions.EnemyKilled, true)
	goalSet = append(goalSet, firstGoal)

	secondGoal := make(plan.StateList)
	secondGoal.Set(actions.AreaPatrolled, true)
	goalSet = append(goalSet, secondGoal)

	agent.SetGoals(goalSet)

	core.List.SenseManager().Register(agent)
	return agent
}
