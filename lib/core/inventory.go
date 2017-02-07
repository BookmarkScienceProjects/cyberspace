package core

import (
	"github.com/stojg/vivere/lib/components"
)

// NewInventory returns a new Inventory component
func NewInventory() *Inventory {
	return &Inventory{
		list: make(map[string]int),
	}
}

// Graphic component describes various properties on how it's related GameObject should be rendered
type Inventory struct {
	Component
	list map[string]int
}

func (i *Inventory) Add(name string, num int) {
	if _, ok := i.list[name]; !ok {
		i.list[name] = num
	} else {
		i.list[name] += num
	}
}

func (i *Inventory) Remove(name string, num int) {
	if _, ok := i.list[name]; !ok {
		return
	}
	i.list[name] -= num
	if i.list[name] < 0 {
		i.list[name] = 0
	}
}

func (l *ObjectList) AddInventory(id components.Entity, inventory *Inventory) {
	l.Lock()
	inventory.gameObject = l.entities[id]
	l.inventories[id] = inventory
	l.Unlock()
}

// Inventory returns the inventory component for a GameObject
func (l *ObjectList) Inventory(id components.Entity) *Inventory {
	return l.inventories[id]
}

// Inventories returns all registered inventory components
func (l *ObjectList) Inventories() []*Inventory {
	l.Lock()
	var result []*Inventory
	for i := range l.inventories {
		result = append(result, l.inventories[i])
	}
	l.Unlock()
	return result
}
