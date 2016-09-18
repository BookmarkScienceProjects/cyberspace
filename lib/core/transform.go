package core

import "github.com/stojg/vector"

// Transform - Position, rotation and scale of an object.
// Every object in a scene has a Transform. It's used to store and manipulate the position, rotation
// and scale of the object.
type Transform struct {
	position    *vector.Vector3
	orientation *vector.Quaternion
	scale       *vector.Vector3
	// parent object
	parent *GameObject
}

// Position returns the current position of the transform
func (t *Transform) Position() *vector.Vector3 {
	return t.position
}

// Orientation returns the current orientation (heading) of the transform
func (t *Transform) Orientation() *vector.Quaternion {
	return t.orientation
}

// Scale returns the current scale (width, height, depth) of the transform
func (t *Transform) Scale() *vector.Vector3 {
	return t.scale
}
