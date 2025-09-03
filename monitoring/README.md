# Monitoring Stack

This directory contains the monitoring configuration for Rust vs Golang performance comparison.

## Components

- **Prometheus**: Metrics collection and storage
- **Grafana**: Visualization and dashboards

## Quick Start

1. Start all services:
```bash
docker-compose up --build
```

2. Access services:
   - Rust app: http://localhost:3000
   - Golang app: http://localhost:8000
   - Prometheus: http://localhost:9090
   - Grafana: http://localhost:3001 (admin/admin)

## Load Testing

Generate some load to see metrics:

```bash
# Test Rust app
for i in {1..100}; do
  curl -X POST http://localhost:3000/ \
    -H "Content-Type: application/json" \
    -d '{"count": 10, "data": "test"}' &
done

# Test Golang app
for i in {1..100}; do
  curl -X POST http://localhost:8000/ \
    -H "Content-Type: application/json" \
    -d '{"count": 10, "data": "test"}' &
done
```

## Grafana Dashboard

The dashboard includes:
- HTTP Requests per Second comparison
- Response time percentiles (p50, p95, p99)
- Total request counters
- Real-time metrics with 5s refresh

## Prometheus Queries

Example queries to use in Grafana:
- Request rate: `rate(http_get_requests_total[5m])`
- Response time p95: `histogram_quantile(0.95, rate(http_response_time_seconds_bucket[5m]))`
- Total requests: `http_get_requests_total + http_post_requests_total`