package core

import (
	"github.com/stojg/vivere/lib/components"
	"math"
)

var List *ObjectList

func init() {
	List = &ObjectList{
		entities:   make(map[components.Entity]*GameObject),
		graphics:   make(map[components.Entity]*Graphic),
		bodies:     make(map[components.Entity]*Body),
		collisions: make(map[components.Entity]*Collision),
		ais:        make(map[components.Entity]AI),
		deleted:    make([]components.Entity, 0),
	}
}

type ObjectList struct {
	nextID     components.Entity
	entities   map[components.Entity]*GameObject
	graphics   map[components.Entity]*Graphic
	bodies     map[components.Entity]*Body
	collisions map[components.Entity]*Collision
	ais        map[components.Entity]AI
	deleted    []components.Entity
}

func (l *ObjectList) Add(g *GameObject) {
	if l.nextID == math.MaxUint32 {
		panic("Out of entity ids, implement GC")
	}
	l.nextID++
	g.id = l.nextID
	l.entities[g.id] = g
}

func (l *ObjectList) Remove(g *GameObject) {
	if _, found := l.entities[g.id]; found {
		delete(l.entities, g.id)
	}
	if _, found := l.graphics[g.id]; found {
		delete(l.graphics, g.id)
	}
	if _, found := l.bodies[g.id]; found {
		delete(l.bodies, g.id)
	}
	l.deleted = append(l.deleted, g.id)
}

func (l *ObjectList) All() []*GameObject {
	var result []*GameObject
	for i := range l.entities {
		result = append(result, l.entities[i])
	}
	return result
}

// Returns one active GameObject tagged tag. Returns nil if no GameObject was found.
func (l *ObjectList) FindWithTag(tag string) []*GameObject {
	var result []*GameObject
	for i := range l.entities {
		if l.entities[i].CompareTag(tag) {
			result = append(result, l.entities[i])
		}
	}
	return result
}

func (l *ObjectList) AddGraphic(id components.Entity, graphic *Graphic) {
	graphic.gameObject = l.entities[id]
	graphic.transform = l.entities[id].transform
	l.graphics[id] = graphic
}

func (l *ObjectList) Graphic(id components.Entity) *Graphic {
	return l.graphics[id]
}

func (l *ObjectList) Graphics() []*Graphic {
	var result []*Graphic
	for i := range l.graphics {
		result = append(result, l.graphics[i])
	}
	return result
}

func (l *ObjectList) AddBody(id components.Entity, body *Body) {
	body.gameObject = l.entities[id]
	body.transform = l.entities[id].transform
	l.bodies[id] = body
}

func (l *ObjectList) Body(id components.Entity) *Body {
	return l.bodies[id]
}

func (l *ObjectList) Bodies() []*Body {
	var result []*Body
	for i := range l.bodies {
		result = append(result, l.bodies[i])
	}
	return result
}

func (l *ObjectList) AddCollision(id components.Entity, collision *Collision) {
	collision.gameObject = l.entities[id]
	collision.transform = l.entities[id].transform
	l.collisions[id] = collision
}

func (l *ObjectList) Collision(id components.Entity) *Collision {
	return l.collisions[id]
}

func (l *ObjectList) Collisions() []*Collision {
	var result []*Collision
	for i := range l.collisions {
		result = append(result, l.collisions[i])
	}
	return result
}

func (l *ObjectList) AddAI(id components.Entity, ai AI) {
	ai.SetGameObject(l.entities[id])
	l.ais[id] = ai
}

func (l *ObjectList) AI(id components.Entity) AI {
	return l.ais[id]
}

func (l *ObjectList) AIs() []AI {
	var result []AI
	for i := range l.ais {
		result = append(result, l.ais[i])
	}
	return result
}

// Finds a game object by name and returns it.
func (l *ObjectList) Find() *GameObject {
	return nil
}

// Returns a list of active GameObjects tagged tag. Returns empty array if no GameObject was found.
func (l *ObjectList) FindGameObjectsWithTag() []*GameObject {
	return nil
}

func (l *ObjectList) Deleted() []components.Entity {
	return l.deleted
}

func (l *ObjectList) ClearDeleted() {
	l.deleted = make([]components.Entity, 0)
}
