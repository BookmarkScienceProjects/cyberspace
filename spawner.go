package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/stojg/cyberspace/lib/core"
	"github.com/stojg/vector"
)

type objectModel struct {
	Model           int             `json:"model"`
	Scale           *vector.Vector3 `json:"scale"`
	Weight          float64         `json:"weight"`
	MaxAcceleration *vector.Vector3 `json:"max_acceleration"`
}

func loadFromFile(name string) *objectModel {
	file, e := ioutil.ReadFile("./data/" + name + ".json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	var obj *objectModel
	err := json.Unmarshal(file, &obj)
	if err != nil {
		Println(err)
		return nil
	}
	return obj
}

func spawn(name string) *core.GameObject {
	data := loadFromFile(name)
	if data == nil {
		return nil
	}
	object := core.NewGameObject(name, core.List)
	object.AddTags([]string{name})
	object.Transform().Position().Set(0, 0, 0)
	object.Transform().Scale().Set(data.Scale[0], data.Scale[1], data.Scale[2])

	object.AddGraphic(core.NewGraphic(data.Model))
	object.AddBody(core.NewBody(1/data.Weight, true))
	if data.MaxAcceleration != nil {
		object.Body().MaxAcceleration().Copy(data.MaxAcceleration)
	}
	object.AddInventory(core.NewInventory())
	object.AddCollision(core.NewCollisionRectangle(data.Scale[0], data.Scale[1], data.Scale[2]))
	return object
}

func spawnNonCollidable(name string) *core.GameObject {
	data := loadFromFile(name)
	if data == nil {
		return nil
	}
	object := core.NewGameObject(name, core.List)
	object.AddTags([]string{name})
	object.Transform().Position().Set(0, 0, 0)
	object.Transform().Scale().Set(data.Scale[0], data.Scale[1], data.Scale[2])

	object.AddGraphic(core.NewGraphic(data.Model))
	object.AddBody(core.NewBody(1/data.Weight, false))
	if data.MaxAcceleration != nil {
		object.Body().MaxAcceleration().Copy(data.MaxAcceleration)
	}
	object.AddInventory(core.NewInventory())
	object.AddCollision(core.NewCollisionRectangle(data.Scale[0], data.Scale[1], data.Scale[2]))
	return object
}
