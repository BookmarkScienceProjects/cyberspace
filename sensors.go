package main

import (
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/cyberspace/lib/percepts"
	"github.com/stojg/vector"
)

func NewSight(object *core.GameObject, strength float64) *core.Signal {
	return &core.Signal{
		ID: object.ID(),
		Modality: &SightModality{
			invTransmissionSpeed: 0,
			maximumRange:         40,
			attenuation:          1,
			orientation:          object.Transform().Orientation(),
		},
		Position: object.Transform().Position(),
		Strength: strength,
	}
}

type SightModality struct {
	orientation  *vector.Quaternion
	maximumRange float64
	// attenuation describes how e.g. the volume of a sound drops over distance
	attenuation float64
	// how far will this signal travel over one second in per unit
	invTransmissionSpeed float64
}

func (m *SightModality) Attenuation() float64 {
	return m.attenuation
}

func (m *SightModality) InverseTransmissionSpeed() float64 {
	return m.invTransmissionSpeed
}

func (m *SightModality) MaximumRange() float64 {
	return m.maximumRange
}

func (mod *SightModality) Check(signal *core.Signal, sensor core.Sensor) bool {
	if !percepts.InSight(sensor.Position(), sensor.Orientation(), signal.Position, 3.14) {
		return false
	}

	res := percepts.InLineOfSight(sensor.Position(), signal.Position)
	// there should be only two objects in the line, the signal and sensor
	if len(res) != 2 {
		//fmt.Printf("not in line of sight %v\n", res)
		//for i := range res {
		//	fmt.Printf("%v\n", res[i])
		//}
		return false
	}

	return true
}
