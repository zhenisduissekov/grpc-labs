.PHONY: all init generate clean run-server run-client jaeger

# Go parameters
GOCMD = go
GOBIN = $(shell go env GOPATH)/bin
PROTOC_GEN_GO = $(GOBIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC = $(GOBIN)/protoc-gen-go-grpc

# Initialize the project
init:
	rm -rf go.mod go.sum
	$(GOCMD) mod init step-09_opentelemetry_tracing
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

# Start Jaeger in a container
jaeger:
	docker run -d --name jaeger \
	  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
	  -e COLLECTOR_OTLP_ENABLED=true \
	  -p 6831:6831/udp \
	  -p 6832:6832/udp \
	  -p 5778:5778 \
	  -p 16686:16686 \
	  -p 4317:4317 \
	  -p 4318:4318 \
	  -p 14250:14250 \
	  -p 14268:14268 \
	  -p 14269:14269 \
	  -p 9411:9411 \
	  jaegertracing/all-in-one:latest

# Stop Jaeger container
stop-jaeger:
	docker stop jaeger
	docker rm jaeger

# Install dependencies
deps:
	$(GOCMD) get -u \
	  go.opentelemetry.io/otel \
	  go.opentelemetry.io/otel/exporters/jaeger \
	  go.opentelemetry.io/otel/sdk/trace \
	  go.opentelemetry.io/otel/sdk/resource \
	  go.opentelemetry.io/otel/propagation \
	  go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc \
	  google.golang.org/grpc

# Run all tests
test:
	$(GOCMD) test -v ./...
