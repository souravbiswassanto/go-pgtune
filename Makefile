.PHONY: help build run example test clean

help: ## Display this help message
@echo "Available targets:"
@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the server binary
go build -o bin/pgtune-server ./cmd/server

run: ## Run the server
go run ./cmd/server

example: ## Run the example program
go run ./cmd/example

test: ## Run tests
go test ./...

clean: ## Clean build artifacts
rm -rf bin/

install: ## Install dependencies
go mod download
go mod tidy

fmt: ## Format code
go fmt ./...

vet: ## Run go vet
go vet ./...
