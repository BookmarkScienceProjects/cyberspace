// Package plan is a Goal Orientated Action Planner
package plan

// Find what sequence of actions can fulfill the goal. Returns an empty list if a plan could not be found, or
// a list of the actions that must be performed, in order, to fulfill the goal.
// @todo if there is performance problem this planner could use an iterative deepening A* search instead of a linear
func Find(agent Agent, worldState StateList, goalState StateList) []Action {
	// we cannot plan without any goals
	if len(goalState) == 0 {
		panic("cannot plan without a goal")
	}

	var result []Action

	// check what actions can run
	var usableActions []Action
	for _, action := range agent.AvailableActions() {
		action.Reset()
		usableActions = append(usableActions, action)
	}

	// early out, this agent doesn't have any actions
	if len(usableActions) == 0 {
		return result
	}

	var plans []*node
	if !buildGraph(&node{state: worldState}, &plans, usableActions, goalState, agent) {
		return result
	}

	// get the cheapest plan
	cheapest := plans[0]
	for i := 1; i < len(plans); i++ {
		if plans[i].runningCost < cheapest.runningCost {
			cheapest = plans[i]
		}
	}

	// invert the list so that we get the end action at the end of the result list
	for n := cheapest; n != nil && n.action != nil; n = n.parent {
		result = append([]Action{n.action}, result...)
	}
	return result
}

// Node is used for building up the graph and holding the running costs of actions.
type node struct {
	parent      *node
	runningCost float64
	state       StateList
	action      Action
}

// buildGraph returns true if at least one solution was found. The possible paths are stored in the
// leaves list. Each leaf has a 'runningCost' value where the lowest cost will be the best action
// sequence.
func buildGraph(parent *node, result *[]*node, actions []Action, goal StateList, agent Agent) bool {

	for i, action := range actions {

		if !parent.state.Includes(action.Preconditions()) {
			continue
		}

		candidate := &node{
			parent:      parent,
			runningCost: parent.runningCost + action.Cost(),
			state:       parent.state.Clone().Apply(action.Effects()),
			action:      action,
		}

		if !action.CheckContextPrecondition(candidate.state) {
			continue
		}

		// found a solution
		if candidate.state.Includes(goal) {
			*result = append(*result, candidate)
			continue
		}

		// not at a solution yet, so test all the remaining actions and branch out the tree
		buildGraph(candidate, result, remove(actions, i), goal, agent)
	}
	return len(*result) > 0
}
