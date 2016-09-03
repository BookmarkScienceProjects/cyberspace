package main

import "github.com/stojg/steering"

type controllerSystem struct{}

func (s *controllerSystem) Update(elapsed float64) {

	for _, instance := range instanceList.All() {

		if instance.Target() != nil {
			arrive := steering.NewArrive(instance.Model, instance.RigidBody, instance.Target().Position(), 200, 0.1, 20)
			steering := arrive.Get()
			instance.RigidBody.AddForce(steering.Linear())
		}
		instance.Model.Position()[1] = instance.Model.Scale[1] / 2

	}

}
