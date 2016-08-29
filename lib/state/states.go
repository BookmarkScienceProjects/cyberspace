package state

import (
	. "github.com/stojg/cyberspace/lib/components"
	. "github.com/stojg/steering"
	. "github.com/stojg/vector"
	. "github.com/stojg/vivere/lib/components"
)

type State interface {
	Steering() *SteeringOutput
	Update() State
}

func siblings(instance *AWSInstance, model *Model, inclusive bool) []*Vector3 {
	var positions []*Vector3
	for _, i := range instance.Tree.Siblings(instance.Name) {
		if i == instance && !inclusive {
			continue
		}
		positions = append(positions, i.Position())
	}
	return positions
}

func GetMidpoint(p []*Vector3) *Vector3 {
	if len(p) < 2 {
		return nil
	}
	midpoint := &Vector3{}
	for _, i := range p {
		midpoint.Add(i)
	}
	midpoint.Scale(1 / float64(len(p)))
	return midpoint

}
