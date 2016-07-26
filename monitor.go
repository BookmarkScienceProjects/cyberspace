package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/stojg/vivere/lib/components"
	"github.com/stojg/vivere/lib/vector"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var typeToCost map[string]float64

var monitor *Monitor

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

	monitor = &Monitor{
		instances: make(map[string]*Instance),
		clusters:  make(map[string]*Cluster),
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

		t, err := json.Marshal(found)
		if err != nil {
			Printf("Error during info json marshalling: %s", err)
		}

		w.Write(t)

		fmt.Println("we got a request for some info at ", r.URL.Path)
	})
}

type Monitor struct {
	sync.Mutex
	instances map[string]*Instance
	clusters  map[string]*Cluster
}

func (m *Monitor) FindByEntityID(id Entity) *Instance {
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

func (m *Monitor) UpdateInstances() {

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

		Printf("Checking region %s", *region)

		sess := session.New()
		svc := ec2.New(sess, &aws.Config{Region: region})
		resp, err := svc.DescribeInstances(nil)
		if err != nil {
			panic(err)
		}
		// resp has all of the response data, pull out instance IDs:
		for idx := range resp.Reservations {
			for _, ec2Inst := range resp.Reservations[idx].Instances {

				inst, ok := m.instances[*ec2Inst.InstanceId]
				if !ok {
					e := entities.Create()
					inst = &Instance{
						CPUCreditBalance: 1000,
						ID: e,
					}
				}
				m.Lock()
				m.instances[*ec2Inst.InstanceId] = inst
				inst.Update(ec2Inst)
				m.Unlock()

				body := modelList.Get(inst.ID)
				if body == nil {
					body = modelList.New(inst.ID, inst.Scale[0], inst.Scale[1], inst.Scale[2], 1)
					body.Position.Set(rand.Float64()*800-400, rand.Float64()*800-400, rand.Float64()*800-400)
				}
				body.Model = 1
				if inst.State != "running" {
					body.Model = 0
				}

				if rigidList.Get(inst.ID) == nil {
					rig := rigidList.New(inst.ID, 1)
					rig.MaxAcceleration = &vector.Vector3{10, 10, 10}
				}

				if collisionList.Get(inst.ID) == nil {
					collisionList.New(inst.ID, inst.Scale[0], inst.Scale[1], inst.Scale[2])
				}

				if controllerList.Get(inst.ID) == nil {
					controllerList.New(inst.ID, NewAI(inst.ID))
				}

				if inst.Cluster != "" {
					cluster, ok := m.clusters[inst.Cluster]
					if !ok {
						cluster = NewCluster(inst.Cluster)
						m.clusters[inst.Cluster] = cluster
					}
					cluster.Add(inst)
				}

				m.SetMetrics(inst, region)
				time.Sleep(time.Millisecond * 100)
			}
		}
		Printf("Checked region %s", *region)
	}
}

func (m *Monitor) SetMetrics(inst *Instance, region *string) {
	cw := cloudwatch.New(session.New(), &aws.Config{Region: region})

	endTime := time.Now()
	startTime := endTime.Add(-(time.Minute * 5))
	period := int64(360)
	statistics := []*string{
		aws.String("Maximum"),
	}
	metrics, err := cw.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EC2"),
		MetricName: aws.String("CPUUtilization"),
		Dimensions: []*cloudwatch.Dimension{
			&cloudwatch.Dimension{
				Name:  aws.String("InstanceId"),
				Value: aws.String(inst.InstanceID),
			},
		},
		StartTime:  &startTime,
		EndTime:    &endTime,
		Period:     &period,
		Statistics: statistics,
	})

	if err != nil {
		Printf("%s", err)
	}

	if len(metrics.Datapoints) > 0 {
		inst.CPUUtilization = *metrics.Datapoints[0].Maximum
		//Printf("%s CPUUtilization: %f", inst.Name, inst.CPUUtilization)
	}

	metrics, err = cw.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/EC2"),
		MetricName: aws.String("CPUCreditBalance"),
		Dimensions: []*cloudwatch.Dimension{
			&cloudwatch.Dimension{
				Name:  aws.String("InstanceId"),
				Value: aws.String(inst.InstanceID),
			},
		},
		StartTime:  &startTime,
		EndTime:    &endTime,
		Period:     &period,
		Statistics: statistics,
	})

	if err != nil {
		Printf("%s", err)
	}

	if len(metrics.Datapoints) > 0 {
		inst.CPUCreditBalance = *metrics.Datapoints[0].Maximum
		//Printf("%s CPUCreditBalance: %f", inst.Name, inst.CPUCreditBalance)
	}

}

func (m *Monitor) something(x, y, z float64) {
}
