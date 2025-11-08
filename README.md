# go-pgtune

A Go library and HTTP API for generating PostgreSQL tuning configurations. This is a Go port of [pgtune](https://github.com/le0pard/pgtune) that replicates the same logic and calculations.

## Features

- 🚀 Go library for programmatic PostgreSQL tuning
- 🌐 HTTP API server for easy integration
- ✅ 100% compatible with pgtune.leopard.in.ua output
- 📦 Zero external dependencies for the core library
- 🔧 Kubernetes-ready for cloud-native controllers

## Installation

### As a Library

```bash
go get github.com/souravbiswassanto/go-pgtune/pkg/pgtune
```

### As a Server

```bash
git clone https://github.com/souravbiswassanto/go-pgtune.git
cd go-pgtune
go run server.go
```

## Usage

### Library Usage

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
        TotalMemory:    2 * pgtune.GB, // 2 GB
        CPUNum:         1,
        MaxConnections: 100,
        StorageType:    pgtune.StorageTypeSSD,
    }

    output, err := pgtune.Tune(input)
    if err != nil {
        panic(err)
    }

    fmt.Printf("max_connections = %d\n", output.MaxConnections)
    fmt.Printf("shared_buffers = %s\n", output.SharedBuffers)
    fmt.Printf("effective_cache_size = %s\n", output.EffectiveCacheSize)
    // ... more parameters
}
```

### HTTP API Usage

Start the server:
```bash
go run server.go
# Server starting on port 8080
```

Make a POST request to `/tune`:

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

Response:
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
    "MaxWalSize": "4GB"
  }
}
```

## Parameters

### Input Parameters

| Parameter | Type | Description | Values |
|-----------|------|-------------|--------|
| `db_version` | float64 | PostgreSQL version | 18, 17, 16, 15, 14, 13, 12, 11, 10 |
| `os_type` | string | Operating system | `linux`, `windows`, `mac` |
| `db_type` | string | Database workload type | `web`, `oltp`, `dw`, `desktop`, `mixed` |
| `total_memory_gb` | float64 | Total RAM in GB | Any positive number |
| `cpu_num` | int | Number of CPU cores | Any positive integer |
| `max_connections` | int | Max connections (0 for default) | 0 or any positive integer |
| `storage_type` | string | Storage type | `ssd`, `hdd`, `san` |

### Output Parameters

The tuning configuration includes:
- `max_connections`
- `shared_buffers`
- `effective_cache_size`
- `maintenance_work_mem`
- `checkpoint_completion_target`
- `wal_buffers`
- `default_statistics_target`
- `random_page_cost`
- `effective_io_concurrency`
- `work_mem`
- `huge_pages`
- `min_wal_size`
- `max_wal_size`
- `max_worker_processes` (when CPU >= 4)
- `max_parallel_workers_per_gather` (when CPU >= 4)
- `max_parallel_workers` (when CPU >= 4, PG >= 10)
- `max_parallel_maintenance_workers` (when CPU >= 4, PG >= 11)

## Example Output

For the following input:
```
DB Version: 18
OS Type: linux
DB Type: web
Total Memory (RAM): 2 GB
CPUs num: 1
Connections num: 100
Data Storage: ssd
```

Output:
```
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

## Running the Example

```bash
go run example.go
```

## Testing

To verify the output matches pgtune website:

1. Visit https://pgtune.leopard.in.ua/
2. Enter the same parameters
3. Compare with the output from go-pgtune

The values should match exactly!

## API Endpoints

### POST /tune
Generate PostgreSQL tuning configuration

**Request Body:**
```json
{
  "db_version": 18,
  "os_type": "linux",
  "db_type": "web",
  "total_memory_gb": 2,
  "cpu_num": 1,
  "max_connections": 100,
  "storage_type": "ssd"
}
```

### GET /health
Health check endpoint

**Response:**
```json
{
  "status": "ok"
}
```

## Use Cases

- Kubernetes Operators for PostgreSQL
- Cloud-native database controllers
- Automated database provisioning
- Configuration management tools
- CI/CD pipelines for database setup

## License

This project is inspired by and maintains compatibility with [pgtune](https://github.com/le0pard/pgtune).

## Credits

Based on the excellent work of [pgtune](https://github.com/le0pard/pgtune) by Alexey Vasiliev.