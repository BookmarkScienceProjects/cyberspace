package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/stojg/cyberspace/lib/components"
	"github.com/stojg/cyberspace/lib/formation"
	"github.com/stojg/vivere/lib/components"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func init() {
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
			Printf("GET /monitor atoi error: %s", err)
		}
		eID := components.Entity(realID)
		//found, ok := instanceList.Get(eID)
		var inst *AWSInstance
		for i, obj := range instanceList.All() {
			if eID == *i {
				inst = obj
				break
			}
		}
		if inst.InstanceID == "" {
			Printf("GET /monitor Could not find instance id %d\n", eID)
			return
		}
		t, err := inst.MarshalJSON()
		if err != nil {
			Printf("GET /monitor : %s", err)
		}
		if _, err := w.Write(t); err != nil {
			Printf("GET /monitor  %s", err)
		}
	})
}

func NewInstanceList() *InstanceList {
	return &InstanceList{
		instances:  make(map[*components.Entity]*AWSInstance),
		formations: formation.NewManager(formation.NewDefensiveCircle(8, 0)),
	}
}

type InstanceList struct {
	sync.Mutex
	instances  map[*components.Entity]*AWSInstance
	formations *formation.Manager
}

func (i *InstanceList) All() map[*components.Entity]*AWSInstance {
	result := make(map[*components.Entity]*AWSInstance, len(i.instances))
	i.Lock()
	for k, v := range i.instances {
		result[k] = v
	}
	i.Unlock()
	return result
}

func (i *InstanceList) Get(eID *components.Entity) (*AWSInstance, bool) {
	i.Lock()
	inst, ok := i.instances[eID]
	i.Unlock()
	return inst, ok
}

func (i *InstanceList) Fetch() {
	ticker := time.NewTicker(time.Second * 60)
	go func() {
		for {
			Println("fetching instances")
			fetchInstances(i)
			<-ticker.C
		}
	}()
}

func (i *InstanceList) SetInstance(t *ec2.Instance) *AWSInstance {
	var inst *AWSInstance
	for _, curr := range i.instances {
		if curr.InstanceID == *t.InstanceId {
			inst = curr
			break
		}
	}
	if inst == nil {
		eID := entities.Create()
		modelID := 0
		if *t.State.Name == "running" {
			modelID = 2
		}
		inst = &AWSInstance{
			ID:        eID,
			Model:     modelList.New(eID, 10, 10, 10, components.EntityType(modelID)),
			RigidBody: rigidList.New(eID, 1),
			Collision: collisionList.New(eID, 10, 10, 10),
			State:     "",
		}
		inst.Model.Position().Set(rand.Float64()*100-50, inst.Scale[2]/2, rand.Float64()*100-50)
		i.Lock()
		i.instances[eID] = inst
		i.Unlock()
		i.formations.AddCharacter(inst)
		i.formations.UpdateSlots()
	}
	inst.Update(t)
	return inst
}

func fetchInstances(list *InstanceList) {

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
				instance := list.SetInstance(ec2Inst)
				fetchMetrics(instance, region)
				time.Sleep(1000 * time.Millisecond)
			}
		}
		time.Sleep(time.Second)
	}
}

func fetchMetrics(inst *AWSInstance, region *string) {

	if inst.State != "running" {
		inst.SetCPUUtilization(0)
		return
	}

	cw := cloudwatch.New(session.New(), &aws.Config{Region: region})

	point, err := getEc2Metric(inst.InstanceID, "CPUUtilization", cw)
	if err != nil {
		Printf("%s", err)
	} else {
		inst.SetCPUUtilization(point)
	}

	if inst.HasCredits {
		point, err = getEc2Metric(inst.InstanceID, "CPUCreditBalance", cw)
		if err != nil {
			Printf("%s", err)
		} else {
			inst.SetCPUCreditBalance(point)
		}
	}
}

func getEc2Metric(instanceID, metricName string, cw *cloudwatch.CloudWatch) (float64, error) {
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
