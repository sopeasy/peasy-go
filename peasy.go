package peasy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// A unique identifier for your website. You can find it in your peasy dashboard.
//
// required
var WebsiteID string

// The URL to send events to. Can be used with [proxies](https://peasy.so/docs/proxying-through-cf-workers)
//
// optional
var IngestURL string = "https://api.peasy.so/v1/ingest/"

// Track is for tracking custom events and associating them with a certain profile.
//
// Example:
//
//	peasy.Track("user_signup", "user123", map[string]any{
//		"email": "user@example.com",
//		"plan":  "premium",
//	})
//
// Note: a visitor with that profileID must exist
func Track(event string, profileID string, data map[string]any) error {
	payload := map[string]any{
		"website_id": WebsiteID,
		"name":       event,
		"profile_id": profileID,
		"metadata":   data,
	}
	return send("e", payload)
}

// SetProfile is for setting a profile for a certain profileID.
//
// Example:
//
//	peasy.SetProfile("user123", map[string]any{
//		"email": "john@peasy.so"
//	})
func SetProfile(profileID string, profile map[string]any) error {
	payload := map[string]any{
		"profile_id": profileID,
		"profile":    profile,
		"website_id": WebsiteID,
	}
	return send("p", payload)
}

func send(endpoint string, payload map[string]any) error {
	u, err := url.JoinPath(IngestURL, endpoint)
	if err != nil {
		return fmt.Errorf("peasy: failed to send request: %w", err)
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", u, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Peasy-Server", "peasy-go")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
