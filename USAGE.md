# Usage Guide

## Quick Start

### 1. Run the Example

The simplest way to see go-pgtune in action:

```bash
go run ./cmd/example
```

Expected output:
```
# DB Version: 18
# OS Type: linux
# DB Type: web
# Total Memory (RAM): 2 GB
# CPUs num: 1
# Connections num: 100
# Data Storage: ssd

max_connections = 100
shared_buffers = 512MB
effective_cache_size = 1536MB
maintenance_work_mem = 128MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
effective_io_concurrency = 200
work_mem = 4854kB
huge_pages = off
min_wal_size = 1GB
max_wal_size = 4GB
```

### 2. Start the API Server

```bash
go run ./cmd/server
```

Or using Make:
```bash
make run
```

The server will start on port 8080.

### 3. Test the API

```bash
curl -X POST http://localhost:8080/tune \
  -H "Content-Type: application/json" \
  -d '{
    "db_version": 18,
    "os_type": "linux",
    "db_type": "web",
    "total_memory_gb": 2,
    "cpu_num": 1,
    "max_connections": 100,
    "storage_type": "ssd"
  }'
```

## Building

### Build Server Binary

```bash
make build
```

This creates `bin/pgtune-server` binary.

### Run Binary

```bash
./bin/pgtune-server
```

## Docker Usage

### Build Docker Image

```bash
docker build -t pgtune-server .
```

### Run Container

```bash
docker run -p 8080:8080 pgtune-server
```

### Test Container

```bash
curl -X POST http://localhost:8080/tune \
  -H "Content-Type: application/json" \
  -d '{
    "db_version": 18,
    "os_type": "linux",
    "db_type": "web",
    "total_memory_gb": 2,
    "cpu_num": 1,
    "max_connections": 100,
    "storage_type": "ssd"
  }'
```

## Library Usage Examples

### Basic Web Application

```go
package main

import (
    "fmt"
    "github.com/souravbiswassanto/go-pgtune/pkg/pgtune"
)

func main() {
    input := pgtune.TuneInput{
        DBVersion:      18,
        OSType:         pgtune.OSLinux,
        DBType:         pgtune.DBTypeWeb,
        TotalMemory:    4 * pgtune.GB,
        CPUNum:         2,
        MaxConnections: 200,
        StorageType:    pgtune.StorageTypeSSD,
    }

    output, _ := pgtune.Tune(input)
    fmt.Printf("shared_buffers = %s\n", output.SharedBuffers)
}
```

### Data Warehouse Configuration

```go
input := pgtune.TuneInput{
    DBVersion:      18,
    OSType:         pgtune.OSLinux,
    DBType:         pgtune.DBTypeDW,
    TotalMemory:    64 * pgtune.GB,
    CPUNum:         16,
    MaxConnections: 0, // Use default (40 for DW)
    StorageType:    pgtune.StorageTypeSSD,
}

output, _ := pgtune.Tune(input)
```

### OLTP Database

```go
input := pgtune.TuneInput{
    DBVersion:      18,
    OSType:         pgtune.OSLinux,
    DBType:         pgtune.DBTypeOLTP,
    TotalMemory:    32 * pgtune.GB,
    CPUNum:         8,
    MaxConnections: 300,
    StorageType:    pgtune.StorageTypeSSD,
}

output, _ := pgtune.Tune(input)
```

### Desktop Database

```go
input := pgtune.TuneInput{
    DBVersion:      18,
    OSType:         pgtune.OSWindows,
    DBType:         pgtune.DBTypeDesktop,
    TotalMemory:    4 * pgtune.GB,
    CPUNum:         2,
    MaxConnections: 0, // Use default (20 for Desktop)
    StorageType:    pgtune.StorageTypeHDD,
}

output, _ := pgtune.Tune(input)
```

## Kubernetes Integration Example

### ConfigMap Generator

```go
package main

import (
    "fmt"
    "github.com/souravbiswassanto/go-pgtune/pkg/pgtune"
)

func GeneratePostgreSQLConfigMap(memory int64, cpu int) string {
    input := pgtune.TuneInput{
        DBVersion:      16,
        OSType:         pgtune.OSLinux,
        DBType:         pgtune.DBTypeWeb,
        TotalMemory:    memory,
        CPUNum:         cpu,
        MaxConnections: 100,
        StorageType:    pgtune.StorageTypeSSD,
    }

    output, _ := pgtune.Tune(input)

    config := fmt.Sprintf(`
apiVersion: v1
kind: ConfigMap
metadata:
  name: postgresql-config
