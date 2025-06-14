GOBIN = $(shell go env GOPATH)/bin
PROTOC_GEN_GO = $(GOBIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC = $(GOBIN)/protoc-gen-go-grpc

.PHONY: all generate init run run-server run-client clean check help

all: generate

init:
	rm -f go.mod go.sum
	rm -f *.pb.go
	go mod init step-14_circuit_breaker
	go get google.golang.org/grpc
	go get google.golang.org/protobuf
	go get github.com/sony/gobreaker
	go mod tidy
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "✅ Installed dependencies and protoc plugins"

check:
	@echo "protoc-gen-go: $(shell which protoc-gen-go || echo Not found)"
	@echo "protoc-gen-go-grpc: $(shell which protoc-gen-go-grpc || echo Not found)"

generate:
	mkdir -p internal/greeter
	PATH="$(shell go env GOPATH)/bin:$$PATH" protoc --go_out=. --go-grpc_out=. proto/greeter.proto

run-server:
	@echo "Starting gRPC server..."
	@go run cmd/server/main.go

run-client:
	@echo "Starting gRPC client with circuit breaker..."
	@go run cmd/client/main.go

run: run-server

clean:
	rm -rf internal/*.pb.go

help:
	@echo "Available targets:"
	@echo "  init        - Initialize the project and install dependencies"
	@echo "  generate    - Generate Go code from .proto files"
	@echo "  run-server  - Start the gRPC server"
	@echo "  run-client  - Run the client with circuit breaker"
	@echo "  run         - Run both server and client (server in background)"
	@echo "  clean       - Remove generated files"
	@echo "  check       - Check if required tools are installed"
