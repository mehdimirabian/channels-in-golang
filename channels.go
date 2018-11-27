// _Closing_ a channel indicates that no more values
// will be sent on it. This can be useful to communicate
// completion to the channel's receivers.

package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
)

const NumberOfMessages = 3

func LogMessagesToFile(fileName string, message string) {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(message)
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// Here's the worker goroutine. It repeatedly receives
// from `jobs` with `j, more := <-jobs`. In this
// special 2-value form of receive, the `more` value
// will be `false` if `jobs` has been `close`d and all
// values in the channel have already been received.
// We use this to notify on `done` when we've worked
// all our jobs.
func Receive(jobs chan int) {
	for {
		select {
		case j, more := <-jobs:
			if more {
				//fmt.Println("received job", j)
				LogMessagesToFile("testLogs", "received job"+strconv.Itoa(j))
			} else {
				fmt.Println("received all jobs")
				return
			}
		}
	}
}

// This sends 3 jobs to the worker over the `jobs`
// channel, then closes it.
func Send(jobs chan int) {
	for j := 1; j <= NumberOfMessages; j++ {
		jobs <- j
		//fmt.Println("sent job", j)
		LogMessagesToFile("testLogs", "sent job"+strconv.Itoa(j))
	}
	close(jobs)
	runtime.GC()
	fmt.Println("sent all jobs")
}

func Serve() {
	jobs := make(chan int)
	go Receive(jobs)
	go Send(jobs)
}

// In this example we'll use a `jobs` channel to
// communicate work to be done from the `main()` goroutine
// to a worker goroutine. When we have no more jobs for
// the worker we'll `close` the `jobs` channel.
func main() {
	for {
		Serve()
		//time.Sleep(2*time.Second)
		PrintMemUsage()
	}
}
