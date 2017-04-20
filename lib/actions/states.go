package actions

import "github.com/stojg/goap"

var (
	EnemyInSight  = goap.State{Name: "enemy_in_sight", Value: true}
	AreaPatrolled = goap.State{Name: "area_patrolled", Value: true}
	EnemyKilled   = goap.State{Name: "enemy_killed", Value: true}
	Healthy       = goap.State{Name: "healthy", Value: true}
)
