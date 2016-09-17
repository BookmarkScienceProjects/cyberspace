package main

import (
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/goap"
	"github.com/stojg/vector"
	"math/rand"
)

func newLevel() *level {
	lvl := &level{
		worldState: make(goap.StateList, 0),
	}
	for i := 0; i < 10; i++ {
		obj := spawn("monster")
		obj.AddAI(&HunterAI{})
		obj.Transform().Position().Set(rand.Float64()*50-24, 0, rand.Float64()*50-25)
		lvl.worldState["monster_exists"] = true
	}
	for i := 0; i < 50; i++ {
		obj := spawn("food")
		obj.Transform().Position().Set(rand.Float64()*50-24, 0, rand.Float64()*50-25)
		lvl.worldState["food_exists"] = true
	}
	return lvl
}

// level contains the overall logic for populating and updating game objects to the level
type level struct {
	worldState goap.StateList
}

// Update gets called from the main game loop
func (l *level) Update(elapsed float64) {
	for _, obj := range core.List.FindWithTag("monster") {
		obj.Body().AddTorque(vector.NewVector3(0, 1, 0))
	}

	for _, obj := range core.List.FindWithTag("food") {
		obj.Body().AddTorque(vector.NewVector3(0, -1, 0))
	}

	UpdateAI(elapsed, l.worldState)
	UpdatePhysics(elapsed)
	UpdateCollisions(elapsed)

}
