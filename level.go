package main

import "github.com/stojg/cyberspace/lib/planning"

func newLevel() *level {
	lvl := &level{
		State: make(planning.StateList),
	}
	return lvl
}

// level contains the overall logic for populating and updating game objects to the level
type level struct {
	State planning.StateList
}
