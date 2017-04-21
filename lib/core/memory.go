package core

import (
	"time"

	"github.com/stojg/vector"
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
		entities: make(map[ID]*Entity),
		internal: &Internal{},
	}
}

type WorkingMemory struct {
	entities map[ID]*Entity
	internal *Internal
}

func (memory *WorkingMemory) Entities() []*Entity {
	var res []*Entity
	for _, v := range memory.entities {
		res = append(res, v)
	}
	return res
}

func (memory *WorkingMemory) Internal() *Internal {
	return memory.internal
}

func (memory *WorkingMemory) tick() {
	now := time.Now()
	for id, value := range memory.entities {
		if value.expiry.After(now) {
			delete(memory.entities, id)
		}
	}
}

// AddEntity adds an entity to the memory, returns true if its an update
func (memory *WorkingMemory) AddEntity(f *Entity) bool {
	_, found := memory.entities[f.ID]
	if found {
		memory.entities[f.ID] = f
		return true
	}
	memory.entities[f.ID] = f
	return false
}
