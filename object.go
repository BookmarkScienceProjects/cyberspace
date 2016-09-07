package main

import (
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

type Object interface {
	Stateful
	ID() *components.Entity
	Kind() Kind
	Position() *vector.Vector3
	Orientation() *vector.Quaternion
	Awake() bool
	Size() *vector.Vector3
	Rendered() bool
	SetRendered()
}

type Kind byte

const (
	_ Kind = iota
	Monster
	Gunk
)

type GameObject struct {
	id *components.Entity
	*components.Model
	*components.RigidBody
	*components.Collision
	state State
	kind  Kind
	sent  bool
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

func (o *GameObject) Kind() Kind {
	return o.kind
}

func (o *GameObject) State() State {
	return o.state
}

func (o *GameObject) SetState(s State) {
	o.state = s
}
