package plan

import (
	"testing"

	"github.com/stojg/vector"
)

func TestStateList(t *testing.T) {

	const (
		isHappy = "isHappy"
		inRange = "inRange"
	)

	states := make(StateList)
	states.Set(isHappy, true)
	states.Set(inRange, vector.Zero())

	if !states.Have(isHappy) {
		t.Error("Expected IsHappy to be set")
	}
	if !states.Get(isHappy).(bool) {
		t.Error("Expected IsHappy to be true")
	}

	states.Set(isHappy, false)
	if states.Get(isHappy).(bool) {
		t.Error("Expected IsHappy to be false")
	}
	states.Unset(isHappy)
	if states.Have(isHappy) {
		t.Error("Expected IsHappy to not be set")
	}
	if val := states.Get(isHappy); val != nil {
		t.Errorf("Expected IsHappy to be nil, got %v", val)
	}

	if !states.Have(inRange) {
		t.Error("Expected that Have would work on inRange")
	}
	if val := states.Get(inRange).(*vector.Vector3); !val.Equals(vector.Zero()) {
		t.Errorf("Expected InRange = %s, got %v", vector.Zero(), val)
	}

	states.Set(inRange, vector.X())
	if val := states.Get(inRange).(*vector.Vector3); !val.Equals(vector.X()) {
		t.Errorf("Expected InRange = %s, got %v", vector.X(), val)
	}
	states.Unset(inRange)
	if val := states.Get(inRange); val != nil {
		t.Errorf("Expected InRange to be nil, got %+v", val)
	}

	// test that removing a non existant value is a NOP
	states.Unset(State("doesnt_exists"))
}

func TestStateList_Includes_true(t *testing.T) {

	state := make(StateList)
	state["food"] = true
	state["temperature"] = false

	test := make(StateList)
	test["food"] = true

	actual := state.Includes(test)
	if !actual {
		t.Error("expected that inState would be true")
	}
}

func TestStateList_Includes_EmptyIsTrue(t *testing.T) {

	state := make(StateList)
	state["food"] = true
	state["temperature"] = false

	test := make(StateList)

	actual := state.Includes(test)
	if !actual {
		t.Error("expected that inState would be true")
	}
}

func TestStateList_Includes_dont_exists(t *testing.T) {

	state := make(StateList)
	state["hasFood"] = true
	state["isFull"] = false

	test := make(StateList)
	test["isHurt"] = true

	actual := state.Includes(test)
	if actual {
		t.Error("expected that inState would be false")
	}
}

func TestStateList_Compare_true(t *testing.T) {
	state := make(StateList)
	state["hasFood"] = true
	state["isFull"] = false

	test := make(StateList)
	test["hasFood"] = true

	actual := state.Compare(test)
	if !actual {
		t.Error("expected that Compare would be true")
	}
}

func TestStateList_Compare_false_false(t *testing.T) {
	state := make(StateList)
	state["hasFood"] = false
	state["isFull"] = false

	test := make(StateList)
	test["hasFood"] = false
	test["isFull"] = false

	actual := state.Compare(test)
	if !actual {
		t.Error("expected that Compare would be true")
	}
}

func TestStateList_Compare_false_value(t *testing.T) {
	state := make(StateList)
	state["hasFood"] = true
	state["isFull"] = false

	test := make(StateList)
	test["hasFood"] = false

	actual := state.Compare(test)
	if actual {
		t.Error("expected that Compare would be false")
	}
}

func TestStateList_Compare_DontExist(t *testing.T) {
	state := make(StateList)
	state["hasFood"] = true
	state["isFull"] = false

	test := make(StateList)
	test["isHurt"] = true

	actual := state.Compare(test)
	if actual {
		t.Error("expected that comapre would be false")
	}

}

func BenchmarkStateList_Includes(b *testing.B) {
	state := make(StateList)
	state["hasFood"] = true
	state["isFull"] = false

	test := make(StateList)
	test["isHurt"] = true

	for i := 0; i < b.N; i++ {
		state.Includes(test)
	}
}

func TestStateList_Clone_Update(t *testing.T) {

	currentState := make(StateList)
	currentState["food"] = false
	currentState["temperature"] = 23.4

	changes := make(StateList)
	changes["food"] = true

	result := currentState.Clone().Apply(changes)

	if len(result) != len(currentState) {
		t.Error("result should have the same # of entries as current state")
	}

	if _, ok := result["food"]; !ok {
		t.Logf("%s", result.String())
		t.Error("could not find 'food' state")
		return
	}

	if !result["food"].(bool) {
		t.Errorf("food state was not changed, expected true, got %t", result["food"])
	}

	if currentState["food"].(bool) {
		t.Error("currentState failed to be treated as an immutable")
	}

	if val := result["temperature"].(float64); val != 23.4 {
		t.Error("unrelated state was changed, temperature should be 23.4, got ", result["temperature"])
	}
}

func BenchmarkStateList_Clone_Apply(b *testing.B) {
	currentState := make(StateList)
	currentState["food"] = false
	currentState["temperature"] = 23.4
	changes := make(StateList)
	changes["food"] = true
	for i := 0; i < b.N; i++ {
		currentState.Clone().Apply(changes)
	}
}

func BenchmarkStateList_Unset(b *testing.B) {
	currentState := make(StateList)
	currentState["food"] = false
	for i := 0; i < b.N; i++ {
		currentState.Unset("banana")
	}
}
