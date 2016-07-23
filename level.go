package main

import (
	"bytes"
	"encoding/binary"
	. "github.com/stojg/vivere/lib/components"
	"github.com/stojg/vivere/lib/vector"
	"math"
	"math/rand"
)

var (
	entities      *EntityManager
	modelList     *ModelList
	collisionList *CollisionList
	rigidList     *RigidBodyList
	controllerList *ControllerList
)

func NewLevel(monitor *Monitor) *Level {

	x := 800.0
	y := 800.0
	z := 400.0

	entities = NewEntityManager()
	modelList = NewModelList()
	rigidList = NewRigidBodyManager()
	collisionList = NewCollisionList()
	controllerList = NewControllerList()

	monitor.UpdateInstances()

	var dudeList []*Entity
	for range monitor.instances {
		e := entities.Create()
		dudeList = append(dudeList, e)

		body := modelList.New(e, 8, 8, 8, ENTITY_PRAY)
		body.Position.Set(x*rand.Float64()-x/2, z*rand.Float64()-z/2, rand.Float64()*y-y/2)
		phi := rand.Float64() * math.Pi * 2
		body.Orientation.RotateByVector(&vector.Vector3{math.Cos(phi), 0, math.Sin(phi)})

		rig := rigidList.New(e, 1)
		rig.MaxAcceleration = &vector.Vector3{10, 0, 10}

		collisionList.New(e, 8, 8, 8)
		//controllerList.New(e, NewAI(e))
	}

	lvl := &Level{}
	lvl.systems = append(lvl.systems, &PhysicSystem{})
	//@todo make an AI component?
	lvl.systems = append(lvl.systems, &ControllerSystem{})
	lvl.systems = append(lvl.systems, &CollisionSystem{})
	return lvl
}

type Level struct {
	systems []System
}

func (l *Level) Update(elapsed float64) {
	for i := range l.systems {
		l.systems[i].Update(elapsed)
	}
}

func (l *Level) Draw() *bytes.Buffer {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, float32(Frame))

	for id, component := range modelList.All() {
		binaryStream(buf, INST_ENTITY_ID, *id)
		binaryStream(buf, INST_SET_POSITION, component.Position)
		binaryStream(buf, INST_SET_ORIENTATION, component.Orientation)
		binaryStream(buf, INST_SET_TYPE, component.Model)
		binaryStream(buf, INST_SET_SCALE, component.Scale)
	}

	return buf
}
