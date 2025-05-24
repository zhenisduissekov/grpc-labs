# gRPC with OpenTelemetry and Jaeger

This step demonstrates how to add distributed tracing to a gRPC service using OpenTelemetry and Jaeger.

## Features

- **Distributed Tracing**: End-to-end tracing of gRPC calls
- **Jaeger Integration**: Visualize traces in the Jaeger UI
- **Context Propagation**: Trace context is propagated between services
- **gRPC Interceptors**: Automatic trace creation for gRPC methods

## Prerequisites

- Docker (for running Jaeger)
- Go 1.16+

## Project Structure

```
.
├── cmd/
│   ├── client/       # gRPC client with tracing
│   └── server/       # gRPC server with tracing
├── internal/
│   └── greeter/     # Generated protobuf code
├── proto/            # Protocol buffer definitions
│   └── greeter.proto
├── go.mod           # Go module configuration
├── Makefile         # Build automation
└── README.md        # This file
```

## Setup and Usage

1. **Start Jaeger** (in a separate terminal):
   ```bash
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
   ```

2. **Run the server**:
   ```bash
   make run-server
   ```

3. **Run the client** (in a separate terminal):
   ```bash
   make run-client
   ```

4. **View traces** in Jaeger UI: http://localhost:16686

## Key Components

### OpenTelemetry Setup

```go
// Initialize tracer provider
tp, err := tracerProvider("greeter-service")
if err != nil {
    log.Fatalf("Failed to create tracer provider: %v", err)
}

// Register our TracerProvider as the global
otel.SetTracerProvider(tp)

// Set up context propagation
otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
    propagation.TraceContext{},
    propagation.Baggage{},
))
```

### gRPC Server with Tracing

```go
// Create gRPC server with OpenTelemetry interceptors
s := grpc.NewServer(
    grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
    grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
)
```

### gRPC Client with Tracing

```go
// Create gRPC connection with OpenTelemetry interceptors
conn, err := grpc.Dial(
    "localhost:50051",
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
    grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
)
```

## Dependencies

- `go.opentelemetry.io/otel` - OpenTelemetry Go SDK
- `go.opentelemetry.io/otel/exporters/jaeger` - Jaeger exporter
- `go.opentelemetry.io/otel/sdk` - OpenTelemetry SDK
- `go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc` - gRPC instrumentation

## Next Steps

- Add custom attributes to spans
- Add logging correlation with trace IDs
- Set up sampling for production use
- Add metrics alongside traces
