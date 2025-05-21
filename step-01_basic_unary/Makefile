GOBIN = $(shell go env GOPATH)/bin
PROTOC_GEN_GO = $(GOBIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC = $(GOBIN)/protoc-gen-go-grpc

.PHONY: all generate init run

all: generate run

init:
	rm -f go.mod go.sum
	rm -f *.pb.go
	go mod init step-01_basic_unary
	go mod tidy
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	export PATH="$(go env GOPATH)/bin:$PATH"
	@echo "Checking protobuf tool paths:"
	which protoc-gen-go
	which protoc-gen-go-grpc

generate:
	mkdir -p internal/greeter
	export PATH="$(go env GOPATH)/bin:$PATH"
	protoc --go_out=. \
		--go-grpc_out=. \
		proto/greeter.proto

run-server:
	@echo "Starting gRPC server..."
	@go run cmd/server/main.go

run-client:
	@echo "Running gRPC client..."
	@go run cmd/client/main.go


