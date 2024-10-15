package main

import (
	"dt-online-kurs/internal/alerts"
	"dt-online-kurs/internal/disruptive"

	//"dt-online-kurs/internal/task1"
	//"dt-online-kurs/internal/task2"
	//"dt-online-kurs/internal/task3"
	"log"
	"log/slog"
	"math"
)

func main() {
	alerts.Location = "Oslo" // E.g. Oslo, New York, Tokyo, etc.
	if alerts.Location == "" {
		panic("Location must be set")
	}

	handler := loggingEventHandler
	//handler = task1.EventHandler
	//handler = task2.EventHandler
	//handler = task3.EventHandler

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
