package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stojg/vivere/lib/vector"
)

var typeToCost map[string]float64

func init() {
	typeToCost = make(map[string]float64)

	typeToCost["t2.nano"] = 0.01
	typeToCost["t2.micro"] = 0.02
	typeToCost["t2.small"] = 0.04
	typeToCost["t2.medium"] = 0.08

	typeToCost["m4.2xlarge"] = 0.336
	typeToCost["m3.large"] = 0.186

	typeToCost["c3.large"] = 0.132
	typeToCost["c4.large"] = 0.137

	typeToCost["t1.micro"] = 0.02
	typeToCost["m1.small"] = 0.058
	typeToCost["m1.medium"] = 0.117

}

type Subnet struct {
	Name      string
	Instances []*Instance
}

type Instance struct {
	InstanceType string
	Scale        vector.Vector3
	State        string
}

type Monitor struct {
	instances map[string]*Instance
	subnets   map[string]*Subnet
}

func (m *Monitor) UpdateInstances() {

	svc := ec2.New(session.New(), &aws.Config{Region: aws.String("ap-southeast-2")})
	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}
	// resp has all of the response data, pull out instance IDs:
	for idx := range resp.Reservations {
		for _, ec2instance := range resp.Reservations[idx].Instances {
			inst, ok := m.instances[*ec2instance.InstanceId]
			if !ok {
				inst = &Instance{}
			}

			inst.InstanceType = *ec2instance.InstanceType
			inst.State = *ec2instance.State.Name

			var subnet *Subnet
			if ec2instance.SubnetId == nil {
				fmt.Println("todo: no subnet for", *ec2instance.InstanceId)
			} else {
				//fmt.Println(":", *ec2instance.SubnetId)
				if subnet, ok = m.subnets[*ec2instance.SubnetId]; !ok {
					subnet = &Subnet{}
					subnet.Name = *ec2instance.SubnetId
					subnet.Instances = make([]*Instance, 0)
					m.subnets[*ec2instance.SubnetId] = subnet
				}
				subnet.Instances = append(subnet.Instances, inst)
			}

			if t, ok := typeToCost[*ec2instance.InstanceType]; ok {
				costToDimension := t * 500 / 3
				inst.Scale = vector.Vector3{costToDimension, costToDimension, costToDimension}
			} else {
				fmt.Println("No typeToCost found for", *ec2instance.InstanceType)
				inst.Scale = vector.Vector3{10, 10, 10}
			}

			m.instances[*ec2instance.InstanceId] = inst
		}
	}

}
