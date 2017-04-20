package core

import (
	"math"
	"time"

	"github.com/stojg/vector"
)

type Modality interface {
	Attenuation() float64
	// Check returns if the Sensor could detect this Signal, for example LOS checks, some might be always true
	Check(*Signal, Sensor) bool
	InverseTransmissionSpeed() float64
	MaximumRange() float64
}

type Signal struct {
	Modality Modality
	Position *vector.Vector3
	Strength float64
	ID       ID
}

type Sensor interface {
	// DetectsModality detects if this sensor can detect this modality, for example eye sensor - sight, ears sesnot - sound etc
	DetectsModality(Modality) bool
	Position() *vector.Vector3
	Orientation() *vector.Quaternion
	Threshold() float64
	Notify(*Signal)
}

type Notification struct {
	sensor    Sensor
	signal    *Signal
	timestamp time.Time
}

type Manager struct {
	sensors []Sensor
	// this could be a priority queue for efficiency
	notifications []*Notification
}

// Add adds a signal in the world
func (m *Manager) Add(signal *Signal) {
	// aggregation phase
	for _, sensor := range m.sensors {
		// testing phase

		// check the modality first
		if !sensor.DetectsModality(signal.Modality) {
			continue
		}

		// find the distance and check range
		distance := signal.Position.NewSub(sensor.Position()).Length()
		if signal.Modality.MaximumRange() < distance {
			continue
		}

		// find the intensity of the signal and check threshold
		intensity := signal.Strength * math.Pow(signal.Modality.Attenuation(), distance)
		if intensity < sensor.Threshold() {
			continue
		}

		// Perform additional modality specific checks
		if !signal.Modality.Check(signal, sensor) {
			continue
		}

		// we're going to notify the senors, work out when
		timestamp := time.Now().Add(time.Duration(distance*signal.Modality.InverseTransmissionSpeed()) * time.Second)

		// create notifications
		notification := &Notification{
			timestamp: timestamp,
			sensor:    sensor,
			signal:    signal,
		}

		// add to list
		m.notifications = append(m.notifications, notification)
	}

	// notifications phase
	m.SendSignals()
}

// SendSignals flushes notification from the queue. up to the current time
func (m *Manager) SendSignals() {
	now := time.Now()
	var new []*Notification
	for i, notification := range m.notifications {
		if now.After(notification.timestamp) {
			notification.sensor.Notify(notification.signal)
			new = append(m.notifications[:i], m.notifications[i+1:]...)
		}
	}
	m.notifications = new
}

func (m *Manager) Register(sensor Sensor) {
	m.sensors = append(m.sensors, sensor)
}

func (m *Manager) Deregister(sensor Sensor) {
	// @todo
}
