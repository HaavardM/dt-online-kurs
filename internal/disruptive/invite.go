package disruptive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type inviteRequest struct {
	Roles []string `json:"roles"`
	Email string   `json:"email"`
}

func InviteUser(email string) error {
	body := inviteRequest{
		Roles: []string{"roles/project.user"},
		Email: email,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://api.d21s.com/v2/projects/%s/members", projectID)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("Unable to create request: %w", err)
	}

	keyID, secret, err := getKey()
	// Set the "Authorization" header with basic auth.
	// NOTE: OAuth2 authentication should be used in production
	req.SetBasicAuth(keyID, secret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
