package core

import (
	"fmt"
	"github.com/stojg/goap"
	"github.com/stojg/steering"
)

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

type Agent struct {
	Component
	fsm        *goap.FSM
	worldState goap.StateList
	innerState goap.StateList
	goals      goap.StateList

	availableActions []goap.Actionable
	currentActions   []goap.Actionable

	havePlan   bool
	timeInPlan float64
}

func (a *Agent) StateMachine() *goap.FSM {
	return a.fsm
}

func (a *Agent) AvailableActions() []goap.Actionable {
	return a.availableActions
}

func (a *Agent) SetAvailableActions(actions []goap.Actionable) {
	a.availableActions = actions
}

func (a *Agent) CurrentAction() goap.Actionable {
	return a.currentActions[0]
}

func (a *Agent) PopCurrentAction() {
	a.currentActions = a.currentActions[1:]
}

func (a *Agent) SetCurrentActions(actions []goap.Actionable) {
	a.currentActions = actions
}

func (a *Agent) HasActionPlan() bool {
	return len(a.currentActions) > 0
}

// The starting state of the Agent and the world.
// Supply what states are needed for actions to run.
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

func (a *Agent) SetWorldState(state goap.StateList) {
	a.worldState = state
}

func (a *Agent) SetState(k string, v interface{}) {
	a.innerState[k] = v
}

// Give the planner a new goal so it can figure out
// the actions needed to fulfill it.
func (a *Agent) CreateGoalState() goap.StateList {
	return a.goals
}

func (a *Agent) AddGoal(name string, value interface{}) {
	a.goals[name] = value
}

// No sequence of actions could be found for the supplied goal.
// You will need to try another goal
func (a *Agent) PlanFailed(failedGoal goap.StateList) {

}

// A plan was found for the supplied goal.
// These are the actions the Agent will perform, in order.
func (a *Agent) PlanFound(goal goap.StateList, actions []goap.Actionable) {
}

// All actions are complete and the goal was reached. Hooray!
func (a *Agent) ActionsFinished() {

}

// One of the actions caused the plan to abort.
func (a *Agent) PlanAborted(abortingAction goap.Actionable) {
	fmt.Printf("plan was aborted by action %s aborted", abortingAction)
}

// Called during Update. Move the agent towards the target in order
// for the next action to be able to perform.
// Return true if the Agent is at the target and the next action can perform.
// False if it is not there yet.
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

	arrive := steering.NewArrive(a.gameObject.Body(), target.transform.Position(), 2, 0.1, 2).Get()
	look := steering.NewLookWhereYoureGoing(a.gameObject.Body()).Get()

	a.gameObject.Body().AddForce(arrive.Linear())
	a.gameObject.Body().AddTorque(look.Angular())

	return false
}

// during IdleState When an agent is in idle state is will do planning
func IdleState(fsm *goap.FSM, agent goap.Agent) {
	worldState := agent.GetWorldState()
	goal := agent.CreateGoalState()
	plan := goap.Plan(agent, agent.AvailableActions(), worldState, goal)
	if plan != nil {
		agent.SetCurrentActions(plan)
		agent.PlanFound(goal, plan)
		// move to PerformAction state
		fsm.PopState()
		fsm.PushState(DoAction)
	} else {
		agent.PlanFailed(goal)
		// move back to IdleAction state
		fsm.PopState()
		fsm.PushState(IdleState)
	}
}

//
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

func DoAction(fsm *goap.FSM, agent goap.Agent) {
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
