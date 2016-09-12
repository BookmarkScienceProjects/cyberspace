package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/vector"
	"github.com/stojg/vivere/lib/components"
	"io"
	"math/rand"
	"sync/atomic"
)

func newLevel() *level {
	lvl := &level{}
	for i := 0; i < 10; i++ {
		obj := spawn("monster")
		t := obj.Transform()
		t.Position().Set(rand.Float64()*100-50, 0, rand.Float64()*100-50)
	}

	for i := 0; i < 200; i++ {
		obj := spawn("food")
		obj.Transform().Position().Set(rand.Float64()*100-50, 0, rand.Float64()*100-50)
	}
	return lvl
}

type level struct {
}

func (l *level) Update(elapsed float64) {
	for _, body := range core.List.Bodies() {
		body.AddForce(vector.NewVector3(rand.Float64()*10-5, 0, rand.Float64()*10-5))
	}

	UpdatePhysics(elapsed)
	UpdateCollisions(elapsed)

}

const (
	_ byte = iota
	instEntityID
	instPosition
	instOrientation
	instType
	instScale
)

func (l *level) draw() *bytes.Buffer {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.LittleEndian, float32(atomic.LoadUint64(&currentFrame)))
	if err != nil {
		Printf("draw() error %s", err)
	}

	for _, graphic := range core.List.Graphics() {
		gameObject := graphic.GameObject()
		if !graphic.IsRendered() {
			serialize(buf, gameObject)
			graphic.SetRendered()
			continue
		}

		body := gameObject.Body()
		if body != nil && body.Awake() {
			serialize(buf, gameObject)
		}
	}
	return buf
}

func (l *level) fullDraw() *bytes.Buffer {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.LittleEndian, float32(atomic.LoadUint64(&currentFrame)))
	if err != nil {
		Printf("fullDraw() error %s", err)
	}
	for _, graphic := range core.List.Graphics() {
		gameObject := graphic.GameObject()
		serialize(buf, gameObject)
	}
	return buf
}

func (l *level) drawDead() *bytes.Buffer {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.LittleEndian, float32(atomic.LoadUint64(&currentFrame)))
	if err != nil {
		Printf("drawDead() error %s", err)
	}

	for _, id := range core.List.Deleted() {
		if err := binaryStream(buf, instEntityID, id); err != nil {
			Printf("binarystream error %s", err)
		}
	}
	core.List.ClearDeleted()

	return buf
}

func serialize(buf *bytes.Buffer, gameObject *core.GameObject) {
	if err := binaryStream(buf, instEntityID, gameObject.ID()); err != nil {
		Printf("binarystream error %s", err)
	}
	if err := binaryStream(buf, instPosition, gameObject.Transform().Position()); err != nil {
		Printf("binarystream error %s", err)
	}
	if err := binaryStream(buf, instOrientation, gameObject.Transform().Orientation()); err != nil {
		Printf("binarystream error %s", err)
	}

	graphic := gameObject.Graphic()
	if graphic != nil {
		if err := binaryStream(buf, instType, graphic.Model()); err != nil {
			Printf("binarystream error %s", err)
		}
	}
	if err := binaryStream(buf, instScale, gameObject.Transform().Scale()); err != nil {
		Printf("binarystream error %s", err)
	}
}

func binaryStream(buf io.Writer, varType byte, value interface{}) error {
	var err error
	if err = binary.Write(buf, binary.LittleEndian, varType); err != nil {
		return err
	}
	switch val := value.(type) {
	case uint8:
		err = binary.Write(buf, binary.LittleEndian, val)
	case uint16:
		err = binary.Write(buf, binary.LittleEndian, float32(val))
	case uint32:
		err = binary.Write(buf, binary.LittleEndian, float32(val))
	case int:
		err = binary.Write(buf, binary.LittleEndian, int32(val))
	case float32:
		err = binary.Write(buf, binary.LittleEndian, val)
	case float64:
		err = binary.Write(buf, binary.LittleEndian, float32(val))
	case components.Entity:
		err = binary.Write(buf, binary.LittleEndian, float32(val))
	case State:
		err = binary.Write(buf, binary.LittleEndian, float32(val))
	case *vector.Vector3:
		err = binary.Write(buf, binary.LittleEndian, float32(val[0]))
		err = binary.Write(buf, binary.LittleEndian, float32(val[1]))
		err = binary.Write(buf, binary.LittleEndian, float32(val[2]))
	case *vector.Quaternion:
		err = binary.Write(buf, binary.LittleEndian, float32(val.R))
		err = binary.Write(buf, binary.LittleEndian, float32(val.I))
		err = binary.Write(buf, binary.LittleEndian, float32(val.J))
		err = binary.Write(buf, binary.LittleEndian, float32(val.K))
	default:
		panic(fmt.Errorf("Havent found out how to serialise literal %v with value of type '%T'", varType, val))
	}
	return err
}
