package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	. "github.com/stojg/cyberspace/lib/components"
	. "github.com/stojg/vector"
	. "github.com/stojg/vivere/lib/components"
	"io"
	"sync/atomic"
	"time"
)

var (
	entities       *EntityManager
	modelList      *ModelList
	collisionList  *CollisionList
	rigidList      *RigidBodyList
	controllerList *ControllerList
)

func newLevel(monitor *awsMonitor) *level {
	entities = NewEntityManager()
	modelList = NewModelList()
	rigidList = NewRigidBodyManager()
	collisionList = NewCollisionList()
	controllerList = NewControllerList()

	ticker := time.NewTicker(time.Second * 60)
	rootNode := NewTree("root", -1)
	go func() {
		for {
			Println("Updating instances")
			monitor.UpdateInstances(rootNode)
			Println("Instances updated")
			<-ticker.C
		}
	}()

	lvl := &level{}
	lvl.systems = append(lvl.systems, &physicSystem{})
	lvl.systems = append(lvl.systems, &controllerSystem{})
	lvl.systems = append(lvl.systems, &collisionSystem{})
	return lvl
}

type level struct {
	systems []system
}

func (l *level) Update(elapsed float64) {
	for i := range l.systems {
		l.systems[i].Update(elapsed)
	}
}

func (l *level) Draw() *bytes.Buffer {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.LittleEndian, float32(atomic.LoadUint64(&currentFrame)))
	if err != nil {
		Printf("Draw() error %s", err)
	}

	for id, component := range modelList.All() {
		if err := binaryStream(buf, instEntityID, *id); err != nil {
			Printf("binarystream error %s", err)
		}
		if err := binaryStream(buf, instPosition, component.Position()); err != nil {
			Printf("binarystream error %s", err)
		}
		if err := binaryStream(buf, instOrientation, component.Orientation()); err != nil {
			Printf("binarystream error %s", err)
		}
		if err := binaryStream(buf, instType, component.Model); err != nil {
			Printf("binarystream error %s", err)
		}
		if err := binaryStream(buf, instScale, component.Scale); err != nil {
			Printf("binarystream error %s", err)
		}
		inst := monitor.FindByEntityID(*id)
		if err := binaryStream(buf, instHealth, inst.Health()); err != nil {
			Printf("binarystream error %s", err)
		}
	}

	return buf
}

func binaryStream(buf io.Writer, lit literal, val interface{}) error {
	var err error
	if err = binary.Write(buf, binary.LittleEndian, lit); err != nil {
		return err
	}
	switch val.(type) {
	case uint8:
		err = binary.Write(buf, binary.LittleEndian, val.(uint8))
	case uint16:
		err = binary.Write(buf, binary.LittleEndian, float32(val.(uint16)))
	case uint32:
		err = binary.Write(buf, binary.LittleEndian, float32(val.(uint32)))
	case EntityType:
		err = binary.Write(buf, binary.LittleEndian, float32(val.(EntityType)))
	case float32:
		err = binary.Write(buf, binary.LittleEndian, val.(float32))
	case float64:
		err = binary.Write(buf, binary.LittleEndian, float32(val.(float64)))
	case Entity:
		err = binary.Write(buf, binary.LittleEndian, float32(val.(Entity)))
	case *Vector3:
		err = binary.Write(buf, binary.LittleEndian, float32(val.(*Vector3)[0]))
		err = binary.Write(buf, binary.LittleEndian, float32(val.(*Vector3)[1]))
		err = binary.Write(buf, binary.LittleEndian, float32(val.(*Vector3)[2]))
	case *Quaternion:
		err = binary.Write(buf, binary.LittleEndian, float32(val.(*Quaternion).R))
		err = binary.Write(buf, binary.LittleEndian, float32(val.(*Quaternion).I))
		err = binary.Write(buf, binary.LittleEndian, float32(val.(*Quaternion).J))
		err = binary.Write(buf, binary.LittleEndian, float32(val.(*Quaternion).K))
	default:
		panic(fmt.Errorf("Havent found out how to serialise literal %v with value of type '%T'", lit, val))
	}
	return err
}
