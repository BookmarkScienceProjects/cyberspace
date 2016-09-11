package GameObject

// Base class for everything attached to GameObjects.
//
// Note that your code will never directly create a Component. Instead, you write script code, and
// attach the script to a GameObject. See Also: ScriptableObject as a way to create scripts that do
// not attach to any GameObject.
type Component struct {
	gameobject *GameObject
	tag string
	transform *Transform
}

func (c *Component) Attach(g *GameObject) {
	c.gameobject = g
}
