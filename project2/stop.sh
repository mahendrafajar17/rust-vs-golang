#!/bin/bash

echo "🛑 Stopping Project 2 services..."

# Stop Rust processor
if pkill -f "target/debug/project2-rust" 2>/dev/null; then
    echo "✅ Rust queue processor stopped"
else
    echo "⚠️  Rust queue processor not running"
fi

# Stop Golang processor
if pkill -f "go run.*main.go" 2>/dev/null; then
    echo "✅ Golang queue processor stopped"
else
    echo "⚠️  Golang queue processor not running"
fi

# Stop main service (usually running on port 8082)
if pkill -f "./main" 2>/dev/null || pkill -f "main$" 2>/dev/null; then
    echo "✅ Main service stopped"
else
    echo "⚠️  Main service not running"
fi

# Note: RabbitMQ service is managed separately via Homebrew
# Use 'brew services stop rabbitmq' to stop RabbitMQ if needed

echo "🏁 All Project 2 services stopped"