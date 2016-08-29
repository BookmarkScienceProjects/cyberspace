package main

import (
	"fmt"
	"github.com/stojg/cyberspace/lib/components"
	"github.com/stojg/vector"
	"testing"
)

type LevelAI struct {
	list []*ai
}

func (a *LevelAI) Add(b *ai) {
	a.list = append(a.list, b)
}

func (ai *LevelAI) Update(elapsed float64) {

}

func TestSomething(t *testing.T) {
	root := components.NewTree("root", 0)

	inst := &components.AWSInstance{}
	root.Add(inst)

	levelAI := &LevelAI{}
	levelAI.Update(0.01)

}

type TestStruct struct {
	position    *vector.Vector3
	orientation *vector.Quaternion
}

func (t *TestStruct) Update(val float64) {
	t.position.Set(val, val, val)
}

var oMapIntList map[int]*TestStruct
var oMapStringList map[string]*TestStruct
var oSliceList []*TestStruct

func BenchmarkIntMap(b *testing.B) {

	oMapIntList = make(map[int]*TestStruct, 0)
	for i := 0; i < 1000; i++ {
		oMapIntList[i] = &TestStruct{
			position: vector.NewVector3(0, 0, 0),
		}
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for i, _ := range oMapIntList {
			oMapIntList[i].Update(1)
		}
	}
}

func BenchmarkStringMap(b *testing.B) {

	oMapStringList = make(map[string]*TestStruct, 0)
	for i := 0; i < 1000; i++ {
		oMapStringList[fmt.Sprintf("%d", i)] = &TestStruct{
			position: vector.NewVector3(0, 0, 0),
		}
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for i, _ := range oMapStringList {
			oMapStringList[i].Update(1)
		}
	}
}

func BenchmarkSlice(b *testing.B) {

	oSliceList = make([]*TestStruct, 1000)
	for i := range oSliceList {
		oSliceList[i] = &TestStruct{
			position: vector.NewVector3(0, 0, 0),
		}
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for i, _ := range oSliceList {
			oSliceList[i].Update(1)
		}
	}
}
