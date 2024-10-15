package task2

import (
	"dt-online-kurs/internal/disruptive"
	"fmt"
	"sync"
)

// The store manager now gets alerts every time someone opens the fridge and
// gets alert fatigue. He now only wants to be notified if a door is left open
// for more than 5 seconds.

func deviceWorker(deviceID string, c <-chan disruptive.Event) {

	for {

		// https://gobyexample.com/select
		// https://gobyexample.com/time
		// https://gobyexample.com/timers
		// https://gobyexample.com/tickers

		select {
		case event := <-c:
			// Handle event
			fmt.Println(event)

			// Add more select cases to handle different events
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

	// Stuck? A similar problem is already solved in task 3
}
