package plan

import "testing"

func TestActionRemove(t *testing.T) {
	eat := newAction("eat")
	drop := newAction("drop")
	hide := newAction("hide")

	actions := []Action{eat, drop, hide}

	result := remove(actions, 1)

	if len(result) != 2 {
		t.Errorf("expected one action to be removed, got %d", len(result))
	}

	if result[0] != eat {
		t.Error("expected eat action to be in the list")
		return
	}

	if result[1] == drop {
		t.Error("didnt expected drop action to be in the list")
		return
	}

	if result[1] != hide {
		t.Error("did expected hide actiond to be in the list")
		return
	}
}

func TestActionNoRemoval(t *testing.T) {
	eat := newAction("eat")
	actions := []Action{eat}
	result := remove(actions, 1)
	if len(result) != 1 {
		t.Errorf("expected one action to be removed, got %d", len(result))
	}

	if result[0] != eat {
		t.Error("expected eat action to be in the list")
		return
	}
}

func BenchmarkActionRemove(b *testing.B) {
	for i := 0; i < b.N; i++ {
		eat := newAction("eat")
		drop := newAction("drop")
		hide := newAction("hide")
		actions := []Action{eat, drop, hide}
		remove(actions, 1)
	}
}

//
//func newGetFoodAction(cost float64) *getFoodAction {
//	a := &getFoodAction{
//		DefaultAction: NewAction("getFood", cost),
//		inRange:       false,
//	}
//	return a
//}
//
//type getFoodAction struct {
//	DefaultAction
//	hasFood bool
//	inRange bool
//}
//
//func (a *getFoodAction) CheckContextPrecondition(agent Agent) bool {
//	a.SetTarget([]int{10, 0, 200})
//	return true
//}
//
//func (a *getFoodAction) Perform(agent Agent) bool {
//	a.hasFood = true
//	return true
//}
//
//func (a *getFoodAction) IsDone() bool {
//	return a.hasFood
//}
//
//func (a *getFoodAction) InRange(agent Agent) bool {
//	return a.inRange
//}
//
//func newEatAction(cost float64) *eatingAction {
//	a := &eatingAction{
//		DefaultAction: NewAction("eat", cost),
//	}
//	return a
//}
//
//type eatingAction struct {
//	DefaultAction
//}
//
//func (a *eatingAction) Perform(agent Agent) bool {
//	return true
//}
//
//func (a *eatingAction) IsDone() bool {
//	return true
//}
//
//func (a *eatingAction) InRange(agent Agent) bool {
//	return true
//}
//
//func newSleepAction(cost float64) *sleepingAction {
//	return &sleepingAction{
//		DefaultAction: NewAction("sleep", cost),
//	}
//}
//
//type sleepingAction struct {
//	DefaultAction
//}
//
//func (a *sleepingAction) InRange(agent Agent) bool {
//	return true
//}
//
//func (a *sleepingAction) Perform(agent Agent) bool {
//	return true
//}
