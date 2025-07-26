package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// CallWebhook calls the webhook URL with the response data
func CallWebhook(webhookURL string, data interface{}) {
	// This is a simple webhook implementation
	// In production, you might want to add retry logic, authentication, etc.
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	// Make HTTP POST request to webhook URL asynchronously
	go func() {
		// Create HTTP client with timeout
		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		// Create request with JSON data
		req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return
		}

		// Set appropriate headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "CARIK.id-Client/1.0")

		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		// Optional: You can add logging here if needed
		// log.Printf("Webhook sent to %s, status: %d", webhookURL, resp.StatusCode)
	}()
}
