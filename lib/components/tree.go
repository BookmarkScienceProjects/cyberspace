package components

import (
	"strings"
)

func NewTree(name string, level int) *TreeNode {
	return &TreeNode{
		level:     level,
		name:      name,
		children:  make([]*TreeNode, 0),
		instances: make([]*AWSInstance, 0),
	}
}

type TreeNode struct {
	level     int
	name      string
	children  []*TreeNode
	instances []*AWSInstance
}

func (c *TreeNode) Siblings(name string) []*AWSInstance {
	names := strings.Split(name, ".")

	if len(names) == 1 {
		var sib []*AWSInstance
		for _, child := range c.children {
			sib = append(sib, child.instances...)
		}
		return sib
	}
	for _, child := range c.children {
		if child.name == names[0] {
			return child.Siblings(strings.Join(names[1:], "."))
		}
	}
	return nil
}

func (c *TreeNode) Add(i *AWSInstance) {
	names := strings.Split(i.Name, ".")

	if len(names) <= c.level+1 {
		c.instances = append(c.instances, i)
		return
	}

	var existingChild *TreeNode
	for _, child := range c.children {
		if child.name == names[c.level+1] {
			existingChild = child
			break
		}
	}

	if existingChild == nil {
		existingChild = NewTree(names[c.level+1], c.level+1)
		c.children = append(c.children, existingChild)
	}
	existingChild.Add(i)
}
