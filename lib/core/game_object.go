package core

import (
	"github.com/stojg/vector"
	"github.com/stojg/vivere/lib/components"
)

func NewGameObject(name string) *GameObject {

	g := &GameObject{
		name: name,
		transform: &Transform{
			position: vector.Zero(),
			orientation: vector.NewQuaternion(1, 0, 0, 0),
			scale:    vector.NewVector3(1, 1, 1),
		},
	}
	g.transform.parent = g
	List.Add(g)
	// link the transform back to the parent object
	return g
}

type GameObject struct {
	id         components.Entity
	name       string
	transform  *Transform
}

func (g *GameObject) ID() components.Entity {
	return g.id
}

func (g *GameObject) Transform() *Transform {
	return g.transform
}

func (g *GameObject) AddGraphic(graphic *Graphic) {
	graphic.transform = g.transform
	List.AddGraphic(g.id, graphic)
}

func (g *GameObject) Graphic() *Graphic {
	return List.Graphic(g.id)
}

func (g *GameObject) AddBody(body *Body) {
	body.transform = g.transform
	List.AddBody(g.id, body)
}

func (g *GameObject) Body() *Body {
	return List.Body(g.id)
}

func (g *GameObject) AddCollision(collision *Collision) {
	collision.transform = g.transform
	List.AddCollision(g.id, collision)
}

func (g *GameObject) Collision() *Collision {
	return List.Collision(g.id)
}

// Calls the method named methodName on every MonoBehaviour in this game object or any of its children.
func (g *GameObject) BroadcastMessage() {

}

// Is this game object tagged with tag ?
func (g *GameObject) CompareTag() {

}

// Returns the component of Type type if the game object has one attached, null if it doesn't.
func (g *GameObject) GetComponent() {

}

// Returns the component of Type type in the GameObject or any of its children using depth first search.
func (g *GameObject) GetComponentInChildren() {

}

// Returns the component of Type type in the GameObject or any of its parents.
func (g *GameObject) GetComponentInParent() {

}

// Returns all components of Type type in the GameObject.
func (g *GameObject) GetComponents() {

}

// Returns all components of Type type in the GameObject or any of its children.
func (g *GameObject) GetComponentsInChildren() {

}

// Returns all components of Type type in the GameObject or any of its parents.
func (g *GameObject) GetComponentsInParent() {

}

// Calls the method named methodName on every MonoBehaviour in this game object.
func (g *GameObject) SendMessage() {

}

// Calls the method named methodName on every MonoBehaviour in this game object and on every ancestor of the behaviour.
func (g *GameObject) SendMessageUpwards() {

}

// Activates/Deactivates the GameObject.
func (g *GameObject) SetActive() {

}
