# PulseCheck Monitor Service

Monitor management and orchestration service for PulseCheck.

## Overview

The Monitor Service manages URL monitoring configurations and coordinates with the Checker Service to perform health checks.

## Features

- 📝 Create and store monitor configurations
- 🔍 Retrieve monitor status and details
- 🤝 Integrates with Checker Service
- 💾 PostgreSQL database storage
- 🔧 Environment-based configuration

## Prerequisites

- Go 1.23+
- PostgreSQL database
- [pulse-check-apis](https://github.com/impruthvi/pulse-check-apis)
- [checkerd](https://github.com/impruthvi/pulse-check-checker) running

## Installation

```bash
# Clone repository
git clone https://github.com/impruthvi/pulse-check-monitor
cd pulse-check-monitor

# Download dependencies
go mod download
```

## Configuration

Create a `.env` file (use `.env.example` as template):

```env
DB_URL=postgresql://username:password@localhost:5432/pulsecheck?sslmode=disable
CHECKER_SERVICE_URL=localhost:50052
OTLP_ENDPOINT=localhost:4317
```

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DB_URL` | PostgreSQL connection string | - | ✅ |
| `CHECKER_SERVICE_URL` | Address of checkerd service | - | ✅ |
| `OTLP_ENDPOINT` | OpenTelemetry collector endpoint | `localhost:4317` | ❌ |

## Running the Service

### Locally

```bash
go run main.go
```

The service will start on port **50051**.

### Using Docker

```bash
# Build image
docker build -t monitord .

# Run container (replace with your actual database credentials)
docker run --rm \
  -e DB_URL="postgresql://your_username:your_password@host.docker.internal:5432/pulsecheck?sslmode=disable" \
  -e CHECKER_SERVICE_URL="host.docker.internal:50052" \
  -p 50051:50051 \
  monitord
```

**Note:** Replace `your_username` and `your_password` with your actual PostgreSQL credentials.

## Database Setup

```sql
-- Create database
CREATE DATABASE pulsecheck;

-- Tables are auto-migrated on startup via GORM
```

## Testing with grpcurl

**Note:** You need the proto file from [pulse-check-apis](https://github.com/impruthvi/pulse-check-apis). Clone it first:

```bash
# Clone the APIs repo (one-time setup)
git clone https://github.com/impruthvi/pulse-check-apis.git
```

### Create Monitor

```bash
grpcurl -plaintext \
  -d '{"url": "https://example.com", "interval_seconds": 60}' \
  -proto=pulse-check-apis/monitor/v1/monitor.proto \
  localhost:50051 \
  monitor.v1.MonitorService/CreateMonitor
```

### Get Monitor

```bash
grpcurl -plaintext \
  -d '{"id": "your-monitor-id"}' \
  -proto=pulse-check-apis/monitor/v1/monitor.proto \
  localhost:50051 \
  monitor.v1.MonitorService/GetMonitor
```

## Observability

### OpenTelemetry Tracing

The service exports distributed traces to an OTLP collector. Traces include:

- **gRPC Server Spans** - Automatic instrumentation of all RPC calls
- **CreateMonitor Operations** - Manual spans with attributes:
  - `monitor.url` - URL being monitored
  - `monitor.interval_seconds` - Check interval
  - `monitor.id` - Created monitor ID
- **gRPC Client Spans** - Outgoing calls to checkerd service
- **Error Recording** - Automatic error capture and tagging

### Running with Jaeger

```bash
# Start Jaeger all-in-one
docker run -d --name jaeger \
  -p 4317:4317 \
  -p 16686:16686 \
  jaegertracing/all-in-one:latest

# Set OTLP endpoint
export OTLP_ENDPOINT=localhost:4317

# Run service
go run main.go

# View traces at http://localhost:16686
```