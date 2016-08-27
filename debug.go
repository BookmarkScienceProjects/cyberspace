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

func printFPS(framesPerSec float64) {
	warningFPS := (1 / framesPerSec) - 1

	go func() {
		timer := 1 * time.Second
		prev := atomic.LoadUint64(&currentFrame)
		prevTime := time.Now()
		for {
			frame := atomic.LoadUint64(&currentFrame)
			currentTime := <-time.After(timer)
			fps := float64(frame-prev) / currentTime.Sub(prevTime).Seconds()
			if fps < warningFPS {
				Printf("fps: %0.1f < %0.1f frame %d\n", fps, warningFPS, frame)
			} else {
				dPrintf("fps: %0.1f frame %d\n", fps, frame)
			}
			prev = frame
			prevTime = currentTime
		}
	}()
}

func Printf(format string, a ...interface{}) {
	log.Printf(format, a...)
}

func Println(a ...interface{}) {
	log.Println(a...)
}

func dPrintf(format string, a ...interface{}) {
	if verbosity < verbosityDebug {
		return
	}
	Printf(format, a...)
}
