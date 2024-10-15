package main

import (
	"log"
	"log/slog"
	"math"

	"github.com/HaavardM/dt-online-kurs/internal/alerts"
	"github.com/HaavardM/dt-online-kurs/internal/disruptive"
	//"github.com/HaavardM/dt-online-kurs/internal/task1"
	//"github.com/HaavardM/dt-online-kurs/internal/task2"
	//"github.com/HaavardM/dt-online-kurs/internal/task3"
)

func main() {
	alerts.Location = "" // E.g. Oslo, New York, Tokyo, etc.
	if alerts.Location == "" {
		panic("Location must be set")
	}

	handler := func(event disruptive.Event) {

		// It may be easier to debug by looking at a single device during development.
		// TODO: Remove once you are ready to handle all devices.
		if event.DeviceID != "my-device-id" {
			return
		}

		loggingEventHandler(event)
		//task1.EventHandler(event)
		//task2.EventHandler(event)
		//task3.EventHandler(event)

	}

	// If replay-speed flag is set, replay events from historical data.
	// Otherwise, receive live events.
	err := disruptive.ReceiveLiveEvents(handler)
	if err != nil {
		log.Panic(err)
	}

}

func loggingEventHandler(event disruptive.Event) {
	l := slog.Default().With(
		slog.String("name", event.DeviceLabels["name"]),
		slog.Time(slog.TimeKey, event.Timestamp),
		slog.String("device_id", event.DeviceID),
	)
	switch data := event.Data.(type) {
	case disruptive.TemperatureData:
		l.Info("temperature event received", slog.Float64("celsius", math.Round(float64(data.Celsius))))
	case disruptive.ObjectPresentData:
		text := "OPEN"
		if data.ObjectPresent {
			text = "CLOSED"
		}
		l.Info(text+" event received", "object_present", slog.BoolValue(data.ObjectPresent))
	case disruptive.TouchData:
		l.Info("touch event received")
	default:
	}
}
