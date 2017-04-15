package core

import (
	"fmt"
	"github.com/stojg/goap"
	"github.com/stojg/steering"
)

// NewAgent returns an initialised agent ready for action!
func NewAgent(actions []goap.Actionable) *Agent {
	a := &Agent{
		DefaultAgent: goap.NewDefaultAgent(actions),
	}
	return a
}

// Agent is the core struct that represents an AI entity that plan and execute actions.
type Agent struct {
	goap.DefaultAgent
	Component
	Debug bool
}

// PlanFailed is called when there is no sequence of actions could be found for the supplied goal.
// You will need to try another goal
func (a *Agent) PlanFailed(failedGoal goap.StateList) {
	if a.Debug {
		fmt.Printf("plan failed with goalState: %v and state %v\n", failedGoal, a.State())
	}
}

// PlanFound is called when a plan was found for the supplied goal. The actions contains the plan
// of actions the agent will perform, in order.
func (a *Agent) PlanFound(goal goap.StateList, actions []goap.Actionable) {
	if a.Debug {
		fmt.Printf("Plan found with actions: %v for %v\n", actions, a.State())
	}
}

// ActionsFinished is signaled when all actions are complete and the goal was reached.
func (a *Agent) ActionsFinished() {
	if a.Debug {
		fmt.Println("actions finished")
	}
}

// PlanAborted is called when one of the actions in the plan have discovered that it can no longer
// be done.
func (a *Agent) PlanAborted(abortingAction goap.Actionable) {
	if a.Debug {
		fmt.Printf("plan was aborted by action %s aborted", abortingAction)
	}
}

// Update checks the state machine and updates its if possible
func (a *Agent) Update() {
	a.DefaultAgent.FSM(a, func(msg string) {
		if a.Debug {
			fmt.Println(msg)
		}
	})
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

	arrive := steering.NewArrive(a.gameObject.Body(), target.transform.Position(), 10).Get()
	look := steering.NewLookWhereYoureGoing(a.gameObject.Body()).Get()

	a.gameObject.Body().AddForce(arrive.Linear())
	a.gameObject.Body().AddTorque(look.Angular())

	return false
}
