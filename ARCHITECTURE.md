# Project Architecture

## Why the `cmd/` Structure?

The `cmd/` directory structure is a **Go best practice** for organizing projects with multiple executable binaries. This is the standard layout recommended by the Go community.

## Directory Structure

```
go-pgtune/
├── cmd/                    # Command-line applications
│   ├── server/            # HTTP API server
│   │   └── main.go        # Server entry point
│   └── example/           # Example/demo application
│       └── main.go        # Example entry point
├── pkg/                   # Public library code
│   └── pgtune/           # Core tuning library
│       ├── pgtune.go     # Main logic
│       └── pgtune_test.go # Tests
├── Dockerfile             # Container build
├── Makefile              # Build automation
└── go.mod                # Module definition
```

## Why This Structure?

### 1. **Multiple Binaries**
The `cmd/` directory allows you to have multiple entry points:
- `cmd/server` - Production HTTP server
- `cmd/example` - Demo/testing application
- Future: `cmd/cli` - CLI tool

Each subdirectory in `cmd/` produces a separate binary when built.

### 2. **Clear Separation**
- **`cmd/`** - Executable entry points (thin wrappers)
- **`pkg/`** - Reusable library code (business logic)

This makes it clear what's a program vs. what's a library.

### 3. **Reusable Library**
The `pkg/pgtune` package can be imported by:
- Your own `cmd/server` and `cmd/example`
- External projects that need PostgreSQL tuning logic
- Kubernetes operators/controllers

### 4. **Go Community Standard**
This layout is used by major Go projects:
- Kubernetes
- Docker
- Prometheus
- And many more...

## How to Use

### Run the Example
```bash
go run ./cmd/example
```

### Run the Server
```bash
go run ./cmd/server
```

### Build Binaries
```bash
# Build server
go build -o bin/pgtune-server ./cmd/server

# Build example
go build -o bin/pgtune-example ./cmd/example
```

### Use as Library
```go
import "github.com/souravbiswassanto/go-pgtune/pkg/pgtune"

func main() {
    input := pgtune.TuneInput{...}
    output, _ := pgtune.Tune(input)
}
```

## Benefits

1. **Scalability**: Easy to add new commands (e.g., `cmd/cli`, `cmd/worker`)
2. **Maintainability**: Clear code organization
3. **Testability**: Library code is independent of entry points
4. **Reusability**: `pkg/` can be imported anywhere
5. **Standard**: Follows Go community conventions

## Build Outputs

When you build the project:

```bash
go build ./cmd/server    # Creates 'server' binary
go build ./cmd/example   # Creates 'example' binary
```

Or with custom names:

```bash
go build -o pgtune-server ./cmd/server
go build -o pgtune-example ./cmd/example
```

## Container Build

The Dockerfile builds the server:

```dockerfile
RUN go build -o pgtune-server ./cmd/server
```

This creates a single, optimized binary for deployment.

## Makefile Commands

```bash
make run       # Runs go run ./cmd/server
make example   # Runs go run ./cmd/example
make build     # Builds binary to bin/pgtune-server
```

## Comparison: Before vs. After

### ❌ Before (Not Recommended)
```
go-pgtune/
├── server.go      # main() function
├── example.go     # main() function
└── pgtune.go      # Logic mixed with main
```
Problems:
- Multiple main() in root causes confusion
- Can't have both as importable packages
- Not standard Go layout

### ✅ After (Recommended)
```
go-pgtune/
├── cmd/
│   ├── server/main.go    # Server entry point
│   └── example/main.go   # Example entry point
└── pkg/pgtune/           # Reusable library
    ├── pgtune.go
    └── pgtune_test.go
```
Benefits:
- Clear separation of concerns
- Standard Go layout
- Library can be imported
- Multiple binaries supported

## For Kubernetes Controllers

When using in a Kubernetes operator:

```go
import "github.com/souravbiswassanto/go-pgtune/pkg/pgtune"

func (r *PostgreSQLReconciler) Reconcile(ctx context.Context, req ctrl.Request) {
    // Use the library
    config, _ := pgtune.Tune(pgtune.TuneInput{...})
    
    // Apply config to PostgreSQL pods
}
```

The `cmd/` structure doesn't interfere - you just import `pkg/pgtune`.

## References

- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Organizing Go Code](https://go.dev/blog/organizing-go-code)
- [How to Structure a Go Project](https://www.youtube.com/watch?v=oL6JBUk6tj0)
