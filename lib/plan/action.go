package plan

// remove removes one action from the list, it assumes that actions are unique in the list
func remove(actions []Action, i int) []Action {
	if i > len(actions)-1 {
		return actions
	}
	return append(actions[:i], actions[i+1:]...)
}

// Action is the interface that describes what the planner and
type Action interface {
	// The cost of performing the action.
	Cost() float64

	// Preconditions returns the StateList that this action requires
	Preconditions() StateList

	// Effects returns how this action effects the world state
	Effects() StateList

	// Reset any variables that need to be reset before plan happens again.
	Reset()

	// MoveTo tells the state machine that it needs to move somewhere
	MoveTo() interface{}

	//Procedurally check if this action can run. Not all actions will need
	//this, but some might.
	CheckContextPrecondition(state StateList) bool

	// Run the action. First return values is the success, second returns an error if it can no longer be exectuted
	Run(Agent) (bool, error)
}
