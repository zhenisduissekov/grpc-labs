GOBIN = $(shell go env GOPATH)/bin
PROTOC_GEN_GO = $(GOBIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC = $(GOBIN)/protoc-gen-go-grpc

.PHONY: all generate init run

all: generate run

init:
	rm -f go.mod go.sum
	rm -f *.pb.go
	go mod init step-02_server_streaming
	go mod tidy
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "✅ Installed protoc-gen-go and protoc-gen-go-grpc"
	@echo "👉 Make sure this is in your shell config:"
	@echo 'export PATH="$$PATH:$$($(GO) env GOPATH)/bin"'

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
	@echo "Running gRPC client..."
	@go run cmd/client/main.go
