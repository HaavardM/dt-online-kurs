package disruptive

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Device struct {
	Name   string            `json:"name"`
	Type   string            `json:"type"`
	Labels map[string]string `json:"labels"`
}

func listDevices() ([]Device, error) {
	// Create the request
	url := fmt.Sprintf("https://api.d21s.com/v2/projects/%s/devices", projectID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Unable to create request: %w", err)
	}

	keyID, secret, err := getKey()
	if err != nil {
		return nil, fmt.Errorf("unable to get key: %w", err)
	}
	// Set the "Authorization" header with basic auth.
	// NOTE: OAuth2 authentication should be used in production
	req.SetBasicAuth(keyID, secret)

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to send request: %w", err)
	}
	defer r.Body.Close()

	var resp struct {
		Devices []Device `json:"devices"`
	}
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return nil, fmt.Errorf("Unable to decode response: %w", err)
	}
	return resp.Devices, nil
}
