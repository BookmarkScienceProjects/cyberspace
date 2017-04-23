package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/stojg/cyberspace/lib/core"
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

func printFPS(signal chan os.Signal) {

	ticker := time.NewTicker(time.Second * 2)

	go func() {

		var cpuProfilingFile *os.File
		var memProfilingFile *os.File

		prevFrame := atomic.LoadUint64(&currentFrame)
		prevTime := time.Now()

		for {
			select {
			case currentTime := <-ticker.C:
				frame := atomic.LoadUint64(&currentFrame)
				fps := float64(frame-prevFrame) / currentTime.Sub(prevTime).Seconds()
				if fps < 60 {
					if cpuProfilingFile == nil {
						var err error
						cpuProfilingFile, err = os.Create("cpu_profile.out")
						if err != nil {
							log.Printf("%v\n", err)
							os.Exit(1)
						}
						pprof.StartCPUProfile(cpuProfilingFile)
						Printf("Started cpu profiling")
					}
					Printf("fps: %0.1f frame: %d, objects: %d\n", fps, frame, len(core.List.All()))
				}
				prevFrame = frame
				prevTime = currentTime
			case <-signal:

				fmt.Println("got signal, quitting")

				if cpuProfilingFile != nil {
					fmt.Println("stopped cpu profiling")
					pprof.StopCPUProfile()
					err := cpuProfilingFile.Close()
					if err != nil {
						fmt.Printf("error on closing cpu profile file %v\n", err)
					}
				}
				var err error
				memProfilingFile, err = os.Create("mem_profile.out")
				if err != nil {
					log.Fatal(err)
				}
				pprof.WriteHeapProfile(memProfilingFile)
				err = memProfilingFile.Close()
				if err != nil {
					fmt.Printf("error on closing mem profile file %v\n", err)
				}
				fmt.Println("wrote memory profiling file")
				os.Exit(0)
				return
			}
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
