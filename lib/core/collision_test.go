package core

import (
	"testing"

	"github.com/stojg/vector"
)

func TestCollision_OBB(t *testing.T) {
	list := NewObjectList()
	obj := NewGameObject("test", list)
	rect := NewCollisionRectangle(1, 1, 1)
	obj.AddCollision(rect)
	obj.AddBody(NewBody(1, true))
	obj.Body().Transform().Orientation().AddScaledVector(vector.NewVector3(0.7, 0, 0.7), 1)
	obj.Body().CalculateDerivedData(obj.Body().Transform())

	obb := rect.OBB(1)

	expected := vector.NewVector3(0.85354, 0.71212, 0.85354)
	actual := obb.MaxPoint
	if !actual.Equals(expected) {
		t.Errorf("Expected max point to be %s, got %s", expected, actual)
	}

	expected = vector.NewVector3(-0.85354, -0.71212, -0.85354)
	actual = obb.MinPoint
	if !actual.Equals(expected) {
		t.Errorf("Expected min point to be %s, got %s", expected, actual)
	}
}

var result *OBB

func BenchmarkCollision_OBB(b *testing.B) {
	list := NewObjectList()
	obj := NewGameObject("test", list)
	rect := NewCollisionRectangle(1, 1, 1)
	obj.AddCollision(rect)
	obj.AddBody(NewBody(1, true))
	obj.Body().Transform().Orientation().AddScaledVector(vector.NewVector3(0.7, 0, 0.7), 1)
	obj.Body().CalculateDerivedData(obj.Body().Transform())

	for i := 0; i < b.N; i++ {
		result = rect.OBB(uint64(i))
	}

}
