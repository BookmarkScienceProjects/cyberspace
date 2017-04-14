package actions

import "github.com/stojg/goap"

var (
	HasFood = goap.State{Name: "has_food", Value: true}
	Full = goap.State{Name: "is_full", Value: true}
	Rested = goap.State{Name: "is_rested", Value: true}
)
