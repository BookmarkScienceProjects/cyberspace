package components

import (
	"strings"
)

func NewTree(name string, level int) *TreeNode {
	return &TreeNode{
		level:    level,
		name:     name,
		children: make([]*TreeNode, 0),
	}
}

type TreeNode struct {
	parent *TreeNode

	level int
	name  string

	instance *AWSInstance

	children []*TreeNode
}

func (node *TreeNode) Children() []*TreeNode {
	return node.children
}

func (node *TreeNode) Siblings() []*AWSInstance {
	if node.parent == nil {
		return nil
	}
	return node.parent.Leafs()
}

func (node *TreeNode) Leafs() []*AWSInstance {
	var i []*AWSInstance
	for _, child := range node.children {
		if child.instance != nil {
			i = append(i, child.instance)
		} else {
			i = append(i, child.Leafs()...)
		}
	}
	return i
}

func (node *TreeNode) Name() string {
	return node.name
}

func (node *TreeNode) Parent() *TreeNode {
	return node.parent
}

func (node *TreeNode) Instance() *AWSInstance {
	return node.instance
}

func (node *TreeNode) Add(instance *AWSInstance) *TreeNode {
	names := strings.Split(instance.Name, ".")

	// check if we are at the leaf
	if len(names) == node.level {
		node.instance = instance
		return node
	}

	// try to find a pre existing child node
	for _, childNode := range node.Children() {
		if childNode.name == names[node.level] {
			return childNode.Add(instance)
		}
	}

	child := NewTree(names[node.level], node.level+1)
	child.parent = node

	node.children = append(node.children, child)
	return child.Add(instance)
}
