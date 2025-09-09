# RabbitMQ Queue Listener & Publisher - Rust Implementation

This project implements a high-performance RabbitMQ queue processor in Rust that listens to messages, adds UUID, and forwards them to another queue.

## Architecture

Based on reference patterns from:
- Publish Queue Reference: `/Users/mahendrafajar/Repository/JatisMobile/waba-integrate/webhook-receiver/webhook-receiver`  
- Listen Queue Reference: `/Users/mahendrafajar/Repository/JatisMobile/waba-integrate/moengage-coster-converter/`

## Project Structure

```
├── src/
│   ├── main.rs                 # Application entry point
│   ├── config/
│   │   └── mod.rs             # Configuration management
│   ├── amqp/
│   │   ├── mod.rs
│   │   └── connection.rs      # AMQP connection with retry logic
│   └── messaging/
│       ├── mod.rs
│       ├── publisher.rs       # Message publisher
│       └── consumer.rs        # Message consumer
├── config.yaml                # Configuration file
├── Cargo.toml                 # Rust dependencies
└── README.md
```

## Features

- **High Performance**: Async/await with Tokio runtime
- **Connection Resilience**: Automatic connection retry with backoff
- **Graceful Shutdown**: Proper cleanup on SIGTERM/SIGINT
- **Configurable Concurrency**: Async semaphore-based concurrency control
- **Structured Logging**: JSON formatted tracing with request correlation
- **Message Persistence**: Durable queues and persistent messages
- **Error Recovery**: Automatic message requeue on processing failures

## Message Flow

1. **Listen Message from Queue:**
```json
{
  "user_id": "12345",
  "product_name": "Laptop",
  "quantity": 2,
  "price": 999.99
}
```

2. **Add UUID and Forward to Next Queue:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "12345",
  "product_name": "Laptop", 
  "quantity": 2,
  "price": 999.99
}
```

## Configuration

Edit `config.yaml` to configure RabbitMQ connection and queue settings:

```yaml
app:
  name: "rabbitmq-queue-processor-rust"
  port: 8083

amqp:
  url: "amqp://admin:password@localhost:5672/%2f"
  concurrent: 10
  prefetch_count: 50

queues:
  input_queue: "input_queue"
  output_queue: "output_queue"

logging:
  level: "info"
  format: "json"
```

## Running the Application

1. Install Rust (if not already installed):
```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

2. Build the project:
```bash
cargo build --release
```

3. Start RabbitMQ server

4. Run the application:
```bash
cargo run
```

## Dependencies

- `tokio`: Async runtime
- `lapin`: High-performance RabbitMQ client 
- `serde`: Serialization framework
- `tracing`: Structured logging and diagnostics
- `uuid`: UUID generation
- `anyhow`: Error handling
- `config`: Configuration management

## Performance Characteristics

- **Memory Safety**: Zero-cost abstractions with compile-time guarantees
- **Async Performance**: Non-blocking I/O with efficient task scheduling
- **Low Latency**: Direct memory management without garbage collection
- **High Throughput**: Optimized message processing with configurable concurrency