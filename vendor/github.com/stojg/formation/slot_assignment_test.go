package formation

import (
	"github.com/stojg/vector"
	"testing"
)

type TestCharacter struct {
	position    *vector.Vector3
	orientation *vector.Quaternion
	target      Static
}

func (c *TestCharacter) Position() *vector.Vector3 {
	return c.position
}

func (c *TestCharacter) Orientation() *vector.Quaternion {
	return c.orientation
}

func (c *TestCharacter) SetTarget(t Static) {
	c.target = t
}

func TestSlotAssignments_Add(t *testing.T) {
	assignments := make(SlotAssignments, 0)

	firstChar := &TestCharacter{
		position: vector.NewVector3(10, 10, 10),
	}
	firstAssignment := &SlotAssignment{
		character: firstChar,
	}

	secondChar := &TestCharacter{
		position: vector.NewVector3(20, 20, 20),
	}
	secondAssignment := &SlotAssignment{
		character: secondChar,
	}

	if len(assignments) != 0 {
		t.Errorf("assignments should have a length of 0")
	}

	assignments = append(assignments, firstAssignment)

	if len(assignments) != 1 {
		t.Errorf("assignments should have a length of 1")
	}

	assignments = append(assignments, secondAssignment)
	if len(assignments) != 2 {
		t.Errorf("assignments should have a length of 2")
	}

	i, ok := assignments.find(secondChar)
	if !ok {
		t.Errorf("Should have found char")
	}

	c := assignments[i]
	if !c.character.Position().Equals(secondChar.position) {
		t.Errorf("We didnt get the same character back? %s", c.character)
	}

	assignments.remove(i)

	if len(assignments) != 1 {
		t.Errorf("assignments should have a length of 1")
	}
}
