package main

import (
	"github.com/stojg/vector"
	"github.com/stojg/vivere/lib/components"
)

type GameObject struct {
	id *components.Entity
	*components.Model
	*components.RigidBody
	*components.Collision
	state State
	kind  Kind
	// @todo idea, add tags for searching
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
