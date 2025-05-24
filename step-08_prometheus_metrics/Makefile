.PHONY: all init generate clean run-server run-client

# Go parameters
GOCMD = go
GOBIN = $(shell go env GOPATH)/bin
PROTOC_GEN_GO = $(GOBIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC = $(GOBIN)/protoc-gen-go-grpc

# Initialize the project
init:
	rm -rf go.mod go.sum
	$(GOCMD) mod init step-08_prometheus_metrics
	$(GOCMD) mod tidy
	$(GOCMD) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GOCMD) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate protobuf code
generate:
	mkdir -p internal/greeter
	PATH="$(shell go env GOPATH)/bin:$$PATH" protoc --go_out=. --go-grpc_out=. proto/greeter.proto
	$(GOCMD) mod tidy

# Clean generated files
clean:
	rm -rf internal/
	rm -f go.sum

# Run the server
run-server:
	$(GOCMD) run ./cmd/server

# Run the client
run-client:
	$(GOCMD) run ./cmd/client

# Install dependencies
deps:
	$(GOCMD) get -u \
		github.com/grpc-ecosystem/go-grpc-prometheus \
		github.com/prometheus/client_golang/prometheus \
		github.com/prometheus/client_golang/prometheus/promhttp

# Run all tests
test:
	$(GOCMD) test -v ./...
