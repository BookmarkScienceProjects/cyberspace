package plan

import (
	"testing"

	"github.com/stojg/vector"
)

func TestFind_NoGoals(t *testing.T) {

	defer func() {
		if r := recover(); r == nil {
			t.Error("The code did not panic on no goal state")
		}
	}()
	worldState := make(StateList)

	goalState := make(StateList)

	a := &agent{}

	Find(a, worldState, goalState)
}

func TestFind_NoActions(t *testing.T) {
	const (
		Hungry = "isHungry"
	)
	worldState := make(StateList)

	goalState := make(StateList)
	goalState.Set(Hungry, true)

	a := &agent{}

	actions := Find(a, worldState, goalState)
	if len(actions) != 0 {
		t.Error("plan should have failed with no actions")
	}
}

func TestFind_Single_Action_Failed(t *testing.T) {
	const (
		killEnemy = "action"
		seeEnemy  = "seeEnemy"
		Hungry    = "isHungry"
	)

	action := newAction("killEnemy")
	action.preConditions.Set(seeEnemy, true)
	action.effects.Set(killEnemy, true)

	worldState := make(StateList)

	goalState := make(StateList)
	goalState.Set(Hungry, true)

	a := &agent{}
	a.actions = append(a.actions, action)

	actions := Find(a, worldState, goalState)
	if len(actions) != 0 {
		t.Error("plan should have failed, got:")
		for _, a := range actions {
			t.Errorf(" - %T\n", a)
		}
	}
}

func TestFind_Single_Action_Success(t *testing.T) {
	const (
		killEnemy = "action"
		seeEnemy  = "seeEnemy"
	)

	action := newAction("killEnemy")
	action.preConditions.Set(seeEnemy, true)
	action.effects.Set(killEnemy, true)

	worldState := make(StateList)
	worldState.Set(seeEnemy, true)

	goalState := make(StateList)
	goalState.Set(killEnemy, true)

	a := &agent{}
	a.actions = append(a.actions, action)

	actions := Find(a, worldState, goalState)
	if len(actions) == 0 {
		t.Error("plan should not have failed")
	}
}

func TestFind_ThreeActionPlan(t *testing.T) {
	const (
		armed     = "armed"
		killEnemy = "killEnemy"
		seeEnemy  = "seeEnemy"
	)

	killAction := newAction("killEnemy")
	killAction.preConditions.Set(seeEnemy, true)
	killAction.preConditions.Set(armed, true)
	killAction.effects.Set(killEnemy, true)

	airStrikeAction := newAction("airStrike")
	airStrikeAction.effects.Set(killEnemy, true)
	airStrikeAction.cost = 10

	drawWeapenAction := newAction("drawWeapon")
	drawWeapenAction.effects.Set(armed, true)

	findAction := newAction("findEnemy")
	findAction.effects.Set(seeEnemy, true)

	worldState := make(StateList)

	goalState := make(StateList)
	goalState.Set(killEnemy, true)

	a := &agent{}
	a.actions = append(a.actions, killAction)
	a.actions = append(a.actions, airStrikeAction)
	a.actions = append(a.actions, findAction)
	a.actions = append(a.actions, drawWeapenAction)

	actions := Find(a, worldState, goalState)
	if len(actions) != 3 {
		t.Errorf("plan should have 2 actions, got %d action", len(actions))
		return
	}

	if !(actions[0].(*action).name == findAction.name || actions[0].(*action).name == drawWeapenAction.name) {
		t.Errorf("first action should have been either '%s' or '%s', got '%s'", drawWeapenAction.name, findAction.name, actions[0].(*action).name)
		return
	}

	if actions[0].(*action).name == actions[1].(*action).name {
		t.Errorf("First and second action should not be the same, got '%s' and '%s'", actions[0].(*action).name, actions[1].(*action).name)
		return
	}

	if !(actions[1].(*action).name == findAction.name || actions[1].(*action).name == drawWeapenAction.name) {
		t.Errorf("second action should have been either '%s' or '%s', got '%s'", drawWeapenAction.name, findAction.name, actions[0].(*action).name)
		return
	}

	if actions[2].(*action).name != killAction.name {
		t.Errorf("second action should have been '%s', got '%s'", killAction.name, actions[1].(*action).name)
		return
	}
}

func TestFind_NeedsToBeInRange(t *testing.T) {
	const (
		killEnemy = "killEnemy"
		target    = "target"
	)

	kill := newkillInRangeAction()
	kill.effects.Set(killEnemy, true)

	airStrike := newAction("airStrike")
	airStrike.effects.Set(killEnemy, true)
	airStrike.cost = 10

	move := newMoveToAction()
	move.cost = 2

	worldState := make(StateList)
	worldState.Set(target, 2)

	goalState := make(StateList)
	goalState.Set(killEnemy, true)

	a := &agent{}
	a.actions = append(a.actions, kill)
	a.actions = append(a.actions, airStrike)
	a.actions = append(a.actions, move)

	actions := Find(a, worldState, goalState)
	if len(actions) != 2 {
		t.Errorf("plan should have 2 actions, got %d action", len(actions))
		t.Errorf("first action is %v", actions[0])
		return
	}
	if a := actions[0].(*moveToAction); a.name != move.name {
		t.Errorf("first action should be '%s', got '%s' action", move.name, a.name)
	}
}

