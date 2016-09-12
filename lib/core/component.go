package core

// Base class for everything attached to GameObjects.
//
// Note that your code will never directly create a Component. Instead, you write script code, and
// attach the script to a GameObject. See Also: ScriptableObject as a way to create scripts that do
// not attach to any GameObject.
type Component struct {
	// pointer to the parent class
	gameObject *GameObject
	//
	tag string

	transform *Transform
}

func (c *Component) GameObject() *GameObject {
	return c.gameObject
}

func (c *Component) Transform() *Transform {
	return c.transform
}
