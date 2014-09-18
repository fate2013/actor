/*
Feed events to proxyd which it will proxy to dragond.
*/
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var (
	mode string
)

func init() {
	flag.StringVar(&mode, "mode", "standard", "feed mode")
	flag.Parse()
}

func main() {
	switch mode {
	case "standard":
		runStandardMode()

	case "recall":
		runRecall()

	case "speedup":
		runSpeedup()
	}

}

func runStandardMode() {
	marchId := 0
	events := []int{
		12,  // arrive
		15,  // gather done
		20,  // back home
		501, // speedup
		502, // recall
	}
	for {
		marchId++
		uid := rand.Intn(10000) + 1
		event := events[rand.Intn(len(events))]
		at := time.Now().Add(time.Duration(rand.Intn(1000)) * time.Second)
		sendRequest(uid, marchId, event, at)
		time.Sleep(1 * time.Millisecond)
	}
}

func runRecall() {
	sendRequest(1, 1, 12, time.Now().Add(5*time.Second))
	sendRequest(1, 1, 20, time.Now().Add(8*time.Second))
}

func runSpeedup() {
	sendRequest(1, 1, 12, time.Now().Add(5*time.Second))
	sendRequest(1, 1, 12, time.Now().Add(3*time.Second))
}

func sendRequest(uid, marchId, event int, at time.Time) {
	fmt.Fprintf(os.Stdout, "%d,%d,%d,%d,1,%d,%d,xxx", 1,
		at.Unix(), event, uid, marchId, time.Now().UnixNano())

	fmt.Fprintln(os.Stdout)
}
