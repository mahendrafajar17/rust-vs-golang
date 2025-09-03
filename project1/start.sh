#!/bin/bash

echo "🚀 Starting services..."

# Kill existing processes
pkill -f "target/debug/project1-rust" 2>/dev/null
pkill -f "go run main.go" 2>/dev/null

# Start Rust service
echo "Starting Rust service..."
(cd rust && cargo run > /dev/null 2>&1) &
echo "✅ Rust running at http://127.0.0.1:3002"

# Start Golang service  
echo "Starting Golang service..."
(cd golang && go run main.go > /dev/null 2>&1) &
echo "✅ Golang running at http://127.0.0.1:8000"

echo ""
echo "📊 Dashboard: http://localhost:3001/d/project1-performance"
echo "🔧 Commands:"
echo "  ./status.sh - Check services"
echo "  ./load-test.sh - Run load test"
echo "  ./stop.sh - Stop services"