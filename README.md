# Rust vs Golang Performance Comparison

## Overview
Performance comparison between Rust and Golang with monitoring using Prometheus + Grafana

## Reference Architecture
Based on: `/Users/mahendrafajar/Repository/JatisMobile/waba-integrate/webhook-receiver/webhook-receiver`

## Project Tasks

### 1. Basic HTTP Endpoints

#### GET /
**Request:**
```json
GET /
```

**Response:**
```json
{
  "status": "success",
  "message": "OK",
  "timestamp": "2025-09-02T10:30:00Z"
}
```

#### POST /
**Request:**
```json
{
  "count": 5,
  "data": "sample_data"
}
```

**Response:**
```json
{
  "status": "success",
  "result": ["item_1", "item_2", "item_3", "item_4", "item_5"],
  "processed_count": 5,
  "timestamp": "2025-09-02T10:30:00Z"
}
```

### 2. RabbitMQ Queue Listener & Publisher

**1. Listen Message from Queue:**
```json
{
  "user_id": "12345",
  "product_name": "Laptop",
  "quantity": 2,
  "price": 999.99
}
```

**2. Add UUID and Forward to Next Queue:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "12345",
  "product_name": "Laptop",
  "quantity": 2,
  "price": 999.99
}
```

### 3. MongoDB CRUD Operations

#### POST /users (Create)
**Request:**
```json
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "age": 30
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "60f5b2c8d1e4a123456789ab",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "age": 30,
    "created_at": "2025-09-02T10:30:00Z"
  }
}
```

#### GET /users/:id (Read)
**Request:**
```json
GET /users/60f5b2c8d1e4a123456789ab
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "60f5b2c8d1e4a123456789ab",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "age": 30,
    "created_at": "2025-09-02T10:30:00Z",
    "updated_at": "2025-09-02T10:30:00Z"
  }
}
```

#### PUT /users/:id (Update)
**Request:**
```json
{
  "name": "John Smith",
  "age": 31
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "60f5b2c8d1e4a123456789ab",
    "name": "John Smith",
    "email": "john.doe@example.com",
    "age": 31,
    "created_at": "2025-09-02T10:30:00Z",
    "updated_at": "2025-09-02T11:00:00Z"
  }
}
```

#### DELETE /users/:id (Delete)
**Request:**
```json
DELETE /users/60f5b2c8d1e4a123456789ab
```

**Response:**
```json
{
  "status": "success",
  "message": "User deleted successfully",
  "deleted_id": "60f5b2c8d1e4a123456789ab"
}
```

### 4. Request Proxy Service

#### POST /proxy
**Request:**
```json
{
  "target_url": "https://external-api.example.com/data",
  "payload": {
    "user_id": "12345",
    "action": "get_profile"
  }
}
```

**Modified Request to External:**
```json
{
  "user_id": "12345",
  "action": "get_profile",
  "proxy_timestamp": "2025-09-02T10:30:00Z",
  "proxy_id": "proxy-123-abc",
  "source": "internal_proxy"
}
```

**Response:**
```json
{
  "status": "success",
  "external_response": {
    "profile": {
      "id": "12345",
      "username": "johndoe",
      "email": "john@example.com"
    }
  },
  "proxy_metadata": {
    "forwarded_at": "2025-09-02T10:30:00Z",
    "response_time_ms": 150
  }
}
```

### 5. Database + Message Queue

#### POST /orders
**Request:**
```json
{
  "user_id": "12345",
  "product_id": "prod-789",
  "quantity": 2,
  "price": 99.99
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "order_id": "order-abc-123",
    "user_id": "12345",
    "product_id": "prod-789",
    "quantity": 2,
    "price": 99.99,
    "total": 199.98,
    "created_at": "2025-09-02T10:30:00Z"
  },
  "queue_status": "message_sent",
  "cache_status": "authorized"
}
```

**RabbitMQ Message:**
```json
{
  "event_type": "order_created",
  "order_id": "order-abc-123",
  "user_id": "12345",
  "product_id": "prod-789",
  "quantity": 2,
  "total": 199.98,
  "timestamp": "2025-09-02T10:30:00Z"
}
```
