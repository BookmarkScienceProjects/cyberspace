package plan

import (
	"fmt"
	"strings"
)

// State is a key in a state list, it's a string for easier debugging
type State string

// StateList is a set of key / value pairs. The key is a State and the value is an empty interface
type StateList map[State]interface{}

// Set a key value pair, will overwrite current value if it already exists
func (s StateList) Set(n State, v interface{}) { s[n] = v }

// Get a value from the StateList
func (s StateList) Get(n State) interface{} { return s[n] }

// Unset a state
func (s StateList) Unset(n State) { delete(s, n) }

// Have checks if StateList have the state, doesn't compare values
func (s StateList) Have(n State) bool {
	_, found := s[n]
	return found
}

// Clone copies and returns a new list
func (s StateList) Clone() StateList {
	state := make(StateList, len(s))
	for key, s := range s {
		state[key] = s
	}
	return state
}

// Apply inserts and updates the StateList with the key and values from another StateList
func (s StateList) Apply(n StateList) StateList {
	for key, change := range n {
		s[key] = change
	}
	return s
}

// Includes tests if the StateList have all states from the test StateList
func (s StateList) Includes(test StateList) bool {
	for testKey, _ := range test {
		if !s.Have(testKey) {
			return false
		}
	}
	return true
}

// Includes tests if the StateList have all states from the test StateList
func (s StateList) Compare(test StateList) bool {
	for testKey, val := range test {
		if s.Get(testKey) != val {
			return false
		}
	}
	return true
}

// String returns a string representation of the StateList, note that since the backing storage is a map there is no
// guarantee about order
func (s StateList) String() string {
	var res []string
	for k, v := range s {
		res = append(res, fmt.Sprintf("%v: %v", k, v))
	}
	return "[" + strings.Join(res, ", ") + "]"
}
