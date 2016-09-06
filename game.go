package main

import (
	"github.com/stojg/steering"
	"github.com/stojg/vector"
	"github.com/stojg/vivere/lib/components"
	"math"
	"math/rand"
	"sync"
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
}

type Kind byte

const (
	_ Kind = iota
	Monster
	Gunk
)

type stuffList struct {
	next uint16

	sync.Mutex
	items   [math.MaxUint16]Object
	deleted []Object
}

func (l *stuffList) Add(o Object) {
	// dont add anymore, we are full
	if l.next+1 == math.MaxUint16 {
		Println("stuff list full")
		return
	}
	l.Lock()
	l.items[l.next] = o
	l.next++
	l.Unlock()
}

func (l *stuffList) Remove(i uint16) {
	l.Lock()
	modelList.Delete(l.items[i].ID())
	rigidList.Delete(l.items[i].ID())
	collisionList.Delete(l.items[i].ID())

	// decrement the current next
	l.next--
	// keep a record of deleted entities so they can be flushed to the view
	l.deleted = append(l.deleted, l.items[i])
	// take the object that was last and replace it with object to be removed
	l.items[i] = l.items[l.next]

	l.Unlock()
}

func (l *stuffList) All() map[uint16]Object {
	result := make(map[uint16]Object, 0)
	l.Lock()
	for i := uint16(0); i < l.next; i++ {
		result[i] = l.items[i]
	}
	l.Unlock()
	return result
}

// idea for the future, let k be a bitmask
func (l *stuffList) ofKind(k Kind) map[uint16]Object {
	result := make(map[uint16]Object, 0)
	l.Lock()
	for i := uint16(0); i < l.next; i++ {
		if l.items[i].Kind() == k {
			result[i] = l.items[i]
		}
	}
	l.Unlock()
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

type World struct {
	timer float64
	list  *stuffList
}

func (w *World) Update(elapsed float64) {

	w.timer += elapsed
	if w.timer > 10 {
		w.timer -= 10
	}

	for len(w.list.ofKind(Gunk)) < 1000 {
		w.list.Add(createFood(3, 3, 3))
	}

	for len(w.list.ofKind(Monster)) < 10 {
		monster := createMonster(10, 10, 10)
		monster.state = stateIdle
		w.list.Add(monster)
	}

	// run the AI for the individual entities
	for _, obj := range w.list.ofKind(Monster) {
		// @todo exclude dead entities
		monster := obj.(*GameObject)
		found := false
		var closestIndex uint16
		closestDistance := math.MaxFloat64
		gunks := w.list.ofKind(Gunk)
		for i := range gunks {
			distance := monster.Position().NewSub(gunks[i].Position()).SquareLength()
			if obj.Position().NewSub(gunks[i].Position()).SquareLength() < closestDistance {
				found = true
				closestIndex = i
				closestDistance = distance
			}
		}
		if found {
			if math.Sqrt(closestDistance) < monster.Scale[0] {
				gunks[closestIndex].SetState(stateDead)
				w.list.Remove(closestIndex)
				//monster.MaxAcceleration.Add(vector.NewVector3(1, 1, 1))
			}

			arrive := steering.NewSeek(monster.Model, monster.RigidBody, gunks[closestIndex].Position())
			st := arrive.Get()
			monster.AddForce(st.Linear())
			monster.AddTorque(st.Angular())
		} else {

		}
		monster.Position()[1] = monster.Scale[1] / 2
	}
}
