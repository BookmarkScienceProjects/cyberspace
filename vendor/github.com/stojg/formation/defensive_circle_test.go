package formation_test

import (
	"github.com/stojg/formation"
	"github.com/stojg/vector"
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

func TestDefensiveCirclePattern_DriftOffset(t *testing.T) {
	def := formation.NewDefensiveCircle(10, 2)
	man := formation.NewManager(def)
	firstChar := &TestCharacter{
		position: vector.NewVector3(0, 0, 0),
	}
	man.AddCharacter(firstChar)
	secondChar := &TestCharacter{
		position: vector.NewVector3(0, 0, 0),
	}
	man.AddCharacter(secondChar)

	man.UpdateSlots()

	expects := vector.NewVector3(10, 0, 0)
	if !firstChar.target.Position().Equals(expects) {
		t.Errorf("Pos should be %v, got %v", expects, firstChar.target.Position())
	}
	expects = vector.NewVector3(-10, 0, 0)
	if !secondChar.target.Position().Equals(expects) {
		t.Errorf("Pos should be %v, got %v", expects, secondChar.target.Position())
	}

}

func TestDefensiveCirclePattern_SlotLocation(t *testing.T) {
	var slotLocationTests = []struct {
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

	for _, tt := range slotLocationTests {
		pattern := formation.NewDefensiveCircle(10, tt.initSlotNumber)
		loc := pattern.SlotLocation(tt.slotNumber)
		if !loc.Position().Equals(tt.expected) {
			t.Errorf("Pos should be %v, got %v", tt.expected, loc.Position())
		}
	}
}
