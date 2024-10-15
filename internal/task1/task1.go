package task1

import (
	"dt-online-kurs/internal/alerts"
	"dt-online-kurs/internal/disruptive"
	"fmt"
	"time"
)

// A store manager in a large grocery store is tired of people forgetting to
// close the freezer doors after taking out their items.
// The manager have installed proximity sensors in each door and wants to be
// notified whenever a door is open, so they can close it.

func EventHandler(event disruptive.Event) {
	// Create an event handler which triggers an alert if the door is open and resolves it when the door is closed.
}

func ExampleAlert() {
	alertID := "my-door-alert-id-1"

	err := alerts.Trigger(alertID)
	if err != nil {
		panic(err) // unrecoverable error, panic stops the program
	}
	time.Sleep(time.Second)
	err = alerts.Resolve(alertID)
	if err != nil {
		panic(err)
	}
}

func ExampleProximityEvent() {
	var event disruptive.Event

	objectPresent, ok := event.Data.(disruptive.ObjectPresentData)
	if !ok {
		// Unwanted event type
		return
	}
	msg := "OPEN"
	if objectPresent.ObjectPresent {
		msg = "CLOSED"
	}
	fmt.Println(msg)
}
