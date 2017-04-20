package core

import (
	"container/heap"
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
	index     int
	sensor    Sensor
	signal    *Signal
	timestamp time.Time
}

type NotificationList []*Notification

func (pq NotificationList) Len() int { return len(pq) }

func (pq NotificationList) Less(i, j int) bool {
	return pq[i].timestamp.Before(pq[j].timestamp)
}

func (pq NotificationList) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *NotificationList) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Notification)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *NotificationList) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *NotificationList) Peek() *Notification {
	if len(*pq) > 0 {
		return (*pq)[0]
	}
	return nil
}

func NewManager() *Manager {
	m := &Manager{}
	m.notifications = make(NotificationList, 0)
	heap.Init(&m.notifications)
	return m
}

type Manager struct {
	sensors []Sensor
	// this could be a priority queue for efficiency
	notifications NotificationList
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

		// create signals
		notification := &Notification{
			timestamp: timestamp,
			sensor:    sensor,
			signal:    signal,
		}
		heap.Push(&m.notifications, notification)
	}

	// signals phase
	m.SendSignals()
}

// SendSignals flushes notification from the queue. up to the current time
func (m *Manager) SendSignals() {
	now := time.Now()
	for m.notifications.Len() > 0 {
		peek := m.notifications.Peek()
		if now.Before(peek.timestamp) {
			return
		}
		note := heap.Pop(&m.notifications).(*Notification)
		note.sensor.Notify(note.signal)
	}
}

func (m *Manager) Register(sensor Sensor) {
	m.sensors = append(m.sensors, sensor)
}

func (m *Manager) Deregister(sensor Sensor) {
	// @todo
}
