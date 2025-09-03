#!/bin/bash

echo "🚀 Load Test"
echo "==========="

DURATION=${1:-20}  # Default 20 seconds
RPS=${2:-5}        # Default 5 requests per second

# Check services
if ! curl -s http://127.0.0.1:3002/ > /dev/null 2>&1; then
    echo "❌ Rust service not running"
    exit 1
fi

if ! curl -s http://127.0.0.1:8000/ > /dev/null 2>&1; then
    echo "❌ Golang service not running"  
    exit 1
fi

echo "Duration: ${DURATION}s, Target RPS: ${RPS}"
echo "Starting load test..."
echo ""

# Load test function
test_service() {
    local name=$1
    local url=$2
    local emoji=$3
    local count=0
    local success=0
    local end_time=$((SECONDS + DURATION))
    
    while [[ $SECONDS -lt $end_time ]]; do
        if curl -s --max-time 1 "$url/" > /dev/null 2>&1; then
            ((success++))
        fi
        if curl -s --max-time 1 -X POST -H "Content-Type: application/json" -d '{"count":3,"data":"test"}' "$url/" > /dev/null 2>&1; then
            ((success++))
        fi
        count=$((count + 2))
        
        # Progress
        if [[ $((count % 20)) -eq 0 ]]; then
            echo "$emoji $name: $count requests sent"
        fi
        
        sleep 0.2
    done
    
    local actual_rps=$(echo "scale=1; $count / $DURATION" | bc 2>/dev/null || echo "N/A")
    local success_rate=$(echo "scale=1; $success * 100 / $count" | bc 2>/dev/null || echo "N/A")
    
    echo "$emoji $name Results:"
    echo "  Total: $count, Success: $success, RPS: $actual_rps, Success rate: $success_rate%"
}

# Run tests in parallel
test_service "Rust" "http://127.0.0.1:3002" "🦀" &
test_service "Golang" "http://127.0.0.1:8000" "🐹" &

wait

echo ""
echo "✅ Load test completed!"
echo "📊 View results: http://localhost:3001/d/project1-performance"