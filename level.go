package main

import (
	"bytes"
	"encoding/binary"
	. "github.com/stojg/vivere/lib/components"
	"math"
	"github.com/stojg/vivere/lib/vector"
	"math/rand"
)

var (
	entities       *EntityManager
	modelList      *ModelList
	collisionList  *CollisionList
	rigidList      *RigidBodyList
	controllerList *ControllerList
)

func NewLevel(monitor *Monitor) *Level {

	x := 1800.0
	y := 1800.0
	z := 1800.0

	entities = NewEntityManager()
	modelList = NewModelList()
	rigidList = NewRigidBodyManager()
	collisionList = NewCollisionList()
	controllerList = NewControllerList()

	monitor.UpdateInstances()

	cubeRoot := math.Pow(float64(len(monitor.subnets)), 1.0/3)
	subX := x / cubeRoot
	subY := y / cubeRoot
	subZ := z / cubeRoot

	currX := 0.0
	currY := 0.0
	currZ := 0.0

	var dudeList []*Entity

	for _, subnet := range monitor.subnets {

		instances := subnet.Instances

		midX := currX + subX/2
		midY := currY + subY/2
		midZ := currZ + subZ/2

		cX := 0.2 * subX
		cY := 0.2 * subY
		cZ := 0.2 * subZ

		for _, inst := range instances {

			e := entities.Create()
			dudeList = append(dudeList, e)

			var model EntityType = 1
			if inst.State != "running" {
				model = 0
			}
			body := modelList.New(e, inst.Scale[0], inst.Scale[1], inst.Scale[2], model)
			posX := (midX + (cX * rand.Float64()) - cX/2) - x / 2
			posY := (midY + (cY * rand.Float64()) - cY/2)- y / 2
			posZ := (midZ + (cZ * rand.Float64()) - cZ/2) - z / 2
			body.Position.Set(posX, posY, posZ)
			phi := rand.Float64() * math.Pi * 2
			body.Orientation.RotateByVector(&vector.Vector3{math.Cos(phi), 0, math.Sin(phi)})
			rig := rigidList.New(e, 1)
			rig.MaxAcceleration = &vector.Vector3{10, 0, 10}
			collisionList.New(e, inst.Scale[0], inst.Scale[1], inst.Scale[2])
		}

		currX += subX
		if currX > x {
			currX = 0
			currY += subY
		}
		if currY > y {
			currY = 0
			currZ += subZ
		}
		if currZ > z {
			currZ = 0
		}

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
