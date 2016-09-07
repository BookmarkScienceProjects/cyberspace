package main

import (
	"github.com/stojg/steering"
	"github.com/stojg/vivere/lib/components"
	"math/rand"
)

func create(name string) *GameObject {
	obj := loadFromFile(name)
	if obj == nil {
		return nil
	}
	eID := entities.Create()
	monster := &GameObject{
		id:        eID,
		state:     stateIdle,
		kind:      obj.Kind,
		Model:     modelList.New(eID, obj.Scale[0], obj.Scale[1], obj.Scale[2], components.EntityType(obj.Kind)),
		RigidBody: rigidList.New(eID, obj.Weight),
		Collision: collisionList.New(eID, obj.Scale[0], obj.Scale[1], obj.Scale[2]),
	}
	return monster
}

type World struct {
	timer   float64
	objects *objectList
}

func (w *World) Add(o Object) {
	w.objects.Add(o)
}

func (w *World) Remove(i uint16) {
	w.objects.Remove(i)
}

func (w *World) Update(elapsed float64) {

	w.timer += elapsed
	if w.timer > 10 {
		w.timer -= 10
	}

	w.CreateEntities()

	// run the AI for the individual entities
	for _, obj := range w.objects.ofKind(Monster) {
		monster := obj.(*GameObject)

		id, found := nearest(monster, w.objects.ofKind(Food))
		if found {
			food := w.objects.items[id]
			dir := monster.Position().NewSub(food.Position())
			// monster is close enough, kill food
			if dir.Length() < monster.Scale[0] {
				w.Remove(id)
			}
			arrive := steering.NewSeek(monster.Model, monster.RigidBody, food.Position())
			st := arrive.Get()
			monster.AddForce(st.Linear())
			monster.AddTorque(st.Angular())
		}
	}
}

func (w *World) CreateEntities() {
	for len(w.objects.ofKind(Food)) < 500 {
		food := create("food")
		if food != nil {
			food.Position().Set(rand.Float64()*1000-500, 0, rand.Float64()*1000-500)
			w.Add(food)
		}
	}

	for len(w.objects.ofKind(Monster)) < 10 {
		monster := create("monster")
		if monster != nil {
			monster.Position().Set(rand.Float64()*500-250, 0, rand.Float64()*500-250)
			w.Add(monster)
		}
	}
}
