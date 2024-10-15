package disruptive

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"time"
)

const projectID = "cs6bd32lh064fjopjmeg"

// Events received from the stream will be in this format.
// Either `Result` or `Error` will be set, the other will be `nil`.
type streamEventSchema struct {
	Result *struct {
		Event rawEvent `json:"event"`
	} `json:"result"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Details []struct {
			Help string `json:"help"`
		} `json:"details"`
	} `json:"error"`
}

const (
	pingInterval = time.Second * 10 // How often we want a ping from the server
	pingLeeway   = time.Second * 2  // Allow for pings to arrive a little late
)

var (
	// This timer will time out if we have not received a ping from the
	// server within the expected interval. We'll use this as a sign that
	// we have lost connection to the server, and then reconnect.
	pingTimer = time.NewTimer(pingInterval + pingLeeway)
)

func getKey() (string, string, error) {

	f, err := os.Open("keys.json")
	if err != nil {
		return "", "", fmt.Errorf("unable to open keys.json in project root: %w", err)
	}

	var keys []struct {
		ID     string `json:"id"`
		Secret string `json:"secret"`
	}
	err = json.NewDecoder(f).Decode(&keys)
	if err != nil {
		return "", "", fmt.Errorf("unable to decode keys.json: %w", err)
	}

	key := keys[rand.IntN(len(keys))] // nolint: gosec
	return key.ID, key.Secret, nil
}

// Assumes that these environment variables are set before running the script.

// ReceiveLiveEvents connects to the stream and call the handler function in a
// goroutine for each event received. This methods block until an error occurs.
func ReceiveLiveEvents(handler func(Event), eventTypes ...EventType) error {

	devices, err := listDevices()
	if err != nil {
		return fmt.Errorf("unable to list devices: %w", err)
	}

	labelsMap := make(map[string]map[string]string, len(devices))
	for _, device := range devices {
		labelsMap[device.Name] = device.Labels
	}

	// Create the request
	url := fmt.Sprintf("https://api.d21s.com/v2/projects/%s/devices:stream", projectID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Unable to create request: %w", err)
	}

	keyID, secret, err := getKey()
	if err != nil {
		return fmt.Errorf("Unable to get key: %w", err)
	}
	// Set the "Authorization" header with basic auth.
	// NOTE: OAuth2 authentication should be used in production
	req.SetBasicAuth(keyID, secret)

	// Add query parameters to the request.
	// "pingInterval" tells the server to send us a periodic pings to let us know the connection is still up.
	// "eventTypes" specifies which event types we want to receive. No specified eventTypes implies all types.
	queryParams := req.URL.Query()
	queryParams.Add("pingInterval", fmt.Sprintf("%.0fs", pingInterval.Seconds()))
	if len(eventTypes) == 0 {
		eventTypes = []EventType{Temperature, ObjectPresent, Touch, NetworkStatus}
	}
	for _, eventType := range eventTypes {
		queryParams.Add("eventTypes", string(eventType))
	}
	req.URL.RawQuery = queryParams.Encode()

	client := &http.Client{}
	resp, err := client.Do(req) // nolint:bodyclose
	if err != nil {
		return fmt.Errorf("Unable to connect to stream: %w", err)
	}
	defer resp.Body.Close()
	fmt.Printf("Connected to the stream with status code %d\n", resp.StatusCode)

	// Scan for events.
	// Events are received as JSON blobs separated by "\n". scanner.Scan() will return
	// true when it has received a new "\n", and return false when the stream disconnects.
	// The event payload is available in scanner.Bytes().
	// Make sure we close the body when we're done with it
	// Scanner to read each line in the stream

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		rawEvent := scanner.Bytes()
		// Unmarshal the stream event
		var streamEvent streamEventSchema
		if err := json.Unmarshal(rawEvent, &streamEvent); err != nil {
			return fmt.Errorf("Unable to unmarshal stream event: %w", err)
		}

		if streamEvent.Error != nil {
			return fmt.Errorf("Received error from stream: %s", streamEvent.Error.Message)
		} else if streamEvent.Result.Event.EventType == "ping" {
			// We got a ping, so we know the stream is still up
			pingTimer.Reset(pingInterval + pingLeeway)
			continue
		}
		go handler(parseRawEvent(streamEvent.Result.Event, labelsMap[streamEvent.Result.Event.TargetName]))
	}
	return nil
}
