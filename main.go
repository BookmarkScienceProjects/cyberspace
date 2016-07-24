package main

import (
	"github.com/stojg/vivere/lib/client"
	"github.com/stojg/vivere/lib/components"
	"math/rand"
	"time"
)

const (
	SEC_PER_UPDATE float64 = 0.016
)

type Updatable interface {
	Update(elapsed float64)
}

var (
	Frame uint
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func handleMessage(c *client.Client, msg client.ClientCommand) {
	Printf("received message type: %d seq %d, Actions %d", msg.Type, msg.Sequence, msg.Data)
	switch msg.Type {
	case 2:
		inst := monitor.FindByEntityID(components.Entity(msg.Data))
		Printf("%v", inst)
		//buf := &bytes.Buffer{}
		//binary.Write(buf, binary.LittleEndian, float32(Frame))
		//binaryStream(buf, INST_ENTITY_ID, msg.Data)
		//c.Update(buf, msg.Type)
	default:
		Printf("Can't handle message type: %d ", msg.Type)
	}
}

func main() {

	level := NewLevel(monitor)

	var previous time.Time = time.Now()
	var lag float64 = 0

	Println("Starting the game loop")
	// @todo fix race condition on the global Frame var
	DebugFPS(SEC_PER_UPDATE)

	for {

		for _, client := range clients {
			select {
			case msg := <-client.Input():
				handleMessage(client, msg)
			default:
			}
		}

		Frame += 1
		now := time.Now()
		elapsed := now.Sub(previous).Seconds()
		previous = now
		lag += elapsed

		level.Update(elapsed)

		buf := level.Draw()
		if buf.Len() > 0 {
			for _, client := range clients {
				client.Update(buf, 1)
			}
		}
		lag -= SEC_PER_UPDATE

		// save some CPU cycles by sleeping for a while
		time.Sleep(time.Duration((SEC_PER_UPDATE-lag)*1000) * time.Millisecond)
	}
}
