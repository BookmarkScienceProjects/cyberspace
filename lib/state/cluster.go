package state

import (
	. "github.com/stojg/cyberspace/lib/components"
	. "github.com/stojg/steering"
	"github.com/stojg/vector"
	. "github.com/stojg/vivere/lib/components"
)

func NewCluster(e *Entity, i *AWSInstance, m *Model, r *RigidBody, target *vector.Vector3) *Cluster {
	return &Cluster{
		entity:   e,
		instance: i,
		model:    m,
		body:     r,
		target:   target,
	}
}

type Cluster struct {
	entity       *Entity
	instance     *AWSInstance
	body         *RigidBody
	model        *Model
	target       *vector.Vector3
	prevSteering float64
}

func (s *Cluster) Steering() *SteeringOutput {

	steering := NewSteeringOutput()

	arrive := NewArrive(s.model, s.body, s.target, 100, s.model.Scale[0]*1, s.model.Scale[0]*4).Get()
	steering.Linear().Add(arrive.Linear())

	targets := siblings(s.instance, s.model, false)
	if len(targets) > 0 {
		separation := NewSeparation(s.model, s.body, targets, s.model.Scale[0]*2).Get()
		steering.Linear().Add(separation.Linear())
	}

	s.prevSteering = steering.Linear().Length()

	return steering
}

func (s *Cluster) Update() State {
	if s.prevSteering < 30 {
		//fmt.Printf(" %s switching to idle (%0.2f)\n", s.instance.Name, s.prevSteering)
		return NewIdle(s.entity, s.instance, s.model, s.body)
	}
	return nil
}