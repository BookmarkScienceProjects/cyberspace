package formation_test

import (
	"github.com/stojg/cyberspace/lib/formation"
	"github.com/stojg/vivere/lib/vector"
	"testing"
)

type TestCharacter struct {
	position    *vector.Vector3
	orientation *vector.Quaternion
	target      formation.Static
}

func (c *TestCharacter) Position() *vector.Vector3 {
	return c.position
}

func (c *TestCharacter) Orientation() *vector.Quaternion {
	return c.orientation
}

func (c *TestCharacter) SetTarget(t formation.Static) {
	c.target = t
}

var defTests = []struct {
	initSlotNumber int
	slotNumber     int
	expected       *vector.Vector3
}{
	{1, 1, vector.NewVector3(0, 0, 0)},
	{2, 1, vector.NewVector3(-10, 0, 0)},
	{2, 2, vector.NewVector3(10, 0, 0)},
	{3, 1, vector.NewVector3(-5.77350, 0, 10)},
	{3, 2, vector.NewVector3(-5.77350, 0, -10)},
	{3, 3, vector.NewVector3(11.54701, 0, 0)},
	{4, 1, vector.NewVector3(0, 0, 14.14214)},
	{4, 2, vector.NewVector3(-14.14214, 0, 0)},
	{4, 3, vector.NewVector3(0, 0, -14.14214)},
	{4, 4, vector.NewVector3(14.14214, 0, 0)},
}

func TestDefensiveCirclePattern_SlotLocation(t *testing.T) {
	for _, tt := range defTests {
		pattern := formation.NewDefensiveCircle(10, tt.initSlotNumber)
		loc := pattern.SlotLocation(tt.slotNumber)
		if !loc.Position().Equals(tt.expected) {
			t.Errorf("Pos should be %s, got %s", tt.expected, loc.Position())
		}
	}
}

//func TestSomething(t *testing.T) {
//
//	defPattern := formation.NewDefensiveCircle(10)
//
//	loc := defPattern.SlotLocation(1)
//	pos := loc.Position()
//	t.Errorf("%s", pos)
//
//	manager := formation.NewManager(defPattern)
//
//	char := &TestCharacter{
//		position:    vector.NewVector3(0, 0, 0),
//		orientation: vector.NewQuaternion(0, 0, 0, 1),
//	}
//	added := manager.AddCharacter(char)
//
//	if !added {
//		t.Errorf("it should always to be possible to add a char to the defensive circle")
//	}
//
//	manager.UpdateSlots()
//
//	t.Errorf("%s", char.target.Position())
//
//	//manager.UpdateSlots()
//}
