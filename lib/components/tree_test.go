package components_test

import (
	. "github.com/stojg/cyberspace/lib/components"
	"testing"
)

func TestTreeNode_Add(t *testing.T) {
	instance := &AWSInstance{
		Name:       "first.stack.child1",
		InstanceID: "i-1",
	}

	tree := NewTree("root", 0)
	node := tree.Add(instance)

	if node.Name() != "child1" {
		t.Errorf("Expected name 'child1', got %s", node.Name())
	}

	if node.Instance().InstanceID != instance.InstanceID {
		t.Error("expected the node instance to be the same as the added instance")
	}
}

func TestTreeNode_Parent(t *testing.T) {
	instances := []*AWSInstance{
		{Name: "first.child1", InstanceID: "i-1"},
		{Name: "first.child2", InstanceID: "i-2"},
		{Name: "second.stack.child1", InstanceID: "i-3"},
		{Name: "second.stack.child2", InstanceID: "i-4"},
	}
	tree := NewTree("root", 0)
	var nodes []*TreeNode
	for _, inst := range instances {
		nodes = append(nodes, tree.Add(inst))
	}

	expected := "first"
	parent := nodes[0].Parent()
	if parent.Name() != expected {
		t.Errorf("Expected name '%s', got '%s'", expected, parent.Name())
	}

	expected = "stack"
	parent = nodes[3].Parent()
	if parent.Name() != expected {
		t.Errorf("Expected name '%s', got '%s'", expected, parent.Name())
	}
}

func TestTreeNode_Siblings(t *testing.T) {
	instances := []*AWSInstance{
		{Name: "first.stack.inst1", InstanceID: "i-1"},
		{Name: "first.stack.inst2", InstanceID: "i-2"},
		{Name: "second.stack", InstanceID: "i-3"},
	}

	tree := NewTree("root", 0)

	var nodes []*TreeNode

	for _, inst := range instances {
		nodes = append(nodes, tree.Add(inst))
	}

	expected := 2
	siblings := nodes[0].Siblings()
	if len(siblings) != expected {
		t.Errorf("Expected %d siblings, got %d", expected, len(siblings))
		return
	}

	siblings = tree.Siblings()
	if len(siblings) != 0 {
		t.Errorf("Expected 0 siblings, got %d", len(siblings))
	}
}

func TestTreeNode_Children(t *testing.T) {
	instances := []*AWSInstance{
		{Name: "stack1.prod1", InstanceID: "i-1"},
		{Name: "stack1.prod2", InstanceID: "i-2"},
		{Name: "stack2.prod1", InstanceID: "i-3"},
		{Name: "stack2.prod1", InstanceID: "i-4"},
	}

	tree := NewTree("root", 0)
	for _, inst := range instances {
		tree.Add(inst)
	}

	c := tree.Children()
	if len(c) != 2 {
		t.Errorf("expected 2 children, got %d", len(c))
	}

}
