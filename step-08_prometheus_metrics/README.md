# gRPC with Prometheus Metrics

This step demonstrates how to add Prometheus metrics to a gRPC service, enabling monitoring of RPC counts, latencies, and errors.

## Features

- **Prometheus Metrics**: Exposes gRPC server metrics for monitoring
  - RPC call counts
  - Latency histograms
  - Error counts
- **Metrics Endpoint**: Serves metrics at `/metrics` on port `:9090`
- **gRPC Server**: Runs on port `:50051`

## Project Structure

```
.
├── cmd/              # Command-line applications
│   ├── client/      # gRPC client implementation
│   └── server/      # gRPC server with Prometheus metrics
├── internal/         # Internal packages
│   └── greeter/     # Generated protobuf code
├── proto/           # Protocol buffer definitions
│   └── greeter.proto
├── go.mod           # Go module configuration
├── Makefile         # Build automation
└── README.md        # This file
```

## Setup and Usage

1. **Initialize the project**:
   ```bash
   make init
   ```

2. **Generate protobuf code**:
   ```bash
   make generate
   ```

3. **Run the server**:
   ```bash
   make run-server
   ```
   The server will start and expose:
   - gRPC server on `:50051`
   - Metrics on `:9090/metrics`

4. **Run the client** (in a separate terminal):
   ```bash
   make run-client
   ```

5. **View metrics**:
   ```bash
   curl http://localhost:9090/metrics
   ```

## Key Components

### Prometheus Metrics
Metrics are exposed using `go-grpc-prometheus` and include:
- `grpc_server_started_total`
- `grpc_server_handled_total`
- `grpc_server_msg_received_total`
- `grpc_server_msg_sent_total`
- `grpc_server_handling_seconds` (histogram)

### Server Setup
The server configures Prometheus metrics and serves them alongside the gRPC server:

```go
// Create gRPC server with Prometheus interceptors
s := grpc.NewServer(
    grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
    grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
)

// Register service
greeterpb.RegisterGreeterServer(s, &server{})

// Enable Prometheus metrics
grpc_prometheus.EnableHandlingTimeHistogram()
grpc_prometheus.Register(s)

// Start metrics server
go func() {
    http.Handle("/metrics", promhttp.Handler())
    log.Printf("Starting metrics server on :9090")
    log.Fatal(http.ListenAndServe(":9090", nil))
}()
```

## Dependencies

- `github.com/grpc-ecosystem/go-grpc-prometheus` - Prometheus integration for gRPC
- `github.com/prometheus/client_golang` - Prometheus Go client
- `google.golang.org/grpc` - gRPC framework

## Next Steps

- Set up Prometheus and Grafana for visualization
- Add custom metrics specific to your application
- Configure alerting based on metrics
