package main

import (
	"fmt"
	"github.com/stojg/vivere/lib/components"
	"math/rand"
)

type game struct {
	timer float64
	list  []*Object
}

type Object struct {
	ID *components.Entity
	*components.Model
	*components.RigidBody
	*components.Collision
}


func (g *game) Update(elapsed float64) {

	g.timer += elapsed

	if g.timer < 1 {
		return
	}
	g.timer -= 1
	if len(g.list) <= 10 {
		fmt.Println("create monster")
		g.list = append(g.list, createMonster(g, 10, 10, 10))
	}


}

func createMonster(g *game, width, height, depth float64) *Object {

	eID := entities.Create()
	monster := &Object{
		ID:        eID,
		Model:     modelList.New(eID, width, height, depth, 2),
		RigidBody: rigidList.New(eID, 1),
		Collision: collisionList.New(eID, width, height, depth),
	}
	monster.Position().Set(rand.Float64()*500-250, 5, rand.Float64()*500-250)
	return monster
}
