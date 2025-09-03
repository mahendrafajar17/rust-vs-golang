#!/bin/bash

echo "🛑 Stopping services..."

# Stop Rust service
if pkill -f "target/debug/project1-rust" 2>/dev/null; then
    echo "✅ Rust service stopped"
else
    echo "⚠️  Rust service not running"
fi

# Stop Golang service
if pkill -f "go run main.go" 2>/dev/null; then
    echo "✅ Golang service stopped"
else
    echo "⚠️  Golang service not running"
fi

echo "🏁 All services stopped"