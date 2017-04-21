package actions

import "github.com/stojg/cyberspace/lib/planning"

var (
	EnemyInSight  = planning.State{Name: "enemy_in_sight", Value: true}
	AreaPatrolled = planning.State{Name: "area_patrolled", Value: true}
	EnemyKilled   = planning.State{Name: "enemy_killed", Value: true}
	Healthy       = planning.State{Name: "healthy", Value: true}
)
