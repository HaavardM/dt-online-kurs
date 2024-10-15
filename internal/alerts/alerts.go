package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

var Location = ""

const baseUrl = "https://dt-online-kurs-default-rtdb.europe-west1.firebasedatabase.app"

type notification struct {
	City        string `json:"city"`
	Description string `json:"description"`
	ID          string `json:"id,omitempty"`
	Resolved    bool   `json:"resolved"`
}

var (
	mtx   sync.Mutex
	state = make(map[string]bool)
)

func Trigger(alertID string) error {
	mtx.Lock()
	defer mtx.Unlock()
	if state, ok := state[alertID]; ok && state {
		return nil
	}
	state[alertID] = true

	id := Location + "-" + alertID
	b, err := json.Marshal(notification{
		City:     Location,
		Resolved: false,
		ID:       id,
	})
	if err != nil {
		return fmt.Errorf("unable to encode notification: %w", err)
	}

	url, err := url.Parse(baseUrl + fmt.Sprintf("/notifications/%s.json", id))
	if err != nil {
		return fmt.Errorf("unable to parse url: %w", err)
	}

	resp, err := http.DefaultClient.Do(&http.Request{
		URL:    url,
		Method: http.MethodPut,
		Body:   io.NopCloser(bytes.NewReader(b)),
	})
	if err != nil {
		return fmt.Errorf("unable to post notification: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
func Resolve(alertID string) error {
	mtx.Lock()
	defer mtx.Unlock()
	if state, ok := state[alertID]; ok && !state {
		return nil
	}
	state[alertID] = false

	id := Location + "-" + alertID
	url, err := url.Parse(baseUrl + fmt.Sprintf("/notifications/%s.json", id))
	if err != nil {
		return fmt.Errorf("unable to parse url: %w", err)
	}

	resp, err := http.DefaultClient.Do(&http.Request{
		URL:    url,
		Method: http.MethodDelete,
	})
	if err != nil {
		return fmt.Errorf("unable to post notification: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
