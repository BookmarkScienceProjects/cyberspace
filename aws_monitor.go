package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/stojg/cyberspace/lib/components"
	"github.com/stojg/cyberspace/lib/formation"
	//"github.com/stojg/vector"
	. "github.com/stojg/vivere/lib/components"
	//"math/rand"
	//"github.com/stojg/vector"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func NewInstanceList() *InstanceList {
	return &InstanceList{
		instances:  make(map[string]*AWSInstance),
		formations: formation.NewManager(formation.NewDefensiveCircle(8, 0)),
	}
}

type InstanceList struct {
	sync.Mutex
	instances  map[string]*AWSInstance
	formations *formation.Manager
}

func (i *InstanceList) All() map[string]*AWSInstance {
	result := make(map[string]*AWSInstance, len(i.instances))
	i.Lock()
	for k, v := range i.instances {
		result[k] = v
	}
	i.Unlock()
	return result
}

func (i *InstanceList) Fetch() {
	ticker := time.NewTicker(time.Second * 60)
	go func() {
		for {
			Println("fetching instances")
			GetInstances(i)
			<-ticker.C
		}
	}()
}

func (i *InstanceList) SetInstance(t *ec2.Instance) {
	inst, ok := i.instances[*t.InstanceId]
	if !ok {
		eID := entities.Create()
		inst = &AWSInstance{
			ID:        eID,
			Model:     modelList.New(eID, 10, 10, 10, 1),
			RigidBody: rigidList.New(eID, 1),
			State:     "",
		}
		collisionList.New(eID, 10, 10, 10)
		inst.Model.Position().Set(rand.Float64()*2000-1000, inst.Scale[2]/2, rand.Float64()*2000-1000)
		i.Lock()
		i.instances[*t.InstanceId] = inst
		i.Unlock()
		i.formations.AddCharacter(inst)
		i.formations.UpdateSlots()
	}
	inst.Update(t)
	time.Sleep(time.Second)
}

var monitor *awsMonitor

func init() {

	monitor = &awsMonitor{
		instances: make(map[string]*AWSInstance),
	}

	http.HandleFunc("/monitor", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", 405)
			return
		}

		id := r.FormValue("id")
		if id == "" {
			return
		}

		realID, err := strconv.Atoi(id)
		if err != nil {
			Printf("atoi error: %s", err)
		}

		found := monitor.FindByEntityID(Entity(realID))

		t, err := found.MarshalJSON()
		if err != nil {
			Printf("Error during info json marshalling: %s", err)
		}

		if _, err := w.Write(t); err != nil {
			Printf("error for /monitor endpoint %s", err)
		}
	})
}

func GetInstances(list *InstanceList) {

	regions := []*string{
		aws.String("us-east-1"),
		aws.String("us-west-2"),
		aws.String("us-west-1"),
		aws.String("eu-west-1"),
		aws.String("eu-central-1"),
		aws.String("ap-southeast-1"),
		aws.String("ap-northeast-1"),
		aws.String("ap-southeast-2"),
		aws.String("ap-northeast-2"),
		aws.String("ap-south-1"),
		aws.String("sa-east-1"),
	}
	for _, region := range regions {
		session := session.New()
		svc := ec2.New(session, &aws.Config{Region: region})
		resp, err := svc.DescribeInstances(nil)
		if err != nil {
			Println(err)
			continue
		}
		for idx := range resp.Reservations {
			for _, ec2Inst := range resp.Reservations[idx].Instances {
				list.SetInstance(ec2Inst)
			}
		}
		time.Sleep(time.Second)
	}

}

type awsMonitor struct {
	sync.Mutex
	instances map[string]*AWSInstance
}

func (m *awsMonitor) FindByEntityID(id Entity) *AWSInstance {
	m.Lock()
	defer m.Unlock()
	for _, inst := range m.instances {
		if inst == nil {
			continue
		}
		if *inst.ID == (id) {
			return inst
		}
	}
	return nil
}

