package GameObject

import "github.com/stojg/vector"

// Transform - Position, rotation and scale of an object.
// Every object in a scene has a Transform. It's used to store and manipulate the position, rotation
// and scale of the object.
type Transform struct {
	Component
	position *vector.Vector3
	rotation *vector.Quaternion
	scale *vector.Vector3
	gameobject *GameObject
}

func (t *Transform) Position() *vector.Vector3 {
	return t.position
}

func (t *Transform) Rotation() *vector.Quaternion {
	return t.rotation
}

func (t *Transform) Scale() *vector.Vector3 {
	return t.scale
}

func (t *Transform) GameObject() *GameObject {
	return t.gameobject
}