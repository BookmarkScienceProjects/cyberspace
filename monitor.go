package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"fmt"
	"github.com/stojg/vivere/lib/vector"
)

var typeToSize map[string]vector.Matrix3

type Monitor struct {
	instances []string
	scale vector.Vector3
}

func (m *Monitor) UpdateInstances() {

	svc := ec2.New(session.New(), &aws.Config{Region: aws.String("ap-southeast-2")})
	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}
	// resp has all of the response data, pull out instance IDs:
	for idx := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			fmt.Println(" - Instance ID: ", *inst.InstanceId)
			m.instances = append(m.instances, *inst.InstanceId)
		}
	}

}
