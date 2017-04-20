package main

import (
	"bytes"
	"encoding/binary"

	"testing"

	"github.com/stojg/cyberspace/lib/core"
)

func TestBinaryStreamCoreID(t *testing.T) {
	var b bytes.Buffer
	err := binaryStream(&b, instEntityID, core.ID(13))

	if err != nil {
		t.Error(err)
		return
	}

	var varType byte
	err = binary.Read(&b, binary.LittleEndian, &varType)
	if err != nil {
		t.Error(err)
		return
	}
	if varType != instEntityID {
		t.Error("expected varType to be of type InstEntityID")
		return
	}
	// we expect core.ID to be converted to a float32, because.. javascript
	var val float32
	err = binary.Read(&b, binary.LittleEndian, &val)
	if err != nil {
		t.Error(err)
		return
	}
	if val != 13 {
		t.Errorf("expected val to be %f, got %f\n", 13.0, val)
		return
	}
}
