package core

import (
	"fmt"
	"reflect"
	"strings"

	"time"

	"github.com/stojg/cyberspace/lib/plan"
	"github.com/stojg/steering"
	"github.com/stojg/vector"
)

// NewAgent returns an initialised agent ready for action!
func NewAgent() *Agent {
	a := &Agent{
		statemachine:  plan.NewFSM(),
		workingMemory: NewWorkingMemory(),
		Debug:         false,
	}
	return a
}

// Agent is the core struct that represents an AI entity that plan and execute actions.
type Agent struct {
	Component
	statemachine     *plan.FSM
	availableActions []plan.Action
	state            plan.StateList
	goals            []plan.StateList
	plan             []plan.Action
	Debug            bool
	workingMemory    *WorkingMemory
}

func (a *Agent) AddAction(action plan.Action) {
	a.availableActions = append(a.availableActions, action)
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
		expiry:   time.Now().Add(10 * time.Second),
	}
	a.Memory().AddEntity(ent)
}

func (a *Agent) AvailableActions() []plan.Action {
	return a.availableActions
}

func (a *Agent) SetPlan(plan []plan.Action) {
	a.plan = plan
}

func (a *Agent) Plan() []plan.Action {
	return a.plan
}

func (a *Agent) PopCurrentAction() {
	if len(a.plan) > 0 {
		a.plan = a.plan[1:]
	}
}

func (a *Agent) State() plan.StateList {
	return a.state
}

func (a *Agent) SetState(s plan.StateList) {
	a.state = s
}

func (a *Agent) Goals() []plan.StateList {
	return a.goals
}

func (a *Agent) SetGoals(set []plan.StateList) {
	a.goals = set
}

func (a *Agent) Memory() *WorkingMemory {
	return a.workingMemory
}

func (a *Agent) Replan() {
	a.statemachine.Reset()
}

// PlanFailed is called when there is no sequence of actions could be found for the supplied goal.
// You will need to try another goal
func (a *Agent) PlanFailed(failedGoal []plan.StateList) {
	if a.Debug {
		fmt.Printf("%s #%d: plan failed with goals: %v and state %v, # of actions: %d\n", a.gameObject.name, a.gameObject.ID(), failedGoal, a.State(), len(a.availableActions))
	}
}

// PlanFound is called when a plan was found for the supplied goal. The actions contains the plan
// of actions the agent will perform, in order.
func (a *Agent) PlanFound(goal plan.StateList, actions []plan.Action) {
	if a.Debug {
		fmt.Printf("%s #%d: plan found for goal %s with actions: %+v state: %v\n", a.gameObject.name, a.gameObject.ID(), a.goals, actionsToString(actions), a.State())
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
func (a *Agent) PlanAborted(abortingAction plan.Action, err error) {
	if a.Debug {
		fmt.Printf("%s #%d: plan was aborted by action %+v\n", a.gameObject.name, a.gameObject.ID(), actionToString(abortingAction))
		fmt.Printf("reason: %s\n", err)
	}
}

// Update checks the state machine and updates its if possible
func (a *Agent) Update(elapsed float64) {
	//a.gameObject.Body().AddForce(vector.X())
	a.statemachine.Update(a, func(msg string) {
		if a.Debug {
			fmt.Printf("%s #%d: %s\n", a.gameObject.name, a.gameObject.ID(), msg)
		}
	})
	a.workingMemory.tick()
}

// MoveAgent is when the agent must move towards the target in order for the next action to be able
// to perform. Return true if the Agent is at the target and the next action can perform. False if
// it is not there yet.
func (a *Agent) Move(action plan.Action) (bool, error) {
	target := action.MoveTo()
	if target == nil {
		return true, nil
	}

	obj, ok := target.(*GameObject)
	if !ok {
		return true, fmt.Errorf("Move to taget isnt a game Object, %+v", target)
	}

	arrive := steering.NewArrive(a.gameObject.Body(), obj.transform.Position(), 10).Get()
	look := steering.NewLookWhereYoureGoing(a.gameObject.Body()).Get()

	a.gameObject.Body().AddForce(arrive.Linear())
	a.gameObject.Body().AddTorque(look.Angular())

	return false, nil
}

func actionToString(action plan.Action) string {
	return reflect.TypeOf(action).String()
}

func actionsToString(actions []plan.Action) string {
	var res []string
	for _, action := range actions {
		res = append(res, actionToString(action))
	}
	return "[" + strings.Join(res, ", ") + "]"
}
