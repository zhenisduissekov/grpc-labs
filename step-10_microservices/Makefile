GOBIN = $(shell go env GOPATH)/bin
PROTOC_GEN_GO = $(GOBIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC = $(GOBIN)/protoc-gen-go-grpc

.PHONY: all generate init run-logger run-server run-client clean

all: generate

init:
	rm -f go.mod go.sum
	rm -f internal/server/*.pb.go
	rm -f internal/logger/*.pb.go
	go mod init step-10_microservices
	go mod tidy
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "✅ Installed protoc-gen-go and protoc-gen-go-grpc"
	@echo "👉 Make sure this is in your shell config:"
	@echo 'export PATH="$$PATH:$$(go env GOPATH)/bin"'

check:
	@echo "protoc-gen-go: $(shell which protoc-gen-go || echo Not found)"
	@echo "protoc-gen-go-grpc: $(shell which protoc-gen-go-grpc || echo Not found)"

generate:
	mkdir -p internal/server internal/logger
	PATH="$(shell go env GOPATH)/bin:$$PATH" protoc --go_out=. --go-grpc_out=. proto/server.proto
	PATH="$(shell go env GOPATH)/bin:$$PATH" protoc --go_out=. --go-grpc_out=. proto/logger.proto

run-logger:
	@echo "Starting Logger Service (port 50052)..."
	@go run cmd/logger/main.go

run-server:
	@echo "Starting Server Service (port 50051)..."
	@go run cmd/server/main.go

NAME ?= World
run-client:
	@echo "Running gRPC client with name: $(NAME)"
	@go run cmd/client/main.go $(NAME)

clean:
	rm -f internal/server/*.pb.go
	rm -f internal/logger/*.pb.go
	rm -f go.mod go.sum
