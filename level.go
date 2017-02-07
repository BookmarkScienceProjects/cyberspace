package main

import (
	"github.com/stojg/goap"
)

func newLevel() *level {
	lvl := &level{
		State: make(goap.StateList),
	}
	return lvl
}

// level contains the overall logic for populating and updating game objects to the level
type level struct {
	State goap.StateList
}
