package task2

import (
	"dt-online-kurs/internal/disruptive"
	"fmt"
	"sync"
)

/*
	The sensors are used to monitor doors on fridges in a grocery store. The
	store manager wants to be alerted when a fridge door is open for more than
	10 seconds.

*/

func deviceWorker(deviceID string, c <-chan disruptive.Event) {

	for {

		// https://gobyexample.com/select
		select {
		case event := <-c:
			// Handle event
			fmt.Println(event)
			// Switch on channels
			// https://gobyexample.com/tickers
		}
	}

}

var (
	mtx    sync.Mutex
	events = make(map[string]chan disruptive.Event)
)

func EventHandler(event disruptive.Event) {

	mtx.Lock()
	// Critical region, only one goroutine can access at a time.
	// Create device-specific channel and spawn a deviceWorker goroutine if not already exists.
	// https://gobyexample.com/maps
	// https://gobyexample.com/goroutines
	// https://gobyexample.com/channels

	mtx.Unlock()

	// Send event to a device-specific channel

	// Stuck? An implementation is found in the task 3
}