//func (m *awsMonitor) UpdateInstances(rootNode *TreeNode) {
//
//	regions := []*string{
//		aws.String("us-east-1"),
//		aws.String("us-west-2"),
//		aws.String("us-west-1"),
//		aws.String("eu-west-1"),
//		aws.String("eu-central-1"),
//		aws.String("ap-southeast-1"),
//		aws.String("ap-northeast-1"),
//		aws.String("ap-southeast-2"),
//		aws.String("ap-northeast-2"),
//		aws.String("ap-south-1"),
//		aws.String("sa-east-1"),
//	}
//
//	instanceCount := 0
//
//	for _, region := range regions {
//
//		sess := session.New()
//		svc := ec2.New(sess, &aws.Config{Region: region})
//		resp, err := svc.DescribeInstances(nil)
//		if err != nil {
//			panic(err)
//		}
//
//		// resp has all of the response data, pull out instance IDs:
//		for idx := range resp.Reservations {
//			for _, ec2Inst := range resp.Reservations[idx].Instances {
//				inst, ok := m.instances[*ec2Inst.InstanceId]
//				if !ok {
//					inst = &AWSInstance{
//						ID: entities.Create(),
//					}
//					inst.SetCPUCreditBalance(100)
//				}
//				m.Lock()
//				m.instances[*ec2Inst.InstanceId] = inst
//				inst.Update(ec2Inst)
//				m.Unlock()
//
//				instanceCount++
//
//				model := modelList.Get(inst.ID)
//				if model == nil {
//					model = modelList.New(inst.ID, inst.Scale[0], inst.Scale[1], inst.Scale[2], 1)
//					model.Position().Set(rand.Float64()*2000-1000, inst.Scale[2]/2, rand.Float64()*2000-1000)
//					inst.Position = model.Position()
//				}
//
//				model.Model = 2
//				if inst.State != "running" {
//					model.Model = 0
//				}
//
//				body := rigidList.Get(inst.ID)
//				if body == nil {
//					//invMass := typeToCost["t2.nano"] / typeToCost[inst.InstanceType]
//					invMass := 1.0
//					body = rigidList.New(inst.ID, invMass)
//					body.MaxAcceleration = &vector.Vector3{100, 100, 100}
//					body.SetAwake(true)
//				}
//
//				if collisionList.Get(inst.ID) == nil {
//					collisionList.New(inst.ID, inst.Scale[0], inst.Scale[1], inst.Scale[2])
//				}
//
//				if controllerList.Get(inst.ID) == nil {
//					controllerList.New(inst.ID, newAI(inst.ID, model, body))
//				}
//
//				inst.Tree = rootNode
//				rootNode.Add(inst)
//				m.SetMetrics(inst, region)
//			}
//		}
//	}
//
//	Printf("%d instance found", instanceCount)
//}

func (m *awsMonitor) SetMetrics(inst *AWSInstance, region *string) {

	if inst.State != "running" {
		inst.SetCPUUtilization(0)
		return
	}

	cw := cloudwatch.New(session.New(), &aws.Config{Region: region})

	point, err := m.GetEc2Metric(inst.InstanceID, "CPUUtilization", cw)
	if err != nil {
		Printf("%s", err)
	} else {
		inst.SetCPUUtilization(point)
	}

	time.Sleep(10 * time.Millisecond)

	if inst.HasCredits {
		point, err = m.GetEc2Metric(inst.InstanceID, "CPUCreditBalance", cw)
		if err != nil {
			Printf("%s", err)
		} else {
			inst.SetCPUCreditBalance(point)
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func (m *awsMonitor) GetEc2Metric(instanceID, metricName string, cw *cloudwatch.CloudWatch) (float64, error) {
	endTime := time.Now()
	startTime := endTime.Add(-10 * time.Minute)
	period := int64(3600)
	metrics, err := cw.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EC2"),
		MetricName: aws.String(metricName),
		Dimensions: []*cloudwatch.Dimension{
			&cloudwatch.Dimension{
				Name:  aws.String("InstanceId"),
				Value: aws.String(instanceID),
			},
		},
		StartTime: &startTime,
		EndTime:   &endTime,
		Period:    &period,
		Statistics: []*string{
			aws.String("Average"),
		},
	})

	if err != nil {
		return 0, err
	}

	if len(metrics.Datapoints) > 0 {
		return *metrics.Datapoints[0].Average, nil

	}
	return 0, fmt.Errorf("No datapoints for %s and metric %s", instanceID, metricName)
}
