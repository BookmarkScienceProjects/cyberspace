package formation

import (
	. "github.com/stojg/vector"
)

type Model struct {
	position    *Vector3
	orientation *Quaternion
}

func (m *Model) Position() *Vector3 {
	return m.position
}

func (m *Model) Orientation() *Quaternion {
	return m.orientation
}
