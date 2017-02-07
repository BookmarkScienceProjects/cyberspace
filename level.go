package main

import (
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/goap"
	"math/rand"
)

func newLevel() *level {
	lvl := &level{
		worldState: make(goap.StateList),
	}
	for i := 0; i < 5; i++ {
		obj := spawn("monster")
		obj.AddAgent(NewMonsterAgent())
		obj.Transform().Position().Set(rand.Float64()*50-24, 0, rand.Float64()*50-25)
		lvl.worldState["monster_exists"] = true
	}
	return lvl
}

// level contains the overall logic for populating and updating game objects to the level
type level struct {
	worldState goap.StateList
}

// Update gets called from the main game loop
func (l *level) Update(elapsed float64) {
	UpdateAI(elapsed, l.worldState)
	UpdatePhysics(elapsed)
	UpdateCollisions(elapsed)

	if len(core.List.FindWithTag("food")) < 40 {
		obj := spawn("food")
		obj.Transform().Position().Set(rand.Float64()*50-25, 0, rand.Float64()*50-25)
		l.worldState["food_exists"] = true
	}

}
