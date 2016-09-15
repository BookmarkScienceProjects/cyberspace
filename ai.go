package main

import (
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/goap"
)

type eatAction struct {
	goap.Action
}

func (a *eatAction) Reset() {
}

func (a *eatAction) CheckProceduralPrecondition(agent goap.Agent) bool {
	return true
}

func (a *eatAction) RequiresInRange() bool {
	return false
}

func (a *eatAction) IsInRange() bool {
	return true
}

func (a *eatAction) Perform(agent goap.Agent) bool {
	return true
}

type HunterAI struct {
	worldState goap.StateList
	innerState goap.StateList
	gameObject *core.GameObject
}

func (h *HunterAI) SetWorldState(s goap.StateList) {
	h.worldState = s
}

func (h *HunterAI) SetGameObject(g *core.GameObject) {
	h.gameObject = g
}

func (h *HunterAI) AvailableActions() []goap.Actionable {
	actions := make([]goap.Actionable, 0)
	actions = append(actions, &eatAction{
		Action: goap.NewAction("eat", 1),
	})
	return actions
}

func (h *HunterAI) WorldState() goap.StateList {
	rs := make(goap.StateList, len(h.worldState)+len(h.innerState))
	for k, v := range h.worldState {
		rs[k] = v
	}
	for k, v := range h.innerState {
		rs[k] = v
	}
	return rs
}

func (h *HunterAI) Goals() goap.StateList {
	goals := make(goap.StateList, 1)
	goals["eat"] = true
	return goals
}

func (h *HunterAI) Update(elapsed float64) {
	goap.Plan(h, h.AvailableActions(), h.WorldState(), h.Goals())
}

func UpdateAI(elapsed float64, worldState goap.StateList) {

	for _, obj := range core.List.AIs() {
		obj.SetWorldState(worldState)
		obj.Update(elapsed)
	}

}
