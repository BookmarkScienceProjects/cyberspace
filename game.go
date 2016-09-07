package main

import (
	"github.com/stojg/steering"
	"math"
	"math/rand"
	"sync"
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

	// keep a record of deleted entities so they can be flushed to the view
	l.deleted = append(l.deleted, l.items[i])
	// Take the last element in the list and replace the object to be deleted with that one
	l.items[i] = l.items[l.next-1]
	// we now have one more spot
	l.next--

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

func createFood() *GameObject {
	eID := entities.Create()
	f := &GameObject{
		id:        eID,
		state:     stateIdle,
		kind:      Gunk,
		Model:     modelList.New(eID, 3, 3, 3, 2),
		RigidBody: rigidList.New(eID, 10),
		Collision: collisionList.New(eID, 3, 3, 3),
	}
	return f
}

func createMonster() *GameObject {
	eID := entities.Create()
	monster := &GameObject{
		id:        eID,
		kind:      Monster,
		Model:     modelList.New(eID, 10, 10, 10, 1),
		RigidBody: rigidList.New(eID, 1),
		Collision: collisionList.New(eID, 10, 10, 10),
	}
	return monster
}

type World struct {
	timer float64
	list  *stuffList
}

func (w *World) Add(o Object) {
	w.list.Add(o)
}

func (w *World) Remove(i uint16) {
	w.list.Remove(i)
}

func (w *World) Update(elapsed float64) {

	w.timer += elapsed
	if w.timer > 10 {
		w.timer -= 10
	}

	for len(w.list.ofKind(Gunk)) < 2 {
		food := createFood()
		food.Position().Set(rand.Float64()*1000-500, food.Scale[1]/2, rand.Float64()*1000-500)
		w.Add(food)
	}

	for len(w.list.ofKind(Monster)) < 2 {
		monster := createMonster()
		monster.state = stateIdle
		monster.Position().Set(rand.Float64()*500-250, monster.Scale[1]/2, rand.Float64()*500-250)
		w.Add(monster)
	}

	// run the AI for the individual entities
	for _, obj := range w.list.ofKind(Monster) {
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
				w.Remove(closestIndex)
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
