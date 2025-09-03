# 🚀 Scripts Documentation

Simple scripts for Rust vs Golang performance comparison.

## Available Scripts

### `./start.sh`
Start both services locally.
```bash
./start.sh
```
- 🦀 Rust: http://127.0.0.1:3002
- 🐹 Golang: http://127.0.0.1:8000

### `./status.sh`
Check service status.
```bash
./status.sh
```

### `./load-test.sh [duration] [rps]`
Run load test on both services.
```bash
# Default: 20 seconds, 5 RPS
./load-test.sh

# Custom: 30 seconds, 10 RPS  
./load-test.sh 30 10
```

### `./stop.sh`
Stop all services.
```bash
./stop.sh
```

## Quick Setup

```bash
# 1. Start monitoring (from root directory)
cd .. && docker compose up -d

# 2. Start services (from project1 directory)
cd project1 && ./start.sh

# 3. Run load test
./load-test.sh

# 4. View dashboard
# http://localhost:3001/d/project1-performance
```

## Dashboard

**Grafana:** http://localhost:3001/d/project1-performance
- Login: admin/admin
- Metrics: TPS, CPU, Memory, Response Time

**Prometheus:** http://localhost:9090

## Ports

- 3002: Rust service
- 8000: Golang service  
- 9090: Prometheus
- 3001: Grafana