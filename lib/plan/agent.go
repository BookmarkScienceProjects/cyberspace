package plan

// Debuuger provides an interface for the planner to give feedback to the Agent and report
// success/failure.
type Debugger interface {

	// No sequence of actions could be found for the supplied goal. You will need to try another goal
	PlanFailed(failedGoals []StateList)

	// A plan was found for the supplied goal. These are the actions the Agent will perform, in order.
	PlanFound(goal StateList, actions []Action)

	// All actions are complete and the goal was reached.
	ActionsFinished()

	// One of the actions caused the plan to abort. That action is passed in.
	PlanAborted(Action, error)
}

// Agent must be implemented by agents that wants to use goal oriented action planning
type Agent interface {
	//Debugger

	// Get the actions that this agent can do
	AvailableActions() []Action

	//	Set the actions that will allow this agent to reach it's goal
	SetPlan([]Action)

	// Get current planned actions
	Plan() []Action

	// Move the agent towards the target in order for the next action to be able to perform. Return true if the Agent
	// is at the target and the next action can perform, false if it is not there yet.
	Move(Action) (bool, error)

	// Advance the internal state machine and run actions
	Update(elapsed float64)

	//Remove the currently running DefaultAction
	PopCurrentAction()

	//The starting state of the Agent and the world. Supplies what states are needed for actions to run.
	State() StateList

	// Get the goals for this actor
	Goals() []StateList
}
