#!/bin/bash

echo "📋 Project 2 - Service Status"
echo "=============================="

# Check RabbitMQ
if pgrep -f rabbitmq-server > /dev/null; then
    if curl -s http://127.0.0.1:15672/ > /dev/null 2>&1; then
        echo "🐰 RabbitMQ: ✅ RUNNING (http://127.0.0.1:15672)"
    else
        echo "🐰 RabbitMQ: ⚠️  PROCESS RUNNING but management interface not responding"
    fi
else
    echo "🐰 RabbitMQ: ❌ NOT RUNNING"
fi

# Check Rust processor
if pgrep -f "target/debug/project2-rust" > /dev/null; then
    if curl -s http://127.0.0.1:8083/metrics > /dev/null 2>&1; then
        echo "🦀 Rust Processor: ✅ RUNNING (http://127.0.0.1:8083)"
        # Get some basic metrics
        messages_processed=$(curl -s http://127.0.0.1:8083/metrics | grep "^rabbitmq_messages_processed_total" | awk '{print $2}' 2>/dev/null || echo "0")
        echo "   📈 Messages processed: $messages_processed"
    else
        echo "🦀 Rust Processor: ⚠️  PROCESS RUNNING but not responding"
    fi
else
    echo "🦀 Rust Processor: ❌ NOT RUNNING"
fi

# Check Golang processor
if pgrep -f "go run.*main.go" > /dev/null; then
    if curl -s http://127.0.0.1:8082/metrics > /dev/null 2>&1; then
        echo "🐹 Golang Processor: ✅ RUNNING (http://127.0.0.1:8082)"
        # Get some basic metrics
        messages_processed=$(curl -s http://127.0.0.1:8082/metrics | grep "^rabbitmq_messages_processed_total" | awk '{print $2}' 2>/dev/null || echo "0")
        echo "   📈 Messages processed: $messages_processed"
    else
        echo "🐹 Golang Processor: ⚠️  PROCESS RUNNING but not responding"
    fi
else
    echo "🐹 Golang Processor: ❌ NOT RUNNING"
fi

echo ""
echo "📊 Metrics URLs:"
echo "  🦀 Rust: http://127.0.0.1:8083/metrics"
echo "  🐹 Golang: http://127.0.0.1:8082/metrics"
echo "🐰 RabbitMQ Management: http://127.0.0.1:15672"