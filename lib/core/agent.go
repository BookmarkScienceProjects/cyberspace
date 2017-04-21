package core

import (
	"fmt"

	"github.com/stojg/cyberspace/lib/planning"
	"github.com/stojg/steering"
	"github.com/stojg/vector"
)

// NewAgent returns an initialised agent ready for action!
func NewAgent(actions []planning.Action) *Agent {
	a := &Agent{
		DefaultAgent:  planning.NewDefaultAgent(actions),
		workingMemory: NewWorkingMemory(),
		Debug:         false,
	}
	return a
}

// Agent is the core struct that represents an AI entity that plan and execute actions.
type Agent struct {
	planning.DefaultAgent
	Component
	Debug         bool
	workingMemory *WorkingMemory
}

func (a *Agent) DetectsModality(modality Modality) bool {
	return true
}

func (a *Agent) Position() *vector.Vector3 {
	return a.transform.position
}

func (a *Agent) Orientation() *vector.Quaternion {
	return a.transform.orientation
}

func (a *Agent) Threshold() float64 {
	return 100
}
func (a *Agent) Notify(signal *Signal) {
	other := List.Get(signal.ID)
	if other == nil {
		return
	}

	distance := other.Transform().Position().NewSub(a.Transform().position).Length()
	ent := &Entity{
		ID:       signal.ID,
		Position: signal.Position,
		Distance: distance,
		//Velocity: other.Body().velocity,
		//Name     string
		//Type:
	}
	if !a.Memory().AddEntity(ent) {
		a.Replan()
	}
}

func (a *Agent) Memory() *WorkingMemory {
	return a.workingMemory
}

func (a *Agent) Replan() {
	a.DefaultAgent.StateMachine.Reset(planning.Idle)
}

// PlanFailed is called when there is no sequence of actions could be found for the supplied goal.
// You will need to try another goal
func (a *Agent) PlanFailed(failedGoal planning.StateList) {
	if a.Debug {
		fmt.Printf("%s #%d: plan failed with goalState: %v and state %v\n", a.gameObject.name, a.gameObject.ID(), failedGoal, a.State())
	}
}

// PlanFound is called when a plan was found for the supplied goal. The actions contains the plan
// of actions the agent will perform, in order.
func (a *Agent) PlanFound(goal planning.StateList, actions []planning.Action) {
	if a.Debug {
		fmt.Printf("%s #%d: Plan found with actions: %v for %v\n", a.gameObject.name, a.gameObject.ID(), actions, a.State())
	}
}

// ActionsFinished is signaled when all actions are complete and the goal was reached.
func (a *Agent) ActionsFinished() {
	if a.Debug {
		fmt.Printf("%s #%d: actions finished\n", a.gameObject.name, a.gameObject.ID())
	}
}

// PlanAborted is called when one of the actions in the plan have discovered that it can no longer
// be done.
func (a *Agent) PlanAborted(abortingAction planning.Action) {
	if a.Debug {
		fmt.Printf("%s #%d: plan was aborted by action %s aborted", a.gameObject.name, a.gameObject.ID(), abortingAction.String())
	}
}

// Update checks the state machine and updates its if possible
func (a *Agent) Update() {
	//a.gameObject.Body().AddForce(vector.X())
	a.DefaultAgent.FSM(a, func(msg string) {
		if a.Debug {
			//fmt.Printf("%s #%d: %s\n", a.gameObject.name, a.gameObject.ID(), msg)
		}
	})
	a.workingMemory.tick()
}

// MoveAgent is when the agent must move towards the target in order for the next action to be able
// to perform. Return true if the Agent is at the target and the next action can perform. False if
// it is not there yet.
func (a *Agent) MoveAgent(nextAction planning.Action) bool {
	target, found := nextAction.Target().(*GameObject)
	if !found {
		fmt.Printf("in core.Agent.MoveAgent: %s is not a *GameObject", nextAction.Target())
		return false
	}
	if target == nil {
		fmt.Printf("in core.Agent.MoveAgent: %s is not a *GameObject", nextAction.Target())
		return false
	}

	if nextAction.InRange(a) {
		return true
	}

	arrive := steering.NewArrive(a.gameObject.Body(), target.transform.Position(), 10).Get()
	look := steering.NewLookWhereYoureGoing(a.gameObject.Body()).Get()

	a.gameObject.Body().AddForce(arrive.Linear())
	a.gameObject.Body().AddTorque(look.Angular())

	return false
}
