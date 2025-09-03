#!/bin/bash

echo "📋 Service Status"
echo "================="

# Check Rust
if pgrep -f "target/debug/project1-rust" > /dev/null; then
    if curl -s http://127.0.0.1:3002/ > /dev/null 2>&1; then
        echo "🦀 Rust: ✅ RUNNING (http://127.0.0.1:3002)"
    else
        echo "🦀 Rust: ⚠️  PROCESS RUNNING but not responding"
    fi
else
    echo "🦀 Rust: ❌ NOT RUNNING"
fi

# Check Golang
if pgrep -f "go run main.go" > /dev/null; then
    if curl -s http://127.0.0.1:8000/ > /dev/null 2>&1; then
        echo "🐹 Golang: ✅ RUNNING (http://127.0.0.1:8000)"
    else
        echo "🐹 Golang: ⚠️  PROCESS RUNNING but not responding"
    fi
else
    echo "🐹 Golang: ❌ NOT RUNNING"
fi

echo ""
echo "📊 Dashboard: http://localhost:3001/d/project1-performance"