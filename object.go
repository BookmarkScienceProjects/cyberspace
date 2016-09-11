package main

import (
	"fmt"
	"github.com/stojg/cyberspace/lib/object"
	"github.com/stojg/goap"
	"github.com/stojg/vector"
	"github.com/stojg/vivere/lib/components"
)

type State int

const (
	stateDead State = iota
	stateIdle
	stateMoving
)

type Stateful interface {
	State() State
	SetState(State)
}

type GameObject struct {
	id *components.Entity
	*components.Model
	*components.RigidBody
	*components.Collision
	kind  object.Kind
	sent  bool
	agent *goap.GoapAgent
	state goap.StateList
	// @todo idea, add tags for searching
}

func (o *GameObject) Rendered() bool {
	return o.sent
}
func (o *GameObject) SetRendered() {
	o.sent = true
}

func (o *GameObject) ID() *components.Entity {
	return o.id
}

func (o *GameObject) Size() *vector.Vector3 {
	return o.Scale
}

func (o *GameObject) Kind() object.Kind {
	return o.kind
}

func (p *GameObject) GetWorldState() goap.StateList {
	return p.state
}

func (p *GameObject) CreateGoalState() goap.StateList {
	goal := make(goap.StateList, 0)
	goal["is_full"] = false
	return goal
}

func (p *GameObject) Update() {
	p.agent.Update()
}

func (p *GameObject) PlanFailed(failedGoal goap.StateList) {
	//fmt.Printf("Planning failed: %v\n", failedGoal)
}

func (p *GameObject) PlanFound(goal goap.StateList, actions []goap.Actionable) {
	fmt.Printf("Planning success to goal %v with actions %v\n", goal, actions)
}

func (p *GameObject) ActionsFinished() {
	fmt.Println("all actions finished")
}

func (p *GameObject) PlanAborted(aborter goap.Actionable) {
	fmt.Printf("plan aborted by %v\n", aborter)
}

func (p *GameObject) MoveAgent(nextAction goap.Actionable) bool {
	return true
}
