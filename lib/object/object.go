package object

import (
	"github.com/stojg/vector"
	"github.com/stojg/vivere/lib/components"
)

type Object interface {
	ID() *components.Entity
	Kind() Kind
	Position() *vector.Vector3
	Orientation() *vector.Quaternion
	Awake() bool
	Size() *vector.Vector3
	Rendered() bool
	SetRendered()
}
