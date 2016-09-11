package object

import (
	"fmt"
	"github.com/stojg/vivere/lib/components"
	"math"
	"sync"
)

var List list

func init() {

}

type Deletable interface {
	Delete(*components.Entity)
}

type list struct {
	next uint16

	sync.Mutex
	items      [math.MaxUint16]Object
	Deleted    []Object
	Models     Deletable
	Rigids     Deletable
	Collisions Deletable
}

func (l *list) Add(o Object) {
	// dont add anymore, we are full
	if l.next+1 == math.MaxUint16 {
		fmt.Println("stuff list full")
		return
	}
	l.Lock()
	l.items[l.next] = o
	l.next++
	l.Unlock()
}

func (l *list) Remove(i uint16) {
	l.Lock()
	l.Models.Delete(l.items[i].ID())
	l.Rigids.Delete(l.items[i].ID())
	l.Collisions.Delete(l.items[i].ID())

	// this might trigger call backs or something
	//l.items[i].SetState(stateDead)

	// keep a record of deleted entities so they can be flushed to the view
	l.Deleted = append(l.Deleted, l.items[i])
	// Take the last element in the list and replace the object to be deleted with that one
	l.items[i] = l.items[l.next-1]
	// we now have one more spot
	l.next--

	l.Unlock()
}

func (l *list) All() map[uint16]Object {
	result := make(map[uint16]Object, 0)
	l.Lock()
	for i := uint16(0); i < l.next; i++ {
		result[i] = l.items[i]
	}
	l.Unlock()
	return result
}

// idea for the future, let k be a bitmask
func (l *list) OfKind(k Kind) map[uint16]Object {
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

func (l *list) Nearest(monster Object, kind Kind) (uint16, bool) {
	var closestIndex uint16
	found := false
	closestDistance := math.MaxFloat64
	for i := range l.items {
		if monster.Kind() != kind {
			continue
		}
		distance := monster.Position().NewSub(l.items[i].Position()).SquareLength()
		if monster.Position().NewSub(l.items[i].Position()).SquareLength() < closestDistance {
			found = true
			closestIndex = uint16(i)
			closestDistance = distance
		}
	}
	return closestIndex, found
}
