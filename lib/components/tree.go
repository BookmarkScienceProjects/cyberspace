package components

import (
	"github.com/stojg/formation"
	"github.com/stojg/vector"
	"strings"
)

func NewTree(name string, level int) *TreeNode {
	pattern := formation.NewDefensiveCircle(10, 0).(*formation.DefensiveCirclePattern)
	return &TreeNode{
		level:     level,
		name:      name,
		children:  make([]*TreeNode, 0),
		pattern:   pattern,
		formation: formation.NewManager(pattern),
	}
}

type TreeNode struct {
	parent *TreeNode

	formation *formation.Manager
	pattern   *formation.DefensiveCirclePattern
	position  *vector.Vector3

	level int
	name  string

	instances []*AWSInstance

	children []*TreeNode
}

func (node *TreeNode) Orientation() *vector.Quaternion {
	return vector.NewQuaternion(1, 0, 0, 0)
}

func (node *TreeNode) Position() *vector.Vector3 {

	anchor := vector.NewVector3(0, 0, 0)

	leafs := node.Leafs()

	if len(leafs) == 0 {
		return anchor
	}

	for _, instance := range leafs {
		if instance.Model != nil {
			anchor.Add(instance.Position())
		}
	}
	anchor.Scale(1 / float64(len(leafs)))
	return anchor
}

func (node *TreeNode) SetTarget(t formation.Static) {
	//fmt.Printf("%s %s\n", node.Name(), t.Position())
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
		if len(child.instances) > 0 {
			i = append(i, child.instances...)
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

func (node *TreeNode) Instances() []*AWSInstance {
	return node.instances
}

func (node *TreeNode) Update(elapsed float64) {
	for _, child := range node.Children() {
		child.Update(elapsed)
	}
	node.formation.UpdateSlots()
}

func (node *TreeNode) Add(instance *AWSInstance) *TreeNode {
	names := strings.Split(instance.Name, ".")

	// check if we are at the leaf
	if len(names) == node.level {
		node.instances = append(node.instances, instance)
		if instance.Model != nil {
			node.formation.AddCharacter(instance)
		}
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

	node.formation.AddCharacter(child)

	node.children = append(node.children, child)
	return child.Add(instance)
}
