package disruptive

import (
	"sync"
	"time"
)

// FanOutHandler spawns the function f in a goroutine for each device and routes traffic to them using channels.
func FanOutHandler(f func(string, <-chan Event)) EventHandler {

	// This is a nice place to put shared state between the different goroutines

	var mtx sync.Mutex
	events := make(map[string]chan Event)

	// This new function will be called for each event received by the API.
	return func(event Event) {
		mtx.Lock()
		c, ok := events[event.DeviceID]
		if !ok {
			c = make(chan Event)
			events[event.DeviceID] = c
			go f(event.DeviceID, c)
		}
		mtx.Unlock()
		select {
		case c <- event:
		case <-time.After(10 * time.Second):
			panic("deadlock")
		}
	}
}
