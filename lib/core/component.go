package core

type Entity uint32

// Component is the base struct that should be embedded on all components
type Component struct {
	// pointer to the parent class
	gameObject *GameObject
	// pointer to the transform
	transform *Transform
}

// GameObject returns the GameObject this component is attached to
func (c *Component) GameObject() *GameObject {
	return c.gameObject
}

// Transform is a short code to get the Transform for the parent GameObject
func (c *Component) Transform() *Transform {
	return c.transform
}
