package components

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stojg/vector"
	. "github.com/stojg/vivere/lib/components"
	"math"
	"strings"
	"sync"
)

var typeToCost map[string]float64

func init() {
	typeToCost = map[string]float64{
		"t2.nano":    0.01,
		"t2.micro":   0.02,
		"t2.small":   0.04,
		"t2.medium":  0.08,
		"m4.2xlarge": 0.336,
		"m3.large":   0.186,
		"c3.large":   0.132,
		"c4.large":   0.137,
		"t1.micro":   0.02,
		"m1.small":   0.058,
		"m1.medium":  0.117,
	}
}

type AWSInstance struct {
	sync.Mutex
	ID               *Entity
	Cluster          string
	Stack            string
	Environment      string
	InstanceID       string
	HasCredits       bool
	InstanceType     string
	Scale            vector.Vector3
	State            string
	Name             string
	cpuUtilization   float64
	cpuCreditBalance float64
	PrivateIP        string
	PublicIP         string
	Tree             *TreeNode

	// Position is cached here from the model
	Position *vector.Vector3
}

func (inst *AWSInstance) MarshalJSON() ([]byte, error) {
	inst.Lock()
	defer inst.Unlock()
	
	return json.Marshal(map[string]interface{}{
		"Name":             inst.Name,
		"InstanceID":       inst.InstanceID,
		"InstanceType":     inst.InstanceType,
		"HasCredits":       inst.HasCredits,
		"CPUUtilization":   inst.cpuUtilization,
		"CPUCreditBalance": inst.cpuCreditBalance,
		"PrivateIP":        inst.PrivateIP,
		"PublicIP":         inst.PublicIP,
	})
}

func (inst *AWSInstance) SetCPUUtilization(v float64) {
	inst.Lock()
	defer inst.Unlock()
	inst.cpuUtilization = v
}

func (inst *AWSInstance) CPUUtilization() float64 {
	inst.Lock()
	defer inst.Unlock()
	return inst.cpuUtilization
}

func (inst *AWSInstance) SetCPUCreditBalance(v float64) {
	inst.Lock()
	defer inst.Unlock()
	inst.cpuCreditBalance = v
}

func (inst *AWSInstance) CPUCreditBalance() float64 {
	inst.Lock()
	defer inst.Unlock()
	return inst.cpuCreditBalance
}

func (inst *AWSInstance) Health() float64 {

	maxHealth := 1.0

	if inst.State != "running" {
		return maxHealth
	}
	if inst.HasCredits && inst.CPUCreditBalance() < 10 {
		return 0.0
	}
	return maxHealth - inst.CPUUtilization()/100.0
}

func (inst *AWSInstance) Update(ec2Inst *ec2.Instance) {

	inst.InstanceID = *ec2Inst.InstanceId
	inst.InstanceType = *ec2Inst.InstanceType
	inst.State = *ec2Inst.State.Name
	if ec2Inst.PublicIpAddress != nil {
		inst.PublicIP = *ec2Inst.PublicIpAddress
	}
	if ec2Inst.PrivateIpAddress != nil {
		inst.PrivateIP = *ec2Inst.PrivateIpAddress
	}

	if strings.HasPrefix(inst.InstanceType, "t2") {
		inst.HasCredits = true
	}

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
		fmt.Printf("No typeToCost found for '%s'", *ec2Inst.InstanceType)
		inst.Scale = vector.Vector3{10, 10, 10}
	} else {
		costToDimension := t * 10000
		size := math.Pow(costToDimension, 1/3.0)
		inst.Scale = vector.Vector3{size, size, size}
	}
}

func (i *AWSInstance) String() string {
	return fmt.Sprintf("%s %s %s\t%s\t%s", i.Cluster, i.Stack, i.Environment, i.InstanceType, i.InstanceID)
}
