package core

import (
	"github.com/stojg/goap"
	"github.com/stojg/vector"
)

type FactType int

const (
	Creature FactType = iota
	Item
)

type WorkingMemory struct {
	data []*WorkingMemoryFact
}

func (memory *WorkingMemory) Data() []*WorkingMemoryFact {
	return memory.data
}

func (memory *WorkingMemory) Tick() {
	for i := len(memory.data) - 1; i >= 0; i-- {
		memory.data[i].Confidence -= 0.1
		if memory.data[i].Confidence < 0 {
			memory.data = append(memory.data[:i], memory.data[i+1:]...)
		}
	}
}

func (memory *WorkingMemory) Add(f *WorkingMemoryFact) {
	for i := range memory.data {
		if memory.data[i].ID == f.ID {
			memory.data[i].Confidence = f.Confidence
			memory.data[i].Position = f.Position
			memory.data[i].Type = f.Type
			return
		}
	}
	memory.data = append(memory.data, f)
}

// WorkingMemoryFact examples from http://alumni.media.mit.edu/~jorkin/aiide05OrkinJ.pdf
type WorkingMemoryFact struct {
	Position *vector.Vector3
	//Direction *vector.Vector3
	//Stimulus  *vector.Vector3
	ID ID
	//Desire
	Confidence float64
	Type       FactType
	States     []goap.State
}
