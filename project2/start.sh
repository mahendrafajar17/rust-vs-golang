#!/bin/bash

echo "🚀 Starting Project 2 - RabbitMQ Queue Processors..."

# Kill existing processes
pkill -f "target/debug/project2-rust" 2>/dev/null
pkill -f "go run.*main.go" 2>/dev/null

# Check if RabbitMQ is running
if ! pgrep -f rabbitmq-server > /dev/null; then
    echo "⚠️  RabbitMQ not running. Please start it with: brew services start rabbitmq"
    exit 1
fi

# Start Rust queue processor
echo "Starting Rust queue processor..."
(cd rust && cargo run > /dev/null 2>&1) &
echo "✅ Rust processor running at http://127.0.0.1:8083"

# Start Golang queue processor  
echo "Starting Golang queue processor..."
(cd golang && go run cmd/main.go > /dev/null 2>&1) &
echo "✅ Golang processor running at http://127.0.0.1:8082"

echo ""
echo "📊 Metrics endpoints:"
echo "  🦀 Rust: http://127.0.0.1:8083/metrics"
echo "  🐹 Golang: http://127.0.0.1:8082/metrics"
echo "🐰 RabbitMQ Management: http://127.0.0.1:15672 (guest/guest)"
echo ""
echo "🔧 Commands:"
echo "  ./status.sh - Check services"
echo "  ./stop.sh - Stop services"
echo "  ./load-test.sh - Run load test (if available)"