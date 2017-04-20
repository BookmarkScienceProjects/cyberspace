package main

import (
	"math/rand"
	_ "net/http/pprof"
	"sync/atomic"
	"time"

	"github.com/stojg/cyberspace/lib/core"
)

const (
	frameRate float64 = 0.016
	netRate   float64 = 0.016
)

var (
	currentFrame uint64
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	lvl := newLevel()

	for i := 0; i < 4; i++ {
		obj := spawn("monster")
		obj.AddAgent(NewMonsterAgent())
		obj.Transform().Position().Set(rand.Float64()*50-24, 0, rand.Float64()*50-25)
		lvl.State["monster_exists"] = true
	}
	bed := spawnNonCollidable("bed")
	bed.Transform().Position().Set(-15, 0, -15)
	bed = spawnNonCollidable("bed")
	bed.Transform().Position().Set(15, 0, -15)
	bed = spawnNonCollidable("bed")
	bed.Transform().Position().Set(-15, 0, 15)
	bed = spawnNonCollidable("bed")
	bed.Transform().Position().Set(15, 0, 15)

	bed = spawnNonCollidable("charge_pad")
	bed.Transform().Position().Set(0, 0, 0)

	hub := initNetwork(lvl)

	previous := time.Now()
	frameLag := 0.0
	netLag := 0.0

	Println("Running the game loop")

	// print the FPS when it's below the frameRate
	printFPS()

	for {
		// keep track of which frame we are running
		atomic.AddUint64(&currentFrame, 1)

		// calculate a bunch of time values
		now := time.Now()
		elapsed := now.Sub(previous).Seconds()
		previous = now

		frameLag += elapsed
		netLag += elapsed

		UpdateSight(elapsed)
		UpdateAI(elapsed)
		UpdatePhysics(elapsed)
		core.List.BuildQuadTree()
		UpdateCollisions(elapsed)

		if len(core.List.FindWithTag("food")) < 10 {
			obj := spawn("food")
			obj.Transform().Position().Set(rand.Float64()*50-25, 0, rand.Float64()*50-25)
			//rot := vector.NewVector3(rand.Float64()*2-1, 0, rand.Float64()*2-1)
			//obj.Transform().Orientation().RotateByVector(rot)
			//obj.Transform().Orientation().Normalize()
			//obj.AddAgent(NewFoodAgent())
		}

		// send updates to the network
		if netLag > netRate {
			sendToClients(hub, lvl)
			netLag -= netRate
		}

		frameLag -= frameRate
		// save some CPU cycles by sleeping for a while
		time.Sleep(time.Duration((frameRate-frameLag)*1000) * time.Millisecond)
	}
}

func UpdateSight(elapsed float64) {
	for _, obj := range core.List.All() {
		core.List.SenseManager().Add(NewSight(obj, 100))
	}
}

func sendToClients(hub *clientHub, lvl *level) {
	// send normal entity data
	buf := lvl.draw()
	if buf.Len() > 0 {
		if _, err := hub.Write(1, buf.Bytes()); err != nil {
			Printf("%s", err)
		}
	}
	// we have a separate list that contains 'dead' game objects so that they get flushed to the networks clients
	deadbuf := lvl.drawDead()
	if deadbuf.Len() > 0 {
		if _, err := hub.Write(2, deadbuf.Bytes()); err != nil {
			Printf("%s", err)
		}
	}
}
