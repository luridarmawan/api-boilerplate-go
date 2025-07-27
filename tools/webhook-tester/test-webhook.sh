#!/bin/bash

# Webhook Tester Script
# This script helps you test the CallWebhook function

echo "=== Webhook Tester ==="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed or not in PATH"
    exit 1
fi

# Function to show usage
show_usage() {
    echo "Usage:"
    echo "  $0 server [port]           - Start webhook receiver server (default port: 8080)"
    echo "  $0 test <webhook-url>      - Test CallWebhook function"
    echo "  $0 quick                   - Start server and run test (uses localhost:8080)"
    echo ""
    echo "Examples:"
    echo "  $0 server 9000                                    # Start server on port 9000"
    echo "  $0 test http://localhost:8080/webhook             # Test with local server"
    echo "  $0 test https://webhook.site/your-unique-id       # Test with external service"
    echo "  $0 quick                                          # Quick test with local server"
}

# Check arguments
if [ $# -eq 0 ]; then
    show_usage
    exit 1
fi

case "$1" in
    "server")
        PORT=${2:-8080}
        echo "Starting webhook receiver server on port $PORT..."
        echo "Webhook URL: http://localhost:$PORT/webhook"
        echo "Press Ctrl+C to stop"
        echo ""
        go run main.go server $PORT
        ;;
    
    "test")
        if [ -z "$2" ]; then
            echo "Error: Please provide webhook URL"
            echo "Example: $0 test http://localhost:8080/webhook"
            exit 1
        fi
        
        WEBHOOK_URL="$2"
        echo "Testing CallWebhook function..."
        echo "Target URL: $WEBHOOK_URL"
        echo ""
        go run main.go test "$WEBHOOK_URL"
        ;;
    
    "quick")
        echo "Quick test mode: Starting server and running test"
        echo ""
        
        # Start server in background
        echo "Starting webhook receiver server on port 8080..."
        go run main.go server 8080 &
        SERVER_PID=$!
        
        # Wait for server to start
        echo "Waiting for server to start..."
        sleep 2
        
        # Run test
        echo "Running webhook test..."
        go run main.go test http://localhost:8080/webhook
        
        # Stop server
        echo ""
        echo "Stopping server..."
        kill $SERVER_PID 2>/dev/null
        wait $SERVER_PID 2>/dev/null
        echo "Test completed!"
        ;;
    
    *)
        echo "Error: Unknown command '$1'"
        echo ""
        show_usage
        exit 1
        ;;
esac