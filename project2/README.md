# Project 2: RabbitMQ Queue Listener & Publisher - Performance Comparison

This project compares the performance of RabbitMQ queue processors implemented in Rust and Golang. Both implementations listen to messages, add UUIDs, and forward them to another queue with Prometheus metrics collection.

## Reference Architecture

**Publish Queue Reference:**  
`/Users/mahendrafajar/Repository/JatisMobile/waba-integrate/webhook-receiver/webhook-receiver`

**Listen Queue Reference:**  
`/Users/mahendrafajar/Repository/JatisMobile/waba-integrate/moengage-coster-converter/`

**Note:** Golang project follows the structure patterns from these reference codebases.

## Message Flow

1. **Listen Message from Input Queue:**
```json
{
  "user_id": "12345",
  "product_name": "Laptop",
  "quantity": 2,
  "price": 999.99
}
```

2. **Add UUID and Forward to Output Queue:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "12345",  
  "product_name": "Laptop",
  "quantity": 2,
  "price": 999.99
}
```

## Project Structure

```
project2/
├── start.sh                    # Start both services
├── status.sh                   # Check service status  
├── stop.sh                     # Stop all services
├── load-test.sh               # Run load test
├── golang/                     # Golang implementation
│   ├── cmd/main.go            # Main application
│   ├── config/config.go       # Configuration management
│   ├── provider/
│   │   ├── amqpx/             # Auto-reconnecting AMQP wrapper
│   │   ├── messaging/         # Publisher & Consumer logic
│   │   └── metrics/           # Prometheus metrics
│   └── config.yaml            # Configuration file
├── rust/                      # Rust implementation  
│   ├── src/
│   │   ├── main.rs           # Main application
│   │   ├── config/           # Configuration management
│   │   ├── amqp/             # AMQP connection with retry
│   │   ├── messaging/        # Publisher & Consumer logic
│   │   └── metrics/          # Prometheus metrics
│   └── config.yaml           # Configuration file
└── testing/                   # Load testing tools
    ├── publisher.go          # Message publisher for testing
    └── go.mod               # Go module for testing
```

## Key Features

### Both Implementations Include:
- **Auto-reconnection**: Automatic RabbitMQ connection recovery
- **Graceful shutdown**: Proper cleanup on termination signals
- **Configurable concurrency**: Multiple worker processing
- **Structured logging**: JSON formatted logs with request correlation
- **Message persistence**: Durable queues and persistent messages
- **Prometheus metrics**: Comprehensive performance monitoring
- **Simple HTTP endpoint**: `/metrics` for Prometheus scraping

### Golang Specific Features:
- Logrus structured logging
- Viper configuration management
- Custom AMQP wrapper with reconnection logic
- Object pooling for performance optimization

### Rust Specific Features:
- Tokio async runtime for high performance
- Axum for simple HTTP metrics endpoint
- Tracing structured logging and diagnostics
- Lapin high-performance RabbitMQ client
- Semaphore-based concurrency control
- Zero-cost abstractions with compile-time guarantees

## Metrics Collected:
- **Message Metrics**: Received, processed, failed counts and rates
- **Performance Metrics**: Processing latency percentiles
- **System Metrics**: CPU usage, memory usage, goroutines/tasks
- **Queue Metrics**: Queue depth, active consumers
- **Connection Metrics**: AMQP connections, reconnections

## Quick Start

### Automated Setup (Recommended):

```bash
# Start all services (auto-starts RabbitMQ if needed)
./start.sh

# Check service status
./status.sh

# Run load test
./load-test.sh

# Stop all services
./stop.sh
```

### Manual Setup:

1. **Start RabbitMQ:**
```bash
docker run -d --name rabbitmq \
  -p 5672:5672 -p 15672:15672 \
  -e RABBITMQ_DEFAULT_USER=guest \
  -e RABBITMQ_DEFAULT_PASS=guest \
  rabbitmq:3.13-management
```

2. **Run Golang Implementation:**
```bash
cd golang && go run cmd/main.go
```

3. **Run Rust Implementation:**
```bash  
cd rust && cargo run
```

## Monitoring

- **RabbitMQ Management**: http://localhost:15672 (guest/guest)
- **Golang Metrics**: http://localhost:8082/metrics
- **Rust Metrics**: http://localhost:8083/metrics

Connect these to external Prometheus/Grafana for visualization.

## Performance Characteristics

### Golang Implementation:
- **Memory Management**: Garbage collected with object pooling
- **Concurrency**: Goroutines with configurable worker pools  
- **Network I/O**: Efficient with connection reuse
- **Ecosystem**: Rich library ecosystem with established patterns

### Rust Implementation:
- **Memory Safety**: Zero-cost abstractions, no garbage collection
- **Concurrency**: Async/await with efficient task scheduling
- **Performance**: Direct memory management, compile-time optimizations
- **Type Safety**: Compile-time guarantees preventing runtime errors

## Configuration

Both implementations use YAML configuration:

```yaml
app:
  name: "rabbitmq-queue-processor"
  port: 8082  # or 8083 for Rust

amqp:
  url: "amqp://guest:guest@localhost:5672/%2f" 
  concurrent: 10
  prefetch_count: 50

queues:
  input_queue: "input_queue"
  output_queue: "output_queue"
```

## Scripts

- **`start.sh`** - Starts both processors and RabbitMQ (if needed)
- **`status.sh`** - Shows service status and basic metrics
- **`stop.sh`** - Stops all services
- **`load-test.sh`** - Generates test messages for performance testing

This setup provides a simple platform for evaluating RabbitMQ queue processing performance between Rust and Golang implementations with Prometheus metrics collection.