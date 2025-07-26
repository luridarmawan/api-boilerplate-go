package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// CallWebhook calls the webhook URL with the response data
// This is a copy of the function from internal/utils/net.go for testing purposes
func CallWebhook(webhookURL string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error marshaling data: %v\n", err)
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
			fmt.Printf("Error creating request: %v\n", err)
			return
		}

		// Set appropriate headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Webhook-Client/1.0")

		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending webhook: %v\n", err)
			return
		}
		defer resp.Body.Close()

		fmt.Printf("Webhook sent successfully, status: %d\n", resp.StatusCode)
	}()
}

// TestData represents sample data to send to webhook
type TestData struct {
	ID        int       `json:"id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

// WebhookReceiver handles incoming webhook requests
func webhookReceiver(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("\n=== Webhook Received ===\n")
	fmt.Printf("Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("Headers:\n")
	for key, values := range r.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}
	fmt.Printf("Body: %s\n", string(payload))
	fmt.Printf("=====================\n\n")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "received"}`))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  go run main.go server [port]     - Start webhook receiver server")
		fmt.Println("  go run main.go test <webhook-url> - Test CallWebhook function")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "server":
		port := "8080"
		if len(os.Args) > 2 {
			port = os.Args[2]
		}

		fmt.Printf("Starting webhook receiver server on port %s\n", port)
		fmt.Printf("Webhook URL: http://localhost:%s/webhook\n", port)
		fmt.Println("Press Ctrl+C to stop")

		http.HandleFunc("/webhook", webhookReceiver)
		log.Fatal(http.ListenAndServe(":"+port, nil))

	case "test":
		if len(os.Args) < 3 {
			fmt.Println("Please provide webhook URL")
			fmt.Println("Example: go run main.go test http://localhost:8080/webhook")
			os.Exit(1)
		}

		webhookURL := os.Args[2]
		fmt.Printf("Testing CallWebhook function with URL: %s\n", webhookURL)

		// Test with different types of data
		testCases := []struct {
			name string
			data interface{}
		}{
			{
				name: "Simple Test Data",
				data: TestData{
					ID:        1,
					Message:   "Hello from webhook tester!",
					Timestamp: time.Now(),
					Status:    "success",
				},
			},
			{
				name: "Map Data",
				data: map[string]interface{}{
					"event":   "test_event",
					"user_id": 12345,
					"data": map[string]string{
						"action": "webhook_test",
						"source": "tester_tool",
					},
				},
			},
			{
				name: "Array Data",
				data: []map[string]interface{}{
					{"id": 1, "name": "Item 1"},
					{"id": 2, "name": "Item 2"},
					{"id": 3, "name": "Item 3"},
				},
			},
		}

		for i, testCase := range testCases {
			fmt.Printf("\n--- Test %d: %s ---\n", i+1, testCase.name)
			
			// Print what we're sending
			jsonData, _ := json.MarshalIndent(testCase.data, "", "  ")
			fmt.Printf("Sending data:\n%s\n", string(jsonData))

			// Call the webhook function
			CallWebhook(webhookURL, testCase.data)
			
			fmt.Printf("Webhook called (async)\n")
			
			// Wait a bit between tests
			time.Sleep(1 * time.Second)
		}

		fmt.Println("\nAll tests completed!")
		fmt.Println("Note: CallWebhook runs asynchronously, so check your webhook receiver for results.")

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}