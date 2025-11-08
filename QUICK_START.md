# Quick Start Guide

## 1. Run Example (Fastest way to test)

```bash
cd /home/saurov/go/src/github.com/temp/go-pgtune
go run ./cmd/example
```

Expected output matches pgtune website:
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

## 2. Start API Server

```bash
go run ./cmd/server
```

Server starts on http://localhost:8080

## 3. Test API

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

## 4. Run Tests

```bash
go test ./pkg/pgtune -v
```

## 5. Use as Library

```go
import "github.com/souravbiswassanto/go-pgtune/pkg/pgtune"

input := pgtune.TuneInput{
    DBVersion:      18,
    OSType:         pgtune.OSLinux,
    DBType:         pgtune.DBTypeWeb,
    TotalMemory:    2 * pgtune.GB,
    CPUNum:         1,
    MaxConnections: 100,
    StorageType:    pgtune.StorageTypeSSD,
}

output, err := pgtune.Tune(input)
```

## Constants Available

### OS Types
- `pgtune.OSLinux`
- `pgtune.OSWindows`
- `pgtune.OSMac`

### DB Types
- `pgtune.DBTypeWeb`
- `pgtune.DBTypeOLTP`
- `pgtune.DBTypeDW`
- `pgtune.DBTypeDesktop`
- `pgtune.DBTypeMixed`

### Storage Types
- `pgtune.StorageTypeSSD`
- `pgtune.StorageTypeHDD`
- `pgtune.StorageTypeSAN`

### Memory Units
- `pgtune.KB` = 1024
- `pgtune.MB` = 1048576
- `pgtune.GB` = 1073741824
- `pgtune.TB` = 1099511627776

## Make Commands

```bash
make help      # Show all commands
make run       # Start server
make example   # Run example
make test      # Run tests
make build     # Build binary
make fmt       # Format code
make vet       # Run go vet
```

## Docker Commands

```bash
# Build
docker build -t pgtune-server .

# Run
docker run -p 8080:8080 pgtune-server

# Test
curl http://localhost:8080/health
```

## API Endpoints

- `POST /tune` - Generate configuration
- `GET /health` - Health check

## Sample Configurations

### Small Web App (2GB RAM, 1 CPU)
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

### Medium OLTP (16GB RAM, 4 CPU)
```json
{
  "db_version": 18,
  "os_type": "linux",
  "db_type": "oltp",
  "total_memory_gb": 16,
  "cpu_num": 4,
  "max_connections": 0,
  "storage_type": "ssd"
}
```

### Large DW (128GB RAM, 32 CPU)
```json
{
  "db_version": 18,
  "os_type": "linux",
  "db_type": "dw",
  "total_memory_gb": 128,
  "cpu_num": 32,
  "max_connections": 0,
  "storage_type": "ssd"
}
```

## Files

- `pkg/pgtune/pgtune.go` - Core library
- `server.go` - HTTP server
- `example.go` - Example usage
- `README.md` - Full documentation
- `USAGE.md` - Detailed usage guide
- `PROJECT_SUMMARY.md` - Complete summary

## Verification

Compare output with https://pgtune.leopard.in.ua/

✅ Values should match exactly!
