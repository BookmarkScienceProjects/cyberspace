package main

import (
	"github.com/stojg/vivere/lib/vector"
	"math"
	"strings"
)

type Collidable interface {
	Name() string
	Position() *vector.Vector3
	SetPosition(*vector.Vector3)
	AddPosition(*vector.Vector3)
	Radius() float64
	MinPoint(axis int) float64
	MaxPoint(axis int) float64
	Children() []Collidable
	Leaves() []Collidable
	Add(i *Instance)
	Instance() *Instance
}

type Leaf struct {
	position *vector.Vector3
	scale    *vector.Vector3
	instance *Instance
}

func (l *Leaf) Name() string {
	return l.instance.Name
}

func (l *Leaf) Position() *vector.Vector3 {
	return l.position
}

func (l *Leaf) SetPosition(v *vector.Vector3) {
	l.position = v
}

func (l *Leaf) AddPosition(v *vector.Vector3) {
	l.position.Add(v)
}

func (l *Leaf) Radius() float64 {
	Println(l.instance.Scale[0])
	return l.instance.Scale[0]

	//z := l.instance.Scale[2]
	//return math.Sqrt((x*x)+(z*z))
}

func (l *Leaf) MinPoint(axis int) float64 {
	return l.position[axis] - l.instance.Scale[axis]/2
}

func (l *Leaf) MaxPoint(axis int) float64 {
	return l.position[axis] + l.instance.Scale[axis]/2
}

func (l *Leaf) Children() []Collidable {
	return nil
}

func (l *Leaf) Leaves() []Collidable {
	return nil
}

func (l *Leaf) Add(i *Instance) {}

func (l *Leaf) Instance() *Instance {
	return l.instance
}

func NewTree(name string, level int) *TreeNode {
	return &TreeNode{
		level:    level,
		name:     name,
		children: make([]Collidable, 0),
		leaves:   make([]Collidable, 0),
		position: vector.NewVector3(0, 0, 0),
	}
}

type TreeNode struct {
	parent   Collidable
	level    int
	name     string
	children []Collidable
	leaves   []Collidable
	position *vector.Vector3
}

func (l *TreeNode) SetPosition(v *vector.Vector3) {
	l.position = v
}

func (c *TreeNode) Children() []Collidable {
	return c.children
}

func (c *TreeNode) Leaves() []Collidable {
	return c.leaves
}

func (c *TreeNode) Add(i *Instance) {
	names := strings.Split(i.Name, ".")
	if len(names) <= c.level+1 {
		c.leaves = append(c.leaves, &Leaf{
			instance: i,
			position: vector.NewVector3(0, 0, 0),
		})
		return
	}

	var existingChild Collidable
	for _, child := range c.Children() {
		if child.Name() == names[c.level+1] {
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

func (c *TreeNode) Level() int {
	return c.level
}

func (c *TreeNode) Name() string {
	return c.name
}

func (c *TreeNode) Position() *vector.Vector3 {
	return c.position
}

func (c *TreeNode) AddPosition(p *vector.Vector3) {
	c.position.Add(p)
	for _, child := range c.Children() {
		child.AddPosition(p)
	}
	for _, child := range c.leaves {
		child.AddPosition(p)
	}
}

func (c *TreeNode) Radius() float64 {
	x := c.MaxPoint(0) - c.MinPoint(0)
	z := c.MaxPoint(2) - c.MinPoint(2)
	return math.Sqrt((x * x) + (z * z))
}

func (c *TreeNode) MinPoint(axis int) float64 {
	min := math.MaxFloat64
	for _, child := range c.Children() {
		minPoint := child.MinPoint(axis)
		if minPoint < min {
			min = minPoint
		}
	}
	for _, child := range c.leaves {
		minPoint := child.MinPoint(axis)
		if minPoint < min {
			min = minPoint
		}
	}
	return min * 1.0
}

func (c *TreeNode) MaxPoint(axis int) float64 {
	max := -math.MaxFloat64
	for _, child := range c.Children() {
		maxPoint := child.MaxPoint(axis)
		if maxPoint > max {
			max = maxPoint
		}
	}
	for _, child := range c.leaves {
		maxPoint := child.MaxPoint(axis)
		if maxPoint > max {
			max = maxPoint
		}
	}
	return max * 1.0
}

func (c *TreeNode) Instance() *Instance {
	return nil
}
