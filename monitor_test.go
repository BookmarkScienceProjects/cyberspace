package main

import (
	"github.com/stojg/vivere/lib/vector"
	"testing"
)

func TestSomething(t *testing.T) {

	instA := &Instance{
		Name:  "A",
		Scale: vector.Vector3{20, 20, 20},
	}
	instB := &Instance{
		Name:  "B",
		Scale: vector.Vector3{10, 10, 10},
	}
	instC := &Instance{
		Name:  "C",
		Scale: vector.Vector3{10, 10, 10},
	}

	rootNode := NewTree("root", -1)
	rootNode.Add(instA)
	rootNode.Add(instB)
	rootNode.Add(instC)

	BuildTree(rootNode)

	t.Errorf("%v", rootNode.MinPoint(0))
	t.Errorf("%v", rootNode.MaxPoint(0))
	//t.Errorf("%v", rootNode.MinPoint(1))
	//t.Errorf("%v", rootNode.MaxPoint(1))
	//t.Errorf("%v", rootNode.MinPoint(2))
	//t.Errorf("%v", rootNode.MaxPoint(2))

	t.Errorf("%v", rootNode.children[0].Leaves()[0].Position())
	t.Errorf("%v", rootNode.children[1].Leaves()[0].Position())
	t.Errorf("%v", rootNode.children[2].Leaves()[0].Position())
}
