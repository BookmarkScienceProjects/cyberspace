package core

import (
	"fmt"
	"github.com/stojg/goap"
	"github.com/stojg/steering"
)

// NewAgent returns an initialised agent ready for action!
func NewAgent() *Agent {
	a := &Agent{
		fsm:        &goap.FSM{},
		worldState: make(goap.StateList, 0),
		innerState: make(goap.StateList, 0),
		goals:      make(goap.StateList, 0),
	}
	a.fsm.PushState(IdleState)
	return a
}

// Agent is the core struct that represents an AI entity that plan and execute actions.
type Agent struct {
	Component

	// fsm is the internal state machine that keeps tracks in which state the agent currently 'are'
	// in.
	fsm *goap.FSM

	// This is provided be a higher level AI component and represents what the agent knows about
	// the world in general
	worldState goap.StateList

	// This represents what the agent thinks about itself, is it hungry, tired etc
	innerState goap.StateList

	// what goals are this agent trying to fulfill
	goals goap.StateList

	// A list of actions that the agent can do
	availableActions []goap.Actionable

	// what actions does the agent currently scheduled
	currentActions []goap.Actionable
}

// StateMachine returns the FSM for this agent
func (a *Agent) StateMachine() *goap.FSM {
	return a.fsm
}

// AvailableActions returns a list of actions that this agent can use to fulfill it's goal
func (a *Agent) AvailableActions() []goap.Actionable {
	return a.availableActions
}

// SetAvailableActions sets the available actions for this agent
func (a *Agent) SetAvailableActions(actions []goap.Actionable) {
	a.availableActions = actions
}

// CurrentAction returns the currently running action
func (a *Agent) CurrentAction() goap.Actionable {
	return a.currentActions[0]
}

// PopCurrentAction remove the currently running action
func (a *Agent) PopCurrentAction() {
	a.currentActions = a.currentActions[1:]
}

// SetCurrentActions sets the list of current actions that the agent should run, first action first in list
func (a *Agent) SetCurrentActions(actions []goap.Actionable) {
	a.currentActions = actions
}

// HasActionPlan will return true if this agent have any actions scheduled or running
func (a *Agent) HasActionPlan() bool {
	return len(a.currentActions) > 0
}

// GetWorldState returns the combined state of the world and the agents inner states. is it hungry,
// is it tired, can it see food? etc.
func (a *Agent) GetWorldState() goap.StateList {
	res := make(goap.StateList, 0)
	for k, v := range a.worldState {
		res[k] = v
	}
	for k, v := range a.innerState {
		res[k] = v
	}
	return res
}

// SetWorldState will set the state that the world provides the agent with
func (a *Agent) SetWorldState(state goap.StateList) {
	a.worldState = state
}

// SetState sets the inner state for this agent
func (a *Agent) SetState(k string, v interface{}) {
	a.innerState[k] = v
}

// CreateGoalState gives the goap planner a new goal so it can figure out the actions needed to
// fulfill it.
func (a *Agent) CreateGoalState() goap.StateList {
	return a.goals
}

// AddGoal sets what goals this agent is trying to fulfill
func (a *Agent) AddGoal(name string, value interface{}) {
	a.goals[name] = value
}

// PlanFailed is called when there is no sequence of actions could be found for the supplied goal.
// You will need to try another goal
func (a *Agent) PlanFailed(failedGoal goap.StateList) {}

// PlanFound is called when a plan was found for the supplied goal. The actions contains the plan
// of actions the agent will perform, in order.
func (a *Agent) PlanFound(goal goap.StateList, actions []goap.Actionable) {}

// ActionsFinished is signaled when all actions are complete and the goal was reached.
func (a *Agent) ActionsFinished() {}

// PlanAborted is called when one of the actions in the plan have discovered that it can no longer
// be done.
func (a *Agent) PlanAborted(abortingAction goap.Actionable) {
	fmt.Printf("plan was aborted by action %s aborted", abortingAction)
}

// MoveAgent is when the agent must move towards the target in order for the next action to be able
// to perform. Return true if the Agent is at the target and the next action can perform. False if
// it is not there yet.
func (a *Agent) MoveAgent(nextAction goap.Actionable) bool {
	target, found := nextAction.Target().(*GameObject)
	if !found {
		fmt.Printf("in core.Agent.MoveAgent: %s is not a *GameObject", nextAction.Target())
		return false
	}

	dist := a.gameObject.transform.position.NewSub(target.transform.position).Length()
	if dist < a.gameObject.transform.scale[0]*1.5 {
		nextAction.SetInRange()
		return true
	}

	arrive := steering.NewArrive(a.gameObject.Body(), target.transform.Position(), 3).Get()
	look := steering.NewLookWhereYoureGoing(a.gameObject.Body()).Get()

	a.gameObject.Body().AddForce(arrive.Linear())
	a.gameObject.Body().AddTorque(look.Angular())

	return false
}

// IdleState is the FSM state that the agent is in when it's planning
func IdleState(fsm *goap.FSM, agent goap.Agent) {
	worldState := agent.GetWorldState()
	goal := agent.CreateGoalState()
	plan := goap.Plan(agent, agent.AvailableActions(), worldState, goal)
	if plan != nil {
		agent.SetCurrentActions(plan)
		agent.PlanFound(goal, plan)
		// move to PerformAction state
		fsm.PopState()
		fsm.PushState(DoActionState)
	} else {
		agent.PlanFailed(goal)
		// move back to IdleAction state
		fsm.PopState()
		fsm.PushState(IdleState)
	}
}

// MoveToState is the FSM state the agent is in when it must move to a location before doing an
// action
func MoveToState(fsm *goap.FSM, agent goap.Agent) {
	action := agent.CurrentAction()
	if action.RequiresInRange() && action.Target() == nil {
		fmt.Println("Error: Action requires a target but has none. Planning failed. You did not assign the target in your Action.CheckProceduralPrecondition()")
		fsm.PopState() // move
		fsm.PopState() // perform
		fsm.PushState(IdleState)
		return
	}

	// get the agent to move itself
	if agent.MoveAgent(action) {
		fsm.PopState()
	}
}

// DoActionState is the FSM state the agent is in when is doing actions
func DoActionState(fsm *goap.FSM, agent goap.Agent) {
	// no actions to perform
	if !agent.HasActionPlan() {
		fsm.PopState()
		fsm.PushState(IdleState)
		agent.ActionsFinished()
		return
	}

	action := agent.CurrentAction()
	if action.IsDone() {
		agent.PopCurrentAction()
	}

	if agent.HasActionPlan() {
		// perform the next action
		action = agent.CurrentAction()
		inRange := true
		if action.RequiresInRange() {
			inRange = action.IsInRange()
		}
		if inRange {
			// we are in range, so perform the action
			success := action.Perform(agent)
			if !success {
				// action failed, we need to plan again
				fsm.PopState()
				fsm.PushState(IdleState)
				agent.PlanAborted(action)
			}
		} else {
			// we need to move there first
			fsm.PushState(MoveToState)
		}

	} else {
		// no actions left, move to Plan state
		fsm.PopState()
		fsm.PushState(IdleState)
		agent.ActionsFinished()
	}
}
