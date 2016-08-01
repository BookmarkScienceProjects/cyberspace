package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/stojg/vivere/lib/components"
	"github.com/stojg/vivere/lib/vector"
	"math"
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
	})
}

type Monitor struct {
	sync.Mutex
	instances map[string]*Instance
	clusters  map[string]*TreeNode
	position  *vector.Vector3
	area      float64
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
		//aws.String("us-east-1"),
		//aws.String("us-west-2"),
		//aws.String("us-west-1"),
		//aws.String("eu-west-1"),
		//aws.String("eu-central-1"),
		//aws.String("ap-southeast-1"),
		//aws.String("ap-northeast-1"),
		aws.String("ap-southeast-2"),
		//aws.String("ap-northeast-2"),
		//aws.String("ap-south-1"),
		//aws.String("sa-east-1"),
	}

	rootNode := NewTree("root", -1)

	for _, region := range regions {

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
					//body.Position.Set(rand.Float64()*800-400, inst.Scale[1]/2, rand.Float64()*800-400)
				}
				body.Model = 2
				//if inst.State != "running" {
				//	body.Model = 2
				//}
				//
				//if rigidList.Get(inst.ID) == nil {
				//	rig := rigidList.New(inst.ID, 1)
				//	rig.MaxAcceleration = &vector.Vector3{10, 10, 10}
				//}
				//
				//if collisionList.Get(inst.ID) == nil {
				//	collisionList.New(inst.ID, 10, 100, 10)
				//}
				//
				//if controllerList.Get(inst.ID) == nil {
				//	controllerList.New(inst.ID, NewAI(inst.ID))
				//}

				rootNode.Add(inst)

				m.SetMetrics(inst, region)
				//time.Sleep(time.Millisecond * 1000)
			}
		}
	}

	Some(rootNode)
	Other(rootNode)

}

func Some(c Collidable) {
	for _, c := range c.Children() {
		Some(c)
	}

	// 1. place all instances (leaf nodes)
	for _, leaf := range c.Leaves() {
		//Println("Placing", leaf.Name())
		leaf.SetPosition(vector.NewVector3(
			(rand.Float64()-0.5)*0,
			//leaf.Instance().Scale[1]/2,
			0,
			(rand.Float64()-0.5)*0,
		))
	}

	//2. resolve all collisions among leaf nodes
	if len(c.Leaves()) > 1 {
		num := 1
		for num != 0 {
			//Println("Resolve leaf nodes", c.Name())
			num = CheckLeaves(c.Leaves())
		}
	}

	for _, leaf := range c.Children() {
		//Println("Placing", leaf.Name())
		leaf.AddPosition(vector.NewVector3(
			(rand.Float64()-0.5)*300,
			//leaf.Instance().Scale[1]/2,
			0,
			(rand.Float64()-0.5)*300,
		))
	}

	if len(c.Children()) > 1 {
		num := 1
		for num != 0 {
			//Println("Checking intersection in group", c.Name())
			num = CheckLeaves(c.Children())
		}
	}
}

func Other(c Collidable) {
	for _, c := range c.Children() {
		Other(c)
	}

	for _, leaf := range c.Leaves() {
		i := leaf.Instance()
		if i != nil {
			body := modelList.Get(i.ID)
			body.Position = leaf.Position()
		}
	}
}

func CheckLeaves(leaves []Collidable) int {

	checked := make(map[int]map[int]bool)
	collisions := make([]*rectCol, 0)

	for i, a := range leaves {
		for j, b := range leaves {
			if i == j {
				continue
			}
			if _, ok := checked[i][j]; ok {
				continue
			}
			if _, ok := checked[j][i]; ok {
				continue
			}
			if _, ok := checked[i]; !ok {
				checked[i] = make(map[int]bool)
			}

			if _, ok := checked[j]; !ok {
				checked[j] = make(map[int]bool)
			}
			checked[i][j], checked[j][i] = true, true

			pair := &rectCol{
				A: a,
				B: b,
			}

			pair.intersects()

			if pair.IsIntersecting {
				collisions = append(collisions, pair)
			}
		}
	}

	num := len(collisions)
	var biggest *rectCol
	pen := -math.MaxFloat64
	for _, pair := range collisions {
		if pair.penetration > pen {
			pen = pair.penetration
			biggest = pair
		}

	}
	if pen > 0 {
		biggest.Resolve()
	}

	return num
}

type rectCol struct {
	A, B           Collidable
	penetration    float64
	normal         *vector.Vector3
	IsIntersecting bool
}

func (contact *rectCol) intersects() {
	mtvDistance := math.MaxFloat32 // Set current minimum distance (max float value so next value is always less)
	mtvAxis := &vector.Vector3{}   // Axis along which to travel with the minimum distance

	// [Axes of potential separation]
	// [X Axis]
	if !testAxisSeparation(vector.UnitX, contact.A.MinPoint(0), contact.A.MaxPoint(0), contact.B.MinPoint(0), contact.B.MaxPoint(0), mtvAxis, &mtvDistance) {
		return
	}

	// [Y Axis]
	if !testAxisSeparation(vector.UnitY, contact.A.MinPoint(1), contact.A.MaxPoint(1), contact.B.MinPoint(1), contact.B.MaxPoint(1), mtvAxis, &mtvDistance) {
		return
	}

	// [Z Axis]
	if !testAxisSeparation(vector.UnitZ, contact.A.MinPoint(2), contact.A.MaxPoint(2), contact.B.MinPoint(2), contact.B.MaxPoint(2), mtvAxis, &mtvDistance) {
		return
	}

	contact.penetration = mtvDistance + 1
	contact.normal = mtvAxis.Normalize()
	contact.IsIntersecting = true
}

func (contact *rectCol) Resolve() {
	if contact.penetration <= 0 {
		return
	}
	movePerIMass := contact.normal.NewScale(contact.penetration / 2)
	contact.A.AddPosition(movePerIMass.NewScale(1))
	contact.B.AddPosition(movePerIMass.NewScale(-1))
}

func (m *Monitor) SetMetrics(inst *Instance, region *string) {

	if inst.State != "running" {
		inst.CPUUtilization = 0
		return
	}

	cw := cloudwatch.New(session.New(), &aws.Config{Region: region})

	point, err := m.GetEc2Metric(inst.InstanceID, "CPUUtilization", cw)
	if err != nil {
		Printf("%s", err)
	} else {
		inst.CPUUtilization = point
	}

	//time.Sleep(10 * time.Millisecond)

	if inst.HasCredits {
		point, err = m.GetEc2Metric(inst.InstanceID, "CPUCreditBalance", cw)
		if err != nil {
			Printf("%s", err)
		} else {
			inst.CPUCreditBalance = point
		}
		//time.Sleep(10 * time.Second)
	}
}

func (m *Monitor) GetEc2Metric(instanceID, metricName string, cw *cloudwatch.CloudWatch) (float64, error) {
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
	return 0, errors.New(fmt.Sprintf("No datapoints for %s and metric %s", instanceID, metricName))
}
