package core

import (
	"github.com/stojg/vector"
	"time"
)

type EntityType int

const (
	Creature EntityType = iota
)

type Entity struct {
	ID       ID
	Distance float64
	Position *vector.Vector3
	Velocity *vector.Vector3
	Name     string
	Type     EntityType
	expiry   time.Time
}

type Internal struct {
	Health float64
	Energy float64
	Food   float64
}

func NewWorkingMemory() *WorkingMemory {
	return &WorkingMemory{
		entities: make([]*Entity, 0),
		internal: &Internal{},
	}
}

type WorkingMemory struct {
	entities []*Entity
	internal *Internal
}

func (memory *WorkingMemory) Entities() []*Entity {
	return memory.entities
}

func (memory *WorkingMemory) Internal() *Internal {
	return memory.internal
}

func (memory *WorkingMemory) tick() {
	now := time.Now()
	for i := len(memory.entities) - 1; i >= 0; i-- {
		if now.After(memory.entities[i].expiry) {
			memory.entities = append(memory.entities[:i], memory.entities[i+1:]...)
		}
	}
}

func (memory *WorkingMemory) AddEntity(f *Entity) {
	for _, ent := range memory.entities {
		if ent.ID == f.ID {
			ent.expiry = time.Now()
			ent.Position = f.Position
			ent.Type = f.Type
			return
		}
	}
	memory.entities = append(memory.entities, f)
}
