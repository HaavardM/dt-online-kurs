package task3

import (
	"dt-online-kurs/internal/alerts"
	"dt-online-kurs/internal/disruptive"
	"fmt"
	"math"
	"sync"
	"time"
)

/*
	The sensors are used to monitor doors on fridges in a grocery store. The
	store manager wants to be alerted when a fridge door is open for more than
	10 seconds.
*/

const (
	outsideTemp = float64(20)
	freezerTemp = float64(-20)
	k           = float64(0.01)
)

// Simulates the temperature change in the fridge or freezer using the Newton's Law of Cooling.
// temp0 is the previous temperature and the functions returns the new
// temperature if the door is kept open or closed for a duration.
func newtonsLawOfCooling(temp0 float64, closed bool, duration time.Duration) float64 {

	u := outsideTemp
	if closed {
		u = freezerTemp
	}

	return u + (temp0-u)*math.Exp(-k*duration.Seconds())
}

func deviceWorker(deviceID string, c <-chan disruptive.Event) {

	currentTemp := freezerTemp
	isClosed := true

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case event := <-c:
			objectPresent, ok := event.Data.(disruptive.ObjectPresentData)
			if !ok {
				continue
			}

			// The state changed!
			if objectPresent.ObjectPresent != isClosed {
				isClosed = objectPresent.ObjectPresent
			}

		case <-ticker.C:
			currentTemp = newtonsLawOfCooling(currentTemp, isClosed, time.Second)
			fmt.Printf("Current temperature for %s: %.1f C\n", deviceID, currentTemp)
		}

		if !isClosed && currentTemp > 0 {
			err := alerts.Trigger("door-open-" + deviceID)
			if err != nil {
				panic(err)
			}
		} else {
			err := alerts.Resolve("door-open-" + deviceID)
			if err != nil {
				panic(err)
			}
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
