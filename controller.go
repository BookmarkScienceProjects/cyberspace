package main

type controllerSystem struct{}

func (s *controllerSystem) Update(elapsed float64) {

	for _, move := range controllerList.All() {
		move.Update(elapsed)
	}

}