data:
  postgresql.conf: |
    max_connections = %d
    shared_buffers = %s
    effective_cache_size = %s
    maintenance_work_mem = %s
    checkpoint_completion_target = %.1f
    wal_buffers = %s
    default_statistics_target = %d
    random_page_cost = %.1f
    effective_io_concurrency = %d
    work_mem = %s
    huge_pages = %s
    min_wal_size = %s
    max_wal_size = %s
`,
        output.MaxConnections,
        output.SharedBuffers,
        output.EffectiveCacheSize,
        output.MaintenanceWorkMem,
        output.CheckpointCompletionTarget,
        output.WalBuffers,
        output.DefaultStatisticsTarget,
        output.RandomPageCost,
        *output.EffectiveIOConcurrency,
        output.WorkMem,
        output.HugePages,
        output.MinWalSize,
        output.MaxWalSize,
    )

    return config
}
```

## Environment Variables

The server supports the following environment variables:

- `PORT`: Server port (default: 8080)

Example:
```bash
PORT=9090 go run server.go
```

## Testing

Run all tests:
```bash
make test
```

Or:
```bash
go test ./...
```

Run tests with verbose output:
```bash
go test ./pkg/pgtune -v
```

## API Reference

### POST /tune

**Request:**
```json
{
  "db_version": 18,           // PostgreSQL version
  "os_type": "linux",         // linux, windows, mac
  "db_type": "web",           // web, oltp, dw, desktop, mixed
  "total_memory_gb": 2,       // Total RAM in GB
  "cpu_num": 1,               // Number of CPU cores
  "max_connections": 100,     // Max connections (0 for default)
  "storage_type": "ssd"       // ssd, hdd, san
}
```

**Response (Success):**
```json
{
  "success": true,
  "config": {
    "MaxConnections": 100,
    "SharedBuffers": "512MB",
    "EffectiveCacheSize": "1536MB",
    "MaintenanceWorkMem": "128MB",
    "CheckpointCompletionTarget": 0.9,
    "WalBuffers": "16MB",
    "DefaultStatisticsTarget": 100,
    "RandomPageCost": 1.1,
    "EffectiveIOConcurrency": 200,
    "WorkMem": "4854kB",
    "HugePages": "off",
    "MinWalSize": "1GB",
    "MaxWalSize": "4GB",
    "MaxWorkerProcesses": null,
    "MaxParallelWorkersPerGather": null,
    "MaxParallelWorkers": null,
    "MaxParallelMaintenanceWorkers": null,
    "WalLevel": "",
    "MaxWalSenders": ""
  }
}
```

**Response (Error):**
```json
{
  "success": false,
  "error": "total memory must be greater than 0"
}
```

### GET /health

**Response:**
```json
{
  "status": "ok"
}
```

## Common Scenarios

### Scenario 1: Small Web Application
- Memory: 2 GB
- CPU: 1 core
- Type: web
- Storage: SSD

```bash
curl -X POST http://localhost:8080/tune \
  -H "Content-Type: application/json" \
  -d '{"db_version":18,"os_type":"linux","db_type":"web","total_memory_gb":2,"cpu_num":1,"max_connections":100,"storage_type":"ssd"}'
```

### Scenario 2: Medium OLTP Database
- Memory: 16 GB
- CPU: 4 cores
- Type: OLTP
- Storage: SSD

```bash
curl -X POST http://localhost:8080/tune \
  -H "Content-Type: application/json" \
  -d '{"db_version":18,"os_type":"linux","db_type":"oltp","total_memory_gb":16,"cpu_num":4,"max_connections":0,"storage_type":"ssd"}'
```

### Scenario 3: Large Data Warehouse
- Memory: 128 GB
- CPU: 32 cores
- Type: DW
- Storage: SSD

```bash
curl -X POST http://localhost:8080/tune \
  -H "Content-Type: application/json" \
  -d '{"db_version":18,"os_type":"linux","db_type":"dw","total_memory_gb":128,"cpu_num":32,"max_connections":0,"storage_type":"ssd"}'
```

## Troubleshooting

### Port Already in Use

If port 8080 is already in use:
```bash
PORT=8081 go run server.go
```

### Import Issues

Make sure you're in the correct directory and run:
```bash
go mod tidy
```

### Test Failures

Run verbose tests to see details:
```bash
go test ./pkg/pgtune -v
```

## Additional Resources

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [pgtune Website](https://pgtune.leopard.in.ua/)
- [PostgreSQL Configuration Best Practices](https://wiki.postgresql.org/wiki/Tuning_Your_PostgreSQL_Server)
