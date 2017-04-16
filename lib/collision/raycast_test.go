package collision

import (
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/vector"
	"testing"
)

func TestRaycastFalse(t *testing.T) {
	list := core.NewObjectList()

	me := core.NewGameObject("me", list)
	me.Transform().Position().Set(0, 0, 0)
	me.AddBody(core.NewBody(1))
	me.Body().CalculateDerivedData(me.Transform())
	me.AddCollision(core.NewCollisionRectangle(1, 1, 1))

	blocker := core.NewGameObject("blocker", list)
	blocker.Transform().Position().Set(5, 0, 0)
	blocker.AddBody(core.NewBody(1))
	blocker.Body().CalculateDerivedData(blocker.Transform())
	blocker.AddCollision(core.NewCollisionRectangle(1, 1, 1))

	friend := core.NewGameObject("friend", list)
	friend.Transform().Position().Set(10, 0, 0)
	friend.AddBody(core.NewBody(1))
	friend.Body().CalculateDerivedData(friend.Transform())
	friend.AddCollision(core.NewCollisionRectangle(1, 1, 1))

	direction := vector.NewVector3(5, 0, 0)

	result := Raycast(me.Transform().Position(), direction, list)

	if len(result) != 2 {
		t.Errorf("Expected 2 results in raycast, got %d", len(result))
		return
	}

	if result[0].Collision != me.Collision() {
		t.Errorf("%+v\n", result[0])
	}

	if result[1].Collision != blocker.Collision() {
		t.Errorf("%+v\n", result[0])
	}
}
