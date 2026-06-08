// Package expo provides a client for the Expo push notification API.
package expo

// Message is a single Expo push notification payload.
type Message struct {
	To    string            `json:"to"`
	Title string            `json:"title"`
	Body  string            `json:"body"`
	Data  map[string]string `json:"data,omitempty"`
}

// ReceiptDetails holds the error details returned by the Expo push API.
type ReceiptDetails struct {
	Error string `json:"error,omitempty"` // "DeviceNotRegistered", etc.
}

// Receipt is the per-message result returned by the Expo push API.
type Receipt struct {
	Status  string          `json:"status"` // "ok" or "error"
	Message string          `json:"message,omitempty"`
	Details *ReceiptDetails `json:"details,omitempty"`
}

// ErrDeviceNotRegistered is the Expo error code for a stale push token.
const ErrDeviceNotRegistered = "DeviceNotRegistered"
