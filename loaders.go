package main

import (
	"encoding/json"
	"fmt"
	"github.com/stojg/vector"
	"io/ioutil"
	"os"
)

type fileModel struct {
	Kind   Kind            `json:"kind"`
	Scale  *vector.Vector3 `json:"scale"`
	Weight float64         `json:"weight"`
}

func loadFromFile(name string) *fileModel {
	file, e := ioutil.ReadFile("./data/" + name + ".json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	var obj *fileModel
	err := json.Unmarshal(file, &obj)
	if err != nil {
		Println(err)
		return nil
	}
	return obj
}
