package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client handles sending messages to Slack via Incoming Webhook
type Client struct {
	webhookURL string
}

// NewClient creates a new Slack client
func NewClient(webhookURL string) *Client {
	return &Client{
		webhookURL: webhookURL,
	}
}

// payload represents the JSON payload for Slack
type payload struct {
	Text string `json:"text"`
}

// Send sends a text message to Slack
func (c *Client) Send(ctx context.Context, text string) error {
	if c.webhookURL == "" {
		return fmt.Errorf("slack webhook url is empty")
	}

	p := payload{Text: text}
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.webhookURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("slack API returned status code %d", resp.StatusCode)
	}

	return nil
}
