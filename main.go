package main

import (
	"github.com/stojg/vivere/lib/client"
	"github.com/stojg/vivere/lib/components"
	"math/rand"
	"sync/atomic"
	"time"
)

const (
	frameRate float64 = 0.016
)

var (
	currentFrame uint64
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	level := newLevel(monitor)

	previous := time.Now()
	lag := 0.0

	Println("Starting the game loop")
	// @todo fix race condition on the global Frame var
	printFPS(frameRate)

	for {

		for _, client := range clients {
			select {
			case msg := <-client.Input():
				handleInput(client, msg)
			default:
			}
		}

		atomic.AddUint64(&currentFrame, 1)
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
		lag -= frameRate

		// save some CPU cycles by sleeping for a while
		time.Sleep(time.Duration((frameRate-lag)*1000) * time.Millisecond)
	}
}

func handleInput(c *client.Client, msg client.ClientCommand) {
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