func TestFind_CheckContextPrecondition_Failed(t *testing.T) {
	const (
		killEnemy = "killEnemy"
		target    = "target"
	)

	kill := newkillInRangeAction()
	kill.effects.Set(killEnemy, true)
	kill.contextPrecondition = false

	airStrike := newAction("airStrike")
	airStrike.effects.Set(killEnemy, true)
	airStrike.cost = 10

	move := newMoveToAction()
	move.cost = 2

	worldState := make(StateList)
	worldState.Set(target, 2)

	goalState := make(StateList)
	goalState.Set(killEnemy, true)

	a := &agent{}
	a.actions = append(a.actions, kill)
	a.actions = append(a.actions, airStrike)
	a.actions = append(a.actions, move)

	actions := Find(a, worldState, goalState)
	if len(actions) != 1 {
		t.Errorf("plan should have 1 actions, got %d action", len(actions))
		t.Errorf("first action is %v", actions[0])
		return
	}
	if a := actions[0].(*action); a.name != airStrike.name {
		t.Errorf("first action should be '%s', got '%s' action", airStrike.name, a.name)
	}
}

func BenchmarkFind(b *testing.B) {
	const (
		armed     = "armed"
		killEnemy = "killAction"
		seeEnemy  = "seeEnemy"
	)

	killAction := newAction("killEnemy")
	killAction.cost = 2
	killAction.preConditions.Set(seeEnemy, true)
	killAction.preConditions.Set(armed, true)
	killAction.effects.Set(killEnemy, true)

	airStrikeAction := newAction("airStrike")
	airStrikeAction.effects.Set(killEnemy, true)
	airStrikeAction.cost = 10

	drawWeapenAction := newAction("drawWeapon")
	drawWeapenAction.cost = 1
	drawWeapenAction.effects.Set(armed, true)

	findAction := newAction("findEnemy")
	findAction.cost = 3
	findAction.effects.Set(seeEnemy, true)

	worldState := make(StateList)

	goalState := make(StateList)
	goalState.Set(killEnemy, true)

	a := &agent{}
	a.actions = append(a.actions, killAction)
	a.actions = append(a.actions, airStrikeAction)
	a.actions = append(a.actions, findAction)
	a.actions = append(a.actions, drawWeapenAction)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Find(a, worldState, goalState)
	}
}

func newkillInRangeAction() *killInRangeAction {
	return &killInRangeAction{
		action: newAction("killEnemyInRange"),
	}
}

func newMoveToAction() *moveToAction {
	return &moveToAction{
		action: newAction("moveTo"),
	}
}

func newAction(name string) *action {
	return &action{
		name:                name,
		effects:             make(StateList),
		preConditions:       make(StateList),
		contextPrecondition: true,
	}
}

type killInRangeAction struct {
	*action
}

func (a *killInRangeAction) Preconditions() StateList {
	s := make(StateList)
	// let's say we are not there
	s.Set("inRange", vector.Zero())
	return s
}

type moveToAction struct {
	*action
	target *vector.Vector3
}

func (a *moveToAction) Preconditions() StateList {
	return a.preConditions
}

func (a *moveToAction) Effects() StateList {
	s := make(StateList)
	// let's say we are not there
	s.Set("inRange", vector.Zero())
	return s
}

type action struct {
	name                string
	effects             StateList
	preConditions       StateList
	cost                float64
	contextPrecondition bool
}

func (a *action) Cost() float64 { return a.cost }

func (a *action) Effects() StateList { return a.effects }

func (a *action) Preconditions() StateList { return a.preConditions }

func (a *action) CheckContextPrecondition(state StateList) bool { return a.contextPrecondition }

func (a *action) Reset() {}

func (a *action) MoveTo() interface{} { return true }

func (a *action) Target() interface{} { return nil }

func (a *action) Run(agent Agent) (bool, error) {
	return true, nil
}

type agent struct {
	actions []Action
	plan    []Action
	goals   []StateList
	state   StateList
}

func (a *agent) AvailableActions() []Action { return a.actions }

func (a *agent) Plan() []Action { return a.plan }

func (a *agent) Update(elapsed float64) {}

func (a *agent) PopCurrentAction() {}

func (a *agent) SetPlan(n []Action) {}

func (a *agent) State() StateList {
	return a.state
}

func (a *agent) Goals() []StateList {
	return a.goals
}

func (a *agent) Move(action Action) (bool, error) {
	return true, nil
}
