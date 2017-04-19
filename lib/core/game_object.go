package core

import (
	"github.com/stojg/vector"
)

// NewGameObject returns a new GameObject
func NewGameObject(name string, list *ObjectList) *GameObject {

	g := &GameObject{
		name: name,
		tags: make(map[string]bool),
		transform: &Transform{
			position:    vector.Zero(),
			orientation: vector.NewQuaternion(1, 0, 0, 0),
			scale:       vector.NewVector3(1, 1, 1),
		},
		list: list,
	}
	// link the transform back to the parent object
	g.transform.parent = g
	list.Add(g)
	return g
}

// GameObject is the base struct that all entities in the game should embed or use.
type GameObject struct {
	id        ID
	name      string
	transform *Transform
	tags      map[string]bool
	list      *ObjectList
}

// AddTags tags this object with tags
func (g *GameObject) AddTags(tags []string) {
	for i := range tags {
		g.tags[tags[i]] = true
	}
}

// ID returns the unique ID for this GameObject
func (g *GameObject) ID() ID {
	return g.id
}

// Transform returns the Transform for this GameObject
func (g *GameObject) Transform() *Transform {
	return g.transform
}

// AddGraphic adds a Graphic component to this GameObject
func (g *GameObject) AddGraphic(graphic *Graphic) {
	graphic.transform = g.transform
	g.list.AddGraphic(g.id, graphic)
}

// Graphic returns the Graphic component for this GameObject
func (g *GameObject) Graphic() *Graphic {
	return g.list.Graphic(g.id)
}

// AddBody adds a Body component to this GameObject
func (g *GameObject) AddBody(body *Body) {
	body.transform = g.transform
	g.list.AddBody(g.id, body)
}

// Body returns the Body component for this GameObject
func (g *GameObject) Body() *Body {
	return g.list.Body(g.id)
}

// AddCollision adds a Collision component to this GameObject
func (g *GameObject) AddCollision(collision *Collision) {
	collision.transform = g.transform
	g.list.AddCollision(g.id, collision)
}

// Collision returns the Collision component for this GameObject
func (g *GameObject) Collision() *Collision {
	return g.list.Collision(g.id)
}

// AddAgent adds an Agent component to this GameObject
func (g *GameObject) AddAgent(agent *Agent) {
	agent.transform = g.transform
	g.list.AddAgent(g.id, agent)
}

// Agent returns the Agent component for this GameObject
func (g *GameObject) Agent() *Agent {
	return g.list.Agent(g.id)
}

func (g *GameObject) AddInventory(inv *Inventory) {
	inv.transform = g.transform
	g.list.AddInventory(g.id, inv)
}

func (g *GameObject) Inventory() *Inventory {
	return g.list.Inventory(g.id)
}

// CompareTag returns true if this GameObject is tagged with a tag
func (g *GameObject) CompareTag(tag string) bool {
	_, ok := g.tags[tag]
	return ok
}
