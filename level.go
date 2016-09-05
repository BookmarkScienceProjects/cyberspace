package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	. "github.com/stojg/vector"
	. "github.com/stojg/vivere/lib/components"
	"io"
	"sync/atomic"
)

var (
	entities       *EntityManager
	modelList      *ModelList
	collisionList  *CollisionList
	rigidList      *RigidBodyList
	controllerList *ControllerList
)

func newLevel() *level {
	entities = NewEntityManager()
	modelList = NewModelList()
	rigidList = NewRigidBodyManager()
	collisionList = NewCollisionList()
	controllerList = NewControllerList()

	lvl := &level{}
	lvl.systems = append(lvl.systems, &physicSystem{})
	lvl.systems = append(lvl.systems, &controllerSystem{})
	lvl.systems = append(lvl.systems, &collisionSystem{})
	lvl.world = &World{}
	lvl.systems = append(lvl.systems, lvl.world)
	return lvl
}

type level struct {
	systems []system
	world   *World
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

	for _, model := range l.world.list {
		if err := binaryStream(buf, instEntityID, *model.ID()); err != nil {
			Printf("binarystream error %s", err)
		}
		if err := binaryStream(buf, instPosition, model.Position()); err != nil {
			Printf("binarystream error %s", err)
		}
		if err := binaryStream(buf, instOrientation, model.Orientation()); err != nil {
			Printf("binarystream error %s", err)
		}
		if err := binaryStream(buf, instType, model.Kind()); err != nil {
			Printf("binarystream error %s", err)
		}
		if err := binaryStream(buf, instScale, model.Size()); err != nil {
			Printf("binarystream error %s", err)
		}
		if err := binaryStream(buf, instState, model.State()); err != nil {
			Printf("binarystream error %s", err)
		}
	}

	return buf
}

func binaryStream(buf io.Writer, lit literal, value interface{}) error {
	var err error
	if err = binary.Write(buf, binary.LittleEndian, lit); err != nil {
		return err
	}
	switch val := value.(type) {
	case uint8:
		err = binary.Write(buf, binary.LittleEndian, val)
	case uint16:
		err = binary.Write(buf, binary.LittleEndian, float32(val))
	case uint32:
		err = binary.Write(buf, binary.LittleEndian, float32(val))
	case Kind:
		err = binary.Write(buf, binary.LittleEndian, float32(val))
	case float32:
		err = binary.Write(buf, binary.LittleEndian, val)
	case float64:
		err = binary.Write(buf, binary.LittleEndian, float32(val))
	case Entity:
		err = binary.Write(buf, binary.LittleEndian, float32(val))
	case State:
		err = binary.Write(buf, binary.LittleEndian, float32(val))
	case *Vector3:
		err = binary.Write(buf, binary.LittleEndian, float32(val[0]))
		err = binary.Write(buf, binary.LittleEndian, float32(val[1]))
		err = binary.Write(buf, binary.LittleEndian, float32(val[2]))
	case *Quaternion:
		err = binary.Write(buf, binary.LittleEndian, float32(val.R))
		err = binary.Write(buf, binary.LittleEndian, float32(val.I))
		err = binary.Write(buf, binary.LittleEndian, float32(val.J))
		err = binary.Write(buf, binary.LittleEndian, float32(val.K))
	default:
		panic(fmt.Errorf("Havent found out how to serialise literal %v with value of type '%T'", lit, val))
	}
	return err
}
