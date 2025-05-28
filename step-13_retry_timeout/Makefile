GOBIN = $(shell go env GOPATH)/bin
PROTOC_GEN_GO = $(GOBIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC = $(GOBIN)/protoc-gen-go-grpc

.PHONY: all generate init run-server run-client

all: generate run-server

init:
	rm -f go.mod go.sum
	rm -f internal/greeter/*.pb.go
	go mod init step-13_retry_timeout
	go mod tidy
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate protobuf code
generate:
	mkdir -p internal/greeter
	PATH="$(shell go env GOPATH)/bin:$$PATH" protoc --go_out=. --go-grpc_out=. proto/greeter.proto

# Start the server
run-server:
	@echo "🚀 Starting gRPC server with retry/timeout support on port 50051..."
	@PORT=50051 go run cmd/server/main.go

# Start a second server instance
run-server-2:
	@echo "🚀 Starting second gRPC server with retry/timeout support on port 50052..."
	@PORT=50052 go run cmd/server/main.go

# Run both server instances
run-servers: run-server run-server-2

# Run the client with retry logic
run-client:
	@echo "🚀 Running gRPC client with retry logic..."
	@go run cmd/client/main.go -name $(name) -server localhost:50051

# Run the client with error simulation
run-client-error:
	@echo "🚀 Running gRPC client with error simulation..."
	@go run cmd/client/main.go -name $(name) -server localhost:50051 -error

# Run the client with delay
run-client-delay:
	@echo "🚀 Running gRPC client with delay..."
	@go run cmd/client/main.go -name $(name) -server localhost:50051 -delay $(delay)

# Run the client with custom timeout
run-client-timeout:
	@echo "🚀 Running gRPC client with custom timeout..."
	@go run cmd/client/main.go -name $(name) -server localhost:50051 -timeout $(timeout)

# Run the client with all options
run-client-all:
	@echo "🚀 Running gRPC client with all options..."
	@go run cmd/client/main.go -name $(name) -server localhost:50051 -error=$(error) -delay=$(delay) -timeout=$(timeout)

# Clean generated files
clean:
	rm -f internal/greeter/*.pb.go

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  init      - Initialize the project and install dependencies"
	@echo "  generate  - Generate protobuf code"
	@echo "  run       - Run both server and client (in separate terminals)"
	@echo "  run-server - Run the gRPC server"
	@echo "  run-client - Run the gRPC client"
	@echo "  clean     - Remove generated files"
	@echo "  help      - Show this help message"
