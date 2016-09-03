package main

import (
	"math/rand"
	"sync/atomic"
	"time"
)

const (
	frameRate float64 = 0.016
	netRate   float64 = 0.032
)

var (
	currentFrame uint64
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	level := newLevel()

	hub := initNetwork()

	previous := time.Now()
	frameLag := 0.0
	netLag := 0.0

	Println("Starting the game loop")
	// @todo fix race condition on the global Frame var
	printFPS(frameRate)

	for {

		atomic.AddUint64(&currentFrame, 1)
		now := time.Now()
		elapsed := now.Sub(previous).Seconds()
		previous = now
		frameLag += elapsed
		netLag += elapsed

		level.Update(elapsed)

		if netLag > netRate {
			buf := level.Draw()
			if buf.Len() > 0 {
				if _, err := hub.Write(buf.Bytes()); err != nil {
					Printf("%s", err)

				}
			}
			netLag -= netRate
		}

		frameLag -= frameRate
		// save some CPU cycles by sleeping for a while
		time.Sleep(time.Duration((frameRate-frameLag)*1000) * time.Millisecond)
	}
}

//func handleInput(c *client.Client, msg client.ClientCommand) {
//	Printf("received message type: %d seq %d, Actions %d", msg.Type, msg.Sequence, msg.Data)
//	switch msg.Type {
//	case 2:
//		inst := monitor.FindByEntityID(components.Entity(msg.Data))
//		Printf("%v", inst)
//	//buf := &bytes.Buffer{}
//	//binary.Write(buf, binary.LittleEndian, float32(Frame))
//	//binaryStream(buf, INST_ENTITY_ID, msg.Data)
//	//c.Update(buf, msg.Type)
//	default:
//		Printf("Can't handle message type: %d ", msg.Type)
//	}
//}
