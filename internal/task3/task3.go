package task3

import (
	"dt-online-kurs/internal/disruptive"
	"fmt"
	"math"
	"sync"
	"time"
)

// We want to build an even better solution for the store manager. We know that
// the store manager really wants to be alerted if the food is about to go bad
// because open freezers, he cares little about the doors themselves.
// Normally we would add temperature sensors to solve this problem directly, but
// we want to try to re-use the sensors we already have installed.

// We decide to predict the temperature inside the freezer using Newtons Law of
// Cooling. When the door is open, the temperature will increase, and when the
// door is closed it will decrease. While it will not be a perfect prediction,
// it is good enough for the store manager.

const (
	outsideTemp = float64(20)
	freezerTemp = float64(-20)
	k           = float64(0.01)
)

// Simulates the temperature change in the fridge or freezer using the Newton's Law of Cooling.
// temp0 is the previous temperature and the functions returns the new
// temperature if the door is kept open or closed for a duration.
func newtonsLawOfCooling(temp0 float64, doorClosed bool, duration time.Duration) float64 {

	u := outsideTemp
	if doorClosed {
		u = freezerTemp
	}
	// Newton's Law of Cooling - differential equation
	return u + (temp0-u)*math.Exp(-k*duration.Seconds())
}

func deviceWorker(deviceID string, c <-chan disruptive.Event) {

	for {
		select {
		case event := <-c:
			fmt.Println(event)
		}
	}
}

var (
	mtx    sync.Mutex
	events = make(map[string]chan disruptive.Event)
)

func EventHandler(event disruptive.Event) {
	mtx.Lock()
	c, ok := events[event.DeviceID]
	if !ok {
		c = make(chan disruptive.Event)
		events[event.DeviceID] = c
		go deviceWorker(event.DeviceID, c)
	}
	mtx.Unlock()

	c <- event
}
