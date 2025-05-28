GOBIN = $(shell go env GOPATH)/bin
PROTOC_GEN_GO = $(GOBIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC = $(GOBIN)/protoc-gen-go-grpc

.PHONY: all generate init run-server run-client

all: generate run-server

init:
	rm -f go.mod go.sum
	rm -f *.pb.go
	go mod init step-12_load_balancing
	go mod tidy
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "✅ Installed protoc-gen-go and protoc-gen-go-grpc"

generate:
	mkdir -p internal/greeter
	PATH="$(shell go env GOPATH)/bin:$$PATH" protoc --go_out=. --go-grpc_out=. proto/greeter.proto

run-server:
	@echo "🚀 Starting the load-balanced gRPC server on port 50051..."
	@PORT=50051 go run cmd/server/main.go

run-server-2:
	@echo "🚀 Starting the second load-balanced gRPC server on port 50052..."
	@PORT=50052 go run cmd/server/main.go

run-servers: run-server run-server-2

run-client:
	@echo "🚀 Running the gRPC client with load balancing..."
	@go run cmd/client/main.go