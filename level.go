package main

import "github.com/stojg/cyberspace/lib/plan"

func newLevel() *level {
	lvl := &level{
		State: make(plan.StateList),
	}
	return lvl
}

// level contains the overall logic for populating and updating game objects to the level
type level struct {
	State plan.StateList
}
