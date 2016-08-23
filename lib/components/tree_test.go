package components_test

import (
	. "github.com/stojg/cyberspace/lib/components"
	"testing"
)

func TestTreeAdd(t *testing.T) {

	tree := NewTree("asd", -1)

	tree.Add(&Instance{
		Name:       "first.stack",
		InstanceID: "i-1",
	})
	tree.Add(&Instance{
		Name:       "first.rds",
		InstanceID: "i-2",
	})
	tree.Add(&Instance{
		Name:       "second.stack",
		InstanceID: "i-2",
	})

	siblings := tree.Siblings("first.stack")
	if len(siblings) != 2 {
		t.Errorf("Expected 2 siblings, got %d", len(siblings))
	}

	siblings = tree.Siblings("nah.stack")
	if len(siblings) != 0 {
		t.Errorf("Expected 0 siblings, got %d", len(siblings))
	}

	siblings = tree.Siblings("second.stack")
	if len(siblings) != 1 {
		t.Errorf("Expected 1 siblings, got %d", len(siblings))
	}

}
