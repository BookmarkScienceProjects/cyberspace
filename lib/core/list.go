package core

import (
	"math"
	"sync"

	"github.com/stojg/cyberspace/lib/quadtree"
)

// List is the primary resource for adding, removing and changing GameObjects and their components.
var List *ObjectList

func init() {
	List = NewObjectList()
}

func NewObjectList() *ObjectList {
	return &ObjectList{
		nextID:          1,
		entities:        make(map[ID]*GameObject),
		graphics:        make(map[ID]*Graphic),
		bodies:          make(map[ID]*Body),
		collisions:      make(map[ID]*Collision),
		collisionsFrame: 0,
		collisionCache:  make([]*Collision, 0),
		agents:          make(map[ID]*Agent),
		inventories:     make(map[ID]*Inventory),
		deleted:         make([]ID, 0),
		toDelete:        make([]ID, 0),
		senseManager:    &Manager{},
	}
}

// ObjectList is a struct that contains a list of GameObjects and their components. All creation,
// removal and changes should be handled by this list so they don't get lost or out of sync.
type ObjectList struct {
	sync.Mutex
	quadTree        *quadtree.QuadTree
	nextID          ID
	entities        map[ID]*GameObject
	graphics        map[ID]*Graphic
	bodies          map[ID]*Body
	collisions      map[ID]*Collision
	collisionsFrame uint64
	collisionCache  []*Collision
	agents          map[ID]*Agent
	inventories     map[ID]*Inventory
	senseManager    *Manager
	deleted         []ID
	toDelete        []ID
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
}

func (l *ObjectList) Get(id ID) *GameObject {
	l.Lock()
	defer l.Unlock()
	return l.entities[id]
}

func (l *ObjectList) SenseManager() *Manager {
	return l.senseManager
}

func (l *ObjectList) BuildQuadTree(frame uint64) {
	qTree := quadtree.NewQuadTree(
		quadtree.BoundingBox{
			MinX: -3200 / 2,
			MaxX: 3200 / 2,
			MinY: -3200 / 2,
			MaxY: 3200 / 2,
		},
	)
	// build the quad tree for broad phase collision detection
	for _, collision := range l.Collisions(frame) {
		qTree.Add(collision)
	}

	l.quadTree = &qTree
}

func (l *ObjectList) QuadTree() *quadtree.QuadTree {
	return l.quadTree
}

func (l *ObjectList) Flush() {
	l.Lock()
	for _, id := range l.toDelete {
		if _, found := l.graphics[id]; found {
			delete(l.graphics, id)
		}
		if _, found := l.bodies[id]; found {
			delete(l.bodies, id)
		}
		if _, found := l.agents[id]; found {
			delete(l.agents, id)
		}
		if _, found := l.collisions[id]; found {
			delete(l.collisions, id)
		}
		if _, found := l.entities[id]; found {
			delete(l.entities, id)
		}
		if _, found := l.inventories[id]; found {
			delete(l.inventories, id)
		}
		l.deleted = append(l.deleted, id)
		delete(l.entities, id)
	}
	l.toDelete = make([]ID, 0)
	l.Unlock()
}

// Remove a GameObject and all of it's components
func (l *ObjectList) Remove(g *GameObject) {
	l.Lock()
	l.toDelete = append(l.toDelete, g.id)
	l.Unlock()
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
func (l *ObjectList) AddGraphic(id ID, graphic *Graphic) {
	l.Lock()
	graphic.gameObject = l.entities[id]
	graphic.transform = l.entities[id].transform
	l.graphics[id] = graphic
	l.Unlock()
}

// Graphic returns the Graphic component for a GameObject
func (l *ObjectList) Graphic(id ID) *Graphic {
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
func (l *ObjectList) AddBody(id ID, body *Body) {
	l.Lock()
	body.gameObject = l.entities[id]
	body.transform = l.entities[id].transform
	l.bodies[id] = body
	l.Unlock()
}

// Body returns the Body component for a GameObject
func (l *ObjectList) Body(id ID) *Body {
	if body, found := l.bodies[id]; found {
		return body
	}
	panic("No body exists")
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
func (l *ObjectList) AddCollision(id ID, collision *Collision) {
	l.Lock()
	collision.gameObject = l.entities[id]
	collision.transform = l.entities[id].transform
	l.collisions[id] = collision
	l.Unlock()
}

// Collision returns the Collision component for a GameObject
func (l *ObjectList) Collision(id ID) *Collision {
	return l.collisions[id]
}

// Collisions returns all registered Collision components
func (l *ObjectList) Collisions(frame uint64) []*Collision {
	if frame == l.collisionsFrame {
		return l.collisionCache
	}
	l.collisionCache = make([]*Collision, len(l.collisions))
	count := 0
	for i := range l.collisions {
		l.collisionCache[count] = l.collisions[i]
		count++
	}
	l.collisionsFrame = frame
	return l.collisionCache
}

// AddAgent adds an Agent component to a GameObject
func (l *ObjectList) AddAgent(id ID, agent *Agent) {
	l.Lock()
	agent.gameObject = l.entities[id]
	l.agents[id] = agent
	l.Unlock()
}

// Agent returns the Agent component for a GameObject
func (l *ObjectList) Agent(id ID) *Agent {
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
func (l *ObjectList) Deleted() []ID {
	l.Lock()
	defer l.Unlock()
	return l.deleted
}

// ClearDeleted clears the list of deleted GameObjects
func (l *ObjectList) ClearDeleted() {
	l.Lock()
	defer l.Unlock()
	l.deleted = make([]ID, 0)
}
