package core

import (
	"github.com/stojg/vivere/lib/components"
	"math"
	"sync"
)

var List *ObjectList

func init() {
	List = &ObjectList{
		entities:   make(map[components.Entity]*GameObject),
		graphics:   make(map[components.Entity]*Graphic),
		bodies:     make(map[components.Entity]*Body),
		collisions: make(map[components.Entity]*Collision),
		agents:     make(map[components.Entity]*Agent),
		deleted:    make([]components.Entity, 0),
	}
}

type ObjectList struct {
	sync.Mutex
	nextID     components.Entity
	entities   map[components.Entity]*GameObject
	graphics   map[components.Entity]*Graphic
	bodies     map[components.Entity]*Body
	collisions map[components.Entity]*Collision
	agents     map[components.Entity]*Agent
	deleted    []components.Entity
}

func (l *ObjectList) Add(g *GameObject) {
	l.Lock()
	defer l.Unlock()
	if l.nextID == math.MaxUint32 {
		panic("Out of entity ids, implement GC")
	}
	l.nextID++
	g.id = l.nextID
	l.entities[g.id] = g
}

func (l *ObjectList) Remove(g *GameObject) {
	l.Lock()
	defer l.Unlock()
	if _, found := l.graphics[g.id]; found {
		delete(l.graphics, g.id)
	}
	if _, found := l.bodies[g.id]; found {
		delete(l.bodies, g.id)
	}
	if _, found := l.collisions[g.id]; found {
		delete(l.bodies, g.id)
	}
	if _, found := l.agents[g.id]; found {
		delete(l.bodies, g.id)
	}
	if _, found := l.entities[g.id]; found {
		delete(l.entities, g.id)
	}
	l.deleted = append(l.deleted, g.id)
}

func (l *ObjectList) All() []*GameObject {
	l.Lock()
	defer l.Unlock()
	var result []*GameObject
	for i := range l.entities {
		result = append(result, l.entities[i])
	}
	return result
}

// Returns GameObjects tagged tag. Returns nil if no GameObjects was found.
func (l *ObjectList) FindWithTag(tag string) []*GameObject {
	l.Lock()
	defer l.Unlock()
	var result []*GameObject
	for i := range l.entities {
		if l.entities[i].CompareTag(tag) {
			result = append(result, l.entities[i])
		}
	}
	return result
}

func (l *ObjectList) AddGraphic(id components.Entity, graphic *Graphic) {
	l.Lock()
	defer l.Unlock()
	graphic.gameObject = l.entities[id]
	graphic.transform = l.entities[id].transform
	l.graphics[id] = graphic
}

func (l *ObjectList) Graphic(id components.Entity) *Graphic {
	l.Lock()
	defer l.Unlock()
	return l.graphics[id]
}

func (l *ObjectList) Graphics() []*Graphic {
	l.Lock()
	defer l.Unlock()
	var result []*Graphic
	for i := range l.graphics {
		result = append(result, l.graphics[i])
	}
	return result
}

func (l *ObjectList) AddBody(id components.Entity, body *Body) {
	l.Lock()
	defer l.Unlock()
	body.gameObject = l.entities[id]
	body.transform = l.entities[id].transform
	l.bodies[id] = body
}

func (l *ObjectList) Body(id components.Entity) *Body {
	l.Lock()
	defer l.Unlock()
	return l.bodies[id]
}

func (l *ObjectList) Bodies() []*Body {
	l.Lock()
	defer l.Unlock()
	var result []*Body
	for i := range l.bodies {
		result = append(result, l.bodies[i])
	}
	return result
}

func (l *ObjectList) AddCollision(id components.Entity, collision *Collision) {
	l.Lock()
	defer l.Unlock()
	collision.gameObject = l.entities[id]
	collision.transform = l.entities[id].transform
	l.collisions[id] = collision
}

func (l *ObjectList) Collision(id components.Entity) *Collision {
	l.Lock()
	defer l.Unlock()
	return l.collisions[id]
}

func (l *ObjectList) Collisions() []*Collision {
	l.Lock()
	defer l.Unlock()
	var result []*Collision
	for i := range l.collisions {
		result = append(result, l.collisions[i])
	}
	return result
}

func (l *ObjectList) AddAI(id components.Entity, agent *Agent) {
	l.Lock()
	defer l.Unlock()
	agent.gameObject = l.entities[id]
	l.agents[id] = agent
}

func (l *ObjectList) Agent(id components.Entity) *Agent {
	l.Lock()
	defer l.Unlock()
	return l.agents[id]
}

func (l *ObjectList) Agents() []*Agent {
	l.Lock()
	defer l.Unlock()
	var result []*Agent
	for i := range l.agents {
		result = append(result, l.agents[i])
	}
	return result
}

func (l *ObjectList) Deleted() []components.Entity {
	l.Lock()
	defer l.Unlock()
	return l.deleted
}

func (l *ObjectList) ClearDeleted() {
	l.Lock()
	defer l.Unlock()
	l.deleted = make([]components.Entity, 0)
}
