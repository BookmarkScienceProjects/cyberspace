package main

import (
	"math"
	"sync"
)

type objectList struct {
	next uint16

	sync.Mutex
	items   [math.MaxUint16]Object
	deleted []Object
}

func (l *objectList) Add(o Object) {
	// dont add anymore, we are full
	if l.next+1 == math.MaxUint16 {
		Println("stuff list full")
		return
	}
	l.Lock()
	l.items[l.next] = o
	l.next++
	l.Unlock()
}

func (l *objectList) Remove(i uint16) {
	l.Lock()
	modelList.Delete(l.items[i].ID())
	rigidList.Delete(l.items[i].ID())
	collisionList.Delete(l.items[i].ID())

	// this might trigger call backs or something
	l.items[i].SetState(stateDead)

	// keep a record of deleted entities so they can be flushed to the view
	l.deleted = append(l.deleted, l.items[i])
	// Take the last element in the list and replace the object to be deleted with that one
	l.items[i] = l.items[l.next-1]
	// we now have one more spot
	l.next--

	l.Unlock()
}

func (l *objectList) All() map[uint16]Object {
	result := make(map[uint16]Object, 0)
	l.Lock()
	for i := uint16(0); i < l.next; i++ {
		result[i] = l.items[i]
	}
	l.Unlock()
	return result
}

// idea for the future, let k be a bitmask
func (l *objectList) ofKind(k Kind) map[uint16]Object {
	result := make(map[uint16]Object, 0)
	l.Lock()
	for i := uint16(0); i < l.next; i++ {
		if l.items[i].Kind() == k {
			result[i] = l.items[i]
		}
	}
	l.Unlock()
	return result
}

func nearest(monster Object, gunks map[uint16]Object) (uint16, bool) {
	var closestIndex uint16
	found := false
	closestDistance := math.MaxFloat64
	for i := range gunks {
		distance := monster.Position().NewSub(gunks[i].Position()).SquareLength()
		if monster.Position().NewSub(gunks[i].Position()).SquareLength() < closestDistance {
			found = true
			closestIndex = i
			closestDistance = distance
		}
	}
	return closestIndex, found
}
