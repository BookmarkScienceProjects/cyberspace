package main

import (
	"github.com/stojg/cyberspace/lib/actions"
	"github.com/stojg/cyberspace/lib/object"
	"github.com/stojg/goap"
	"github.com/stojg/vivere/lib/components"
	"math/rand"
)

func createMonster(name string) *GameObject {
	obj := loadFromFile(name)
	if obj == nil {
		return nil
	}
	eID := entities.Create()

	eat := actions.NewEat(4)
	getFood := actions.NewGetFood(5)
	actions := []goap.Actionable{eat, getFood}

	monster := &GameObject{
		id:        eID,
		kind:      obj.Kind,
		Model:     modelList.New(eID, obj.Scale[0], obj.Scale[1], obj.Scale[2], components.EntityType(obj.Kind)),
		RigidBody: rigidList.New(eID, obj.Weight),
		Collision: collisionList.New(eID, obj.Scale[0], obj.Scale[1], obj.Scale[2]),
	}
	monster.agent = goap.NewGoapAgent(monster, actions)
	monster.state = make(goap.StateList, 0)
	monster.state["is_full"] = false
	monster.state["has_food"] = false

	return monster
}

func createFood(name string) *GameObject {
	obj := loadFromFile(name)
	if obj == nil {
		return nil
	}
	eID := entities.Create()

	eat := actions.NewEat(4)
	getFood := actions.NewGetFood(5)
	actions := []goap.Actionable{eat, getFood}

	food := &GameObject{
		id:        eID,
		kind:      obj.Kind,
		Model:     modelList.New(eID, obj.Scale[0], obj.Scale[1], obj.Scale[2], components.EntityType(obj.Kind)),
		RigidBody: rigidList.New(eID, obj.Weight),
		Collision: collisionList.New(eID, obj.Scale[0], obj.Scale[1], obj.Scale[2]),
	}
	food.agent = goap.NewGoapAgent(food, actions)
	food.state = make(goap.StateList, 0)
	return food
}

type World struct {
	timer float64
	state goap.StateList
}

func (w *World) Add(o object.Object) {
	object.List.Add(o)
}

func (w *World) Remove(i uint16) {
	object.List.Remove(i)
}

func (w *World) Update(elapsed float64) {

	// update world state
	w.timer += elapsed
	if w.timer > 10 {
		w.timer -= 10
	}

	w.CreateEntities()

	// run the AI for the individual entities
	for _, obj := range object.List.OfKind(object.Monster) {
		monster := obj.(*GameObject)
		monster.Update()
		//
		//id, found := nearest(monster, w.objects.ofKind(Food))
		//if found {
		//	food := w.objects.items[id]
		//	dir := monster.Position().NewSub(food.Position())
		//	// monster is close enough, kill food
		//	if dir.Length() < monster.Scale[0] {
		//		w.Remove(id)
		//	}
		//	arrive := steering.NewSeek(monster.Model, monster.RigidBody, food.Position())
		//	st := arrive.Get()
		//	monster.AddForce(st.Linear())
		//	monster.AddTorque(st.Angular())
		//}
	}
}

func (w *World) CreateEntities() {
	for len(object.List.OfKind(object.Food)) < 200 {
		food := createFood("food")
		if food != nil {
			food.Position().Set(rand.Float64()*1000-500, 0, rand.Float64()*1000-500)
			w.Add(food)
		}
	}

	for len(object.List.OfKind(object.Monster)) < 10 {
		monster := createMonster("monster")
		if monster != nil {
			monster.Position().Set(rand.Float64()*500-250, 0, rand.Float64()*500-250)
			w.Add(monster)
		}
	}
}
