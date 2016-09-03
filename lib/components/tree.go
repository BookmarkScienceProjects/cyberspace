package components

import (
	"github.com/stojg/cyberspace/lib/formation"
	"github.com/stojg/vector"
	"strings"
	"sync"
)

func NewTree(name string, level int) *TreeNode {
	return &TreeNode{
		level:     level,
		name:      name,
		children:  make([]*TreeNode, 0),
		instances: make([]*AWSInstance, 0),
		manager:   formation.NewManager(formation.NewDefensiveCircle(10, 0)),
	}
}

type TreeNode struct {
	level int
	name  string

	instances []*AWSInstance

	sync.Mutex
	children []*TreeNode
	manager  *formation.Manager
}

func (c *TreeNode) Position() *vector.Vector3 {
	return vector.NewVector3(0, 0, 0)
}

func (c *TreeNode) Orientation() *vector.Quaternion {
	return vector.NewQuaternion(1, 0, 0, 0)
}

func (c *TreeNode) SetTarget(t *vector.Vector3) {

}

func (c *TreeNode) Siblings(name string) []*AWSInstance {
	names := strings.Split(name, ".")

	c.Lock()
	defer c.Unlock()
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

	c.Lock()
	defer c.Unlock()

	// We hit the final level
	if len(names) <= c.level+1 {
		c.instances = append(c.instances, i)
		return
	}

	// check if there already is a child node
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
