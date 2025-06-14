GOBIN = $(shell go env GOPATH)/bin
PROTOC_GEN_GO = $(GOBIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC = $(GOBIN)/protoc-gen-go-grpc

.PHONY: all generate init run build-server build-client run-server run-client run-metrics stop-metrics test vet fmt tidy clean

all: generate run

init:
	rm -f go.mod go.sum
	rm -f *.pb.go
	go mod init step-15_grafana_dashboards
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
	mkdir -p internal/greeter
	PATH="$(shell go env GOPATH)/bin:$$PATH" protoc --go_out=. --go-grpc_out=. proto/greeter.proto

build-server:
	@echo "Building server..."
	go build -o bin/server ./cmd/server

build-client:
	@echo "Building client..."
	go build -o bin/client ./cmd/client

run-server: build-server
	@echo "Starting server..."
	./bin/server

run-client: build-client
	@echo "Starting client..."
	./bin/client

run-metrics:
	@echo "Starting metrics stack (Prometheus + Grafana)..."
	docker compose -f docker-compose.yml up -d

stop-metrics:
	@echo "Stopping metrics stack..."
	docker compose -f docker-compose.yml down

test:
	@echo "Running tests..."
	go test -v ./...

vet:
	@echo "Running go vet..."
	go vet ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...

tidy:
	@echo "Tidying up dependencies..."
	go mod tidy

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf bin/*
	rm -f internal/greeter/*.pb.go
	rm -f internal/greeter/*_grpc.pb.go

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go get \
		github.com/grpc-ecosystem/go-grpc-prometheus \
		github.com/prometheus/client_golang/prometheus \
		github.com/prometheus/client_golang/prometheus/promhttp \
		google.golang.org/grpc \
		google.golang.org/protobuf/...

# Install protoc plugins
install-protoc-plugins:
	@echo "Installing protoc plugins..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Setup everything
setup: init deps generate tidy
