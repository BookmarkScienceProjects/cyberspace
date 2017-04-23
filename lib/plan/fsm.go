package plan

type FSMState func(fsm *FSM, obj Agent, debug func(string))

func NewFSM() *FSM {
	fsm := &FSM{}
	fsm.reset(doPlan)
	return fsm
}

type FSM struct {
	stateStack []FSMState
}

func (fsm *FSM) Reset() {
	fsm.reset(doPlan)
}

func (fsm *FSM) Update(agent Agent, debug func(string)) {
	if len(fsm.stateStack) > 0 {
		fsm.stateStack[len(fsm.stateStack)-1](fsm, agent, debug)
	}
}

func (fsm *FSM) Push(state FSMState) {
	fsm.stateStack = append(fsm.stateStack, state)
}

func (fsm *FSM) reset(state FSMState) {
	states := len(fsm.stateStack)
	for i := 0; i < states; i++ {
		fsm.Pop()
	}
	fsm.Push(state)
}

func (fsm *FSM) Pop() {
	if len(fsm.stateStack) > 0 {
		fsm.stateStack = fsm.stateStack[:len(fsm.stateStack)-1]
	}
}

func doPlan(fsm *FSM, agent Agent, debug func(string)) {

	for _, goal := range agent.Goals() {
		actions := Find(agent, agent.State(), goal)
		if len(actions) == 0 {
			continue
		}
		agent.SetPlan(actions)
		if debugger, ok := agent.(Debugger); ok {
			debugger.PlanFound(goal, actions)
		}
		fsm.reset(doAction)
		return
	}

	if debugger, ok := agent.(Debugger); ok {
		debugger.PlanFailed(agent.Goals())
	}
}

func doAction(fsm *FSM, agent Agent, debug func(string)) {
	// no actions to perform
	if len(agent.Plan()) == 0 {
		fsm.reset(doPlan)
		if debugger, ok := agent.(Debugger); ok {
			debugger.ActionsFinished()
		}
		return
	}

	action := agent.Plan()[0]

	// we need to move there first
	if target := action.MoveTo(); target != nil {
		fsm.Push(doMove)
		return
	}

	done, err := action.Run(agent)
	// action failed, we need to plan again
	if err != nil {
		fsm.reset(doPlan)
		if debugger, ok := agent.(Debugger); ok {
			debugger.PlanAborted(action, err)
		}
	}

	if done {
		agent.PopCurrentAction()
	}
}

func doMove(fsm *FSM, agent Agent, debug func(string)) {
	action := agent.Plan()[0]
	done, err := agent.Move(action)
	if err != nil {
		fsm.reset(doPlan)
		return
	}
	if done {
		fsm.Pop()
	}
}
