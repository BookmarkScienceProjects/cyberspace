package main

import (
	"math/rand"
	"sync/atomic"
	"time"
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

	level := newLevel()

	hub := initNetwork(level)

	previous := time.Now()
	frameLag := 0.0
	netLag := 0.0

	Println("Running the the game loop")

	// print the FPS when it's below the frameRate
	printFPS(frameRate)

	for {

		// keep track of which frame we are running
		atomic.AddUint64(&currentFrame, 1)

		// calculate a bunch of time values
		now := time.Now()
		elapsed := now.Sub(previous).Seconds()
		previous = now

		frameLag += elapsed
		netLag += elapsed

		level.Update(elapsed)

		// send updates to the network
		if netLag > netRate {
			sendToClients(hub, level)
			netLag -= netRate
		}

		frameLag -= frameRate
		// save some CPU cycles by sleeping for a while
		time.Sleep(time.Duration((frameRate-frameLag)*1000) * time.Millisecond)
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

	//// we have a separate list that contains 'dead' game objects so that
	//// they get flushed to the networks clients
	deadbuf := lvl.drawDead()
	if deadbuf.Len() > 0 {
		if _, err := hub.Write(2, deadbuf.Bytes()); err != nil {
			Printf("%s", err)
		}
	}
}
