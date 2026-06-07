package expo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	defaultEndpoint = "https://exp.host/--/expo-push-notification/send"
	batchSize       = 100
)

// sendRequest is the Expo push API request body.
type sendRequest struct {
	Messages []Message `json:"messages"`
}

// sendResponse is the Expo push API response body.
type sendResponse struct {
	Data []Receipt `json:"data"`
}

// Client sends push notifications via the Expo HTTP push API.
// No external dependencies are needed — it speaks plain JSON over HTTPS.
type Client struct {
	endpoint   string
	httpClient *http.Client
}

// NewClient constructs a Client using the default Expo push endpoint.
func NewClient() *Client {
	return &Client{
		endpoint: defaultEndpoint,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// NewClientWithEndpoint constructs a Client with a custom endpoint and HTTP client.
// Intended for tests only.
func NewClientWithEndpoint(endpoint string, httpClient *http.Client) *Client {
	return &Client{endpoint: endpoint, httpClient: httpClient}
}

// Send sends messages in batches of 100 and returns per-message receipts in the
// same order as the input slice. On transient HTTP errors it retries once.
// Receipts with Status == "error" and Details.Error == ErrDeviceNotRegistered
// indicate the token should be removed from storage.
func (c *Client) Send(ctx context.Context, messages []Message) ([]Receipt, error) {
	if len(messages) == 0 {
		return nil, nil
	}

	var all []Receipt
	for i := 0; i < len(messages); i += batchSize {
		batch := messages[i:min(i+batchSize, len(messages))]

		receipts, err := c.sendBatch(ctx, batch)
		if err != nil {
			// Retry once on transient errors.
			receipts, err = c.sendBatch(ctx, batch)
			if err != nil {
				return nil, fmt.Errorf("send batch starting at %d: %w", i, err)
			}
		}
		all = append(all, receipts...)
	}
	return all, nil
}

func (c *Client) sendBatch(ctx context.Context, messages []Message) ([]Receipt, error) {
	body, err := json.Marshal(sendRequest{Messages: messages})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("expo returned status %d", resp.StatusCode)
	}

	var result sendResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return result.Data, nil
}
