package main

import (
	"log"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	verbosity = verbosityNormal
)

const (
	verbosityNormal = 0
	verbosityDebug  = 2
)

func init() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	envVerbosity := os.Getenv("VERBOSITY")
	var err error
	if envVerbosity == "" {
		verbosity = 0
	} else {
		verbosity, err = strconv.Atoi(envVerbosity)
	}
	if err != nil {
		panic(err)
	}
}

func printFPS(frameTime float64) {
	//warningFPS := (1 / frameTime) - 1

	ticker := time.NewTicker(time.Second * 2)

	go func() {
		prevFrame := atomic.LoadUint64(&currentFrame)
		prevTime := time.Now()

		for currentTime := range ticker.C {
			frame := atomic.LoadUint64(&currentFrame)
			fps := float64(frame-prevFrame) / currentTime.Sub(prevTime).Seconds()
			//if fps < warningFPS {
			//	Printf("fps: %0.1f < %0.1f frame %d\n", fps, warningFPS, frame)
			//} else {
			Printf("fps: %0.1f frame: %d\n", fps, frame)
			//}
			prevFrame = frame
			prevTime = currentTime

		}
	}()
}

// Printf is a package proxy of fmt.Printf so that we don't need to import fmt all over the place
func Printf(format string, a ...interface{}) {
	log.Printf(format, a...)
}

// Println is a package proxy of fmt.Println so that we don't need to import fmt all over the place
func Println(a ...interface{}) {
	log.Println(a...)
}

func dPrintf(format string, a ...interface{}) {
	if verbosity < verbosityDebug {
		return
	}
	Printf(format, a...)
}
