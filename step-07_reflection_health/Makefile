.PHONY: all init generate clean run-server run-client

# Go parameters
GOCMD = go
GOBIN = $(shell go env GOPATH)/bin
PROTOC_GEN_GO = $(GOBIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC = $(GOBIN)/protoc-gen-go-grpc

# Initialize the project
init:
	rm -rf go.mod go.sum
	go mod init step-07_reflection_health
	go mod tidy
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate protobuf code
generate:
	mkdir -p internal/greeter
	PATH="$(shell go env GOPATH)/bin:$$PATH" protoc --go_out=. --go-grpc_out=. proto/greeter.proto
	go mod tidy

# Clean generated files
clean:
	rm -rf internal/
	rm -f go.sum

# Run the server
run-server:
	go run cmd/server/main.go

# Run the client
run-client:
	go run cmd/client/main.go
