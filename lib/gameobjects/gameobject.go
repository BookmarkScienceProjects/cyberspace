package GameObject

func New(name string) *GameObject {
	return &GameObject{
		name: name,
	}
}

type GameObject struct {
	name string
	components []Component
}

// Adds a component class named className to the game object.
func (g *GameObject) AddComponent(c Component) {
	g.components = append(g.components, c)
	c.Attach(g)
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



