# RabbitMQ Queue Listener & Publisher - Golang Implementation

This project implements a RabbitMQ queue processor that listens to messages, adds UUID, and forwards them to another queue.

## Architecture

Based on reference patterns from:
- Publish Queue Reference: `/Users/mahendrafajar/Repository/JatisMobile/waba-integrate/webhook-receiver/webhook-receiver`
- Listen Queue Reference: `/Users/mahendrafajar/Repository/JatisMobile/waba-integrate/moengage-coster-converter/`

## Project Structure

```
├── cmd/
│   └── main.go                 # Application entry point
├── config/
│   └── config.go              # Configuration management
├── provider/
│   ├── amqpx/
│   │   └── connection.go      # AMQP connection with auto-reconnect
│   └── messaging/
│       ├── publisher.go       # Message publisher
│       └── consumer.go        # Message consumer
├── config.yaml                # Configuration file
├── go.mod                     # Go modules
└── README.md
```

## Features

- **Auto-reconnection**: Automatic AMQP connection recovery
- **Graceful shutdown**: Proper cleanup on termination signals
- **Configurable concurrency**: Multiple worker goroutines
- **Structured logging**: JSON formatted logs with request correlation
- **Message persistence**: Durable queues and persistent messages
- **Error handling**: Proper message acknowledgment and requeue on errors

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
  name: "rabbitmq-queue-processor"
  port: 8082

amqp:
  scheme: "amqp"
  host: "localhost"
  port: 5672
  username: "admin"
  password: "password"
  concurrent: 10
  prefetch_count: 50

queues:
  input_queue: "input_queue"
  output_queue: "output_queue"
```

## Running the Application

1. Install dependencies:
```bash
go mod tidy
```

2. Start RabbitMQ server

3. Run the application:
```bash
go run cmd/main.go
```

## Dependencies

- `github.com/rabbitmq/amqp091-go`: RabbitMQ client
- `github.com/sirupsen/logrus`: Structured logging
- `github.com/spf13/viper`: Configuration management  
- `github.com/google/uuid`: UUID generation