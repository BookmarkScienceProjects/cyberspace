package main

import (
	"github.com/stojg/steering"
	"github.com/stojg/vector"
	"github.com/stojg/vivere/lib/components"
	"math"
	"math/rand"
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
	Size() *vector.Vector3
}

type Kind byte

const (
	_ Kind = iota
	Monster
	Gunk
)

type stuffList []Object

// idea for the future, let k be a bitmask
func (l stuffList) ofKind(k Kind) stuffList {
	result := make(stuffList, 0)
	for _, o := range l {
		if o.Kind() == k && o.State() != stateDead {
			result = append(result, o)
		}
	}
	return result
}

func createFood(width, height, depth float64) *GameObject {
	eID := entities.Create()
	f := &GameObject{
		id:        eID,
		state:     stateIdle,
		kind:      Gunk,
		Model:     modelList.New(eID, width, height, depth, 2),
		RigidBody: rigidList.New(eID, 10),
		Collision: collisionList.New(eID, width, height, depth),
	}
	f.Position().Set(rand.Float64()*1000-500, height/2, rand.Float64()*1000-500)
	return f
}

func createMonster(width, height, depth float64) *GameObject {
	eID := entities.Create()
	monster := &GameObject{
		kind:      Monster,
		id:        eID,
		Model:     modelList.New(eID, width, height, depth, 1),
		RigidBody: rigidList.New(eID, 1),
		Collision: collisionList.New(eID, width, height, depth),
	}
	monster.Position().Set(rand.Float64()*500-250, height/2, rand.Float64()*500-250)
	return monster
}

type GameObject struct {
	id *components.Entity
	*components.Model
	*components.RigidBody
	*components.Collision
	state State
	kind  Kind
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

type World struct {
	timer float64
	list  stuffList
}

func (w *World) Update(elapsed float64) {

	w.timer += elapsed
	if w.timer > 0.1 {
		w.timer -= 0.1
	}

	if len(w.list.ofKind(Gunk)) <= 40 {
		w.list = append(w.list, createFood(3, 3, 3))
	}

	if len(w.list.ofKind(Monster)) <= 10 {
		monster := createMonster(10, 10, 10)
		monster.state = stateIdle
		w.list = append(w.list, monster)
	}

	// run the AI for the individual entities
	for _, obj := range w.list.ofKind(Monster) {
		// @todo exclude dead entities
		monster, ok := obj.(*GameObject)
		if !ok {
			continue
		}
		var closestF Object
		closestDistance := math.MaxFloat64
		for _, f := range w.list.ofKind(Gunk) {
			distance := monster.Position().NewSub(f.Position()).SquareLength()
			if obj.Position().NewSub(f.Position()).SquareLength() < closestDistance {
				closestF = f
				closestDistance = distance
			}
		}
		if closestF != nil {

			if math.Sqrt(closestDistance) < monster.Scale[0] {
				closestF.SetState(stateDead)
				monster.MaxAcceleration.Add(vector.NewVector3(5, 5, 5))
			}

			arrive := steering.NewSeek(monster.Model, monster.RigidBody, closestF.Position())
			st := arrive.Get()
			monster.AddForce(st.Linear())
			monster.AddTorque(st.Angular())
		}
		monster.Position()[1] = monster.Scale[1] / 2
	}
}
