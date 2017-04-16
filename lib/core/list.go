package core

import (
	"github.com/stojg/cyberspace/lib/quadtree"
	"math"
	"sync"
)

// List is the primary resource for adding, removing and changing GameObjects and their components.
var List *ObjectList

func init() {
	List = NewObjectList()
}

func NewObjectList() *ObjectList {
	return &ObjectList{
		entities:    make(map[Entity]*GameObject),
		graphics:    make(map[Entity]*Graphic),
		bodies:      make(map[Entity]*Body),
		collisions:  make(map[Entity]*Collision),
		agents:      make(map[Entity]*Agent),
		inventories: make(map[Entity]*Inventory),
		deleted:     make([]Entity, 0),
	}
}

// ObjectList is a struct that contains a list of GameObjects and their components. All creation,
// removal and changes should be handled by this list so they don't get lost or out of sync.
type ObjectList struct {
	sync.Mutex
	quadTree    *quadtree.QuadTree
	nextID      Entity
	entities    map[Entity]*GameObject
	graphics    map[Entity]*Graphic
	bodies      map[Entity]*Body
	collisions  map[Entity]*Collision
	agents      map[Entity]*Agent
	inventories map[Entity]*Inventory
	deleted     []Entity
}

// Add a GameObject to this list and assign it an unique ID
func (l *ObjectList) Add(g *GameObject) {
	l.Lock()
	defer l.Unlock()
	if l.nextID == math.MaxUint32 {
		panic("Out of entity ids, implement GC")
	}
	l.nextID++
	g.id = l.nextID
	l.entities[g.id] = g

	for _, a := range l.agents {
		a.Replan()
	}
}

func (l *ObjectList) Get(id Entity) *GameObject {
	l.Lock()
	defer l.Unlock()
	return l.entities[id]
}

func (l *ObjectList) BuildQuadTree() {
	qTree := quadtree.NewQuadTree(
		quadtree.BoundingBox{
			MinX: -3200 / 2,
			MaxX: 3200 / 2,
			MinY: -3200 / 2,
			MaxY: 3200 / 2,
		},
	)
	// build the quad tree for broad phase collision detection
	for _, collision := range l.Collisions() {
		qTree.Add(collision)
	}

	l.quadTree = &qTree
}

func (l *ObjectList) QuadTree() *quadtree.QuadTree {
	return l.quadTree
}

// Remove a GameObject and all of it's components
func (l *ObjectList) Remove(g *GameObject) {
	l.Lock()
	if _, found := l.graphics[g.id]; found {
		delete(l.graphics, g.id)
	}
	if _, found := l.bodies[g.id]; found {
		delete(l.bodies, g.id)
	}
	if _, found := l.agents[g.id]; found {
		delete(l.agents, g.id)
	}
	if _, found := l.collisions[g.id]; found {
		delete(l.collisions, g.id)
	}
	if _, found := l.entities[g.id]; found {
		delete(l.entities, g.id)
	}
	if _, found := l.inventories[g.id]; found {
		delete(l.inventories, g.id)
	}
	l.deleted = append(l.deleted, g.id)
	delete(l.entities, g.id)
	l.Unlock()
	for _, a := range l.agents {
		a.Replan()
	}
}

// All returns all GameObjects in this list
func (l *ObjectList) All() []*GameObject {
	l.Lock()
	var result []*GameObject
	for i := range l.entities {
		result = append(result, l.entities[i])
	}
	l.Unlock()
	return result
}

// FindWithTag returns all GameObjects tagged with tag.
func (l *ObjectList) FindWithTag(tag string) []*GameObject {
	l.Lock()
	var result []*GameObject
	for i := range l.entities {
		if l.entities[i].CompareTag(tag) {
			result = append(result, l.entities[i])
		}
	}
	l.Unlock()
	return result
}

// AddGraphic adds a Graphic component to a GameObject
func (l *ObjectList) AddGraphic(id Entity, graphic *Graphic) {
	l.Lock()
	graphic.gameObject = l.entities[id]
	graphic.transform = l.entities[id].transform
	l.graphics[id] = graphic
	l.Unlock()
}

// Graphic returns the Graphic component for a GameObject
func (l *ObjectList) Graphic(id Entity) *Graphic {
	return l.graphics[id]
}

// Graphics returns all Graphic components
func (l *ObjectList) Graphics() []*Graphic {
	l.Lock()
	var result []*Graphic
	for i := range l.graphics {
		result = append(result, l.graphics[i])
	}
	l.Unlock()
	return result
}

// AddBody adds a Body component to a GameObject
func (l *ObjectList) AddBody(id Entity, body *Body) {
	l.Lock()
	body.gameObject = l.entities[id]
	body.transform = l.entities[id].transform
	l.bodies[id] = body
	l.Unlock()
}

// Body returns the Body component for a GameObject
func (l *ObjectList) Body(id Entity) *Body {
	return l.bodies[id]
}

// Bodies returns all Body components
func (l *ObjectList) Bodies() []*Body {
	l.Lock()
	var result []*Body
	for i := range l.bodies {
		result = append(result, l.bodies[i])
	}
	l.Unlock()
	return result
}

// AddCollision adds a Collision component to a GameObject
func (l *ObjectList) AddCollision(id Entity, collision *Collision) {
	l.Lock()
	collision.gameObject = l.entities[id]
	collision.transform = l.entities[id].transform
	l.collisions[id] = collision
	l.Unlock()
}

// Collision returns the Collision component for a GameObject
func (l *ObjectList) Collision(id Entity) *Collision {
	return l.collisions[id]
}

// Collisions returns all registered Collision components
func (l *ObjectList) Collisions() []*Collision {
	l.Lock()
	var result []*Collision
	for i := range l.collisions {
		result = append(result, l.collisions[i])
	}
	l.Unlock()
	return result
}

// AddAgent adds an Agent component to a GameObject
func (l *ObjectList) AddAgent(id Entity, agent *Agent) {
	l.Lock()
	agent.gameObject = l.entities[id]
	l.agents[id] = agent
	l.Unlock()
}

// Agent returns the Agent component for a GameObject
func (l *ObjectList) Agent(id Entity) *Agent {
	return l.agents[id]
}

// Agents returns all registered Agent components
func (l *ObjectList) Agents() []*Agent {
	l.Lock()
	var result []*Agent
	for i := range l.agents {
		result = append(result, l.agents[i])
	}
	l.Unlock()
	return result
}

// Deleted returns a list of GameObject IDs that has been deleted/removed
func (l *ObjectList) Deleted() []Entity {
	l.Lock()
	defer l.Unlock()
	return l.deleted
}

// ClearDeleted clears the list of deleted GameObjects
func (l *ObjectList) ClearDeleted() {
	l.Lock()
	defer l.Unlock()
	l.deleted = make([]Entity, 0)
}
