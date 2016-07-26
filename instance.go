package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/stojg/vivere/lib/components"
	"github.com/stojg/vivere/lib/vector"
	"math"
	"strings"
)

type Instance struct {
	ID               *Entity
	Cluster          string
	Stack            string
	Environment      string
	InstanceID       string
	InstanceType     string
	Scale            vector.Vector3
	State            string
	Name             string
	CPUUtilization   float64
	CPUCreditBalance float64
}

func (inst *Instance) Update(ec2Inst *ec2.Instance) {

	inst.InstanceID = *ec2Inst.InstanceId
	inst.InstanceType = *ec2Inst.InstanceType
	inst.State = *ec2Inst.State.Name

	for _, tag := range ec2Inst.Tags {
		if *tag.Key == "Name" && len(*tag.Value) > 0 {
			inst.Name = *tag.Value
			nameParts := strings.Split(inst.Name, ".")
			if (len(nameParts)) > 2 {
				inst.Environment = nameParts[2]
			}
			if (len(nameParts)) > 1 {
				inst.Stack = nameParts[1]
			}
			if (len(nameParts)) > 0 {
				inst.Cluster = nameParts[0]
			}
			break
		}
	}

	if t, ok := typeToCost[*ec2Inst.InstanceType]; !ok {
		Printf("No typeToCost found for '%s'", *ec2Inst.InstanceType)
		inst.Scale = vector.Vector3{10, 10, 10}
	} else {
		costToDimension := t * 10000
		size := math.Pow(costToDimension, 1/3.0)
		inst.Scale = vector.Vector3{size, size, size}
	}
}

func (i *Instance) String() string {
	return fmt.Sprintf("%s %s %s\t%s\t%s", i.Cluster, i.Stack, i.Environment, i.InstanceType, i.InstanceID)
}
