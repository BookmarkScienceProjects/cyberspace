package core

import (
	"testing"

	"github.com/stojg/vector"
)

func TestManager_Add(t *testing.T) {

	list := NewObjectList()
	source := NewGameObject("source", list)

	sig := newTestSignal(1, source, 1)

	listener := &listener{}
	manager := NewManager()
	manager.Register(listener)

	manager.Add(sig)

	if len(manager.notifications) != 0 {
		t.Errorf("Expected that the manager.notifications list to be empty, instead %d\n", len(manager.notifications))
	}

	if len(listener.signals) != 1 {
		t.Error("expected to have been notified about a signal")
		return
	}

	recieved := listener.signals[0]
	if recieved.ID != 1 {
		t.Errorf("Expected signal.ID %d, got %d", 1, recieved.ID)
		return
	}

}

func BenchmarkManager_Add(b *testing.B) {
	list := NewObjectList()
	source := NewGameObject("source", list)

	sig := newTestSignal(1, source, 1)

	listener := &listener{}
	manager := NewManager()
	manager.Register(listener)

	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			manager.Add(sig)
		}
	}
}

type listener struct {
	signals []*Signal
}

func (l *listener) DetectsModality(modality Modality) bool {
	return true
}

func (l *listener) Notify(sig *Signal) {
	l.signals = append(l.signals, sig)
}

func (l *listener) Orientation() *vector.Quaternion {
	return vector.NewQuaternion(1, 0, 0, 0)
}

func (l *listener) Position() *vector.Vector3 {
	return vector.Zero()
}

func (l *listener) Threshold() float64 {
	return 0
}

func newTestSignal(sourceID ID, object *GameObject, strength float64) *Signal {
	return &Signal{
		ID: sourceID,
		Modality: &TestModality{
			invTransmissionSpeed: 0,
			maximumRange:         40,
			attenuation:          1,
			orientation:          object.Transform().Orientation(),
		},
		Position: object.Transform().Position(),
		Strength: strength,
	}
}

type TestModality struct {
	orientation  *vector.Quaternion
	maximumRange float64
	// attenuation describes how e.g. the volume of a sound drops over distance
	attenuation float64
	// how far will this signal travel over one second in per unit
	invTransmissionSpeed float64
}

func (m *TestModality) Attenuation() float64 {
	return m.attenuation
}

func (m *TestModality) InverseTransmissionSpeed() float64 {
	return m.invTransmissionSpeed
}

func (m *TestModality) MaximumRange() float64 {
	return m.maximumRange
}

func (mod *TestModality) Check(signal *Signal, sensor Sensor) bool {
	return true
}
