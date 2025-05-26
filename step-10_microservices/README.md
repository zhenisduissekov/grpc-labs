# gRPC Microservices - Step 10: Service-to-Service Communication

This project demonstrates gRPC service-to-service communication between two microservices: `ServerService` and `LoggerService`. It simulates a real-world microservice architecture where one service calls another.

## Project Structure

```
.
├── cmd/                  # Command-line applications
│   ├── server/         # Server service (calls LoggerService)
│   └── logger/          # Logger service (handles logging requests)
├── internal/             # Internal packages
│   ├── server/          # Server service protobuf definitions
│   └── logger/          # Logger service protobuf definitions
├── proto/               # Protocol buffer definitions
│   ├── server.proto     # Server service definition
│   └── logger.proto     # Logger service definition
├── go.mod              # Go module configuration
├── Makefile            # Build automation
└── README.md           # This file
```

## Services

### 1. Server (port 50051)
- Provides a simple greeting functionality
- Calls the Logger service for each request
- Implements the `Server` defined in `proto/server.proto`

### 2. Logger Service (port 50052)
- Handles logging requests from other services
- Implements the `Logger` service defined in `proto/logger.proto`

## Setup and Usage

1. Initialize the project:
```bash
make init
```
This command will:
- Clean up existing Go module files and generated code
- Initialize a new Go module
- Install required protobuf and gRPC plugins

2. Generate protobuf code:
```bash
make generate
```
This will generate Go code from the protobuf definitions.

3. Start the services:
```bash
# Terminal 1 - Start Logger Service
make run-logger

# Terminal 2 - Start Server Service
make run-server

# Terminal 3 - Run the client
make run-client
```

## Protocol Buffer Definitions

### `proto/logger.proto`
Defines the Logger service with a simple Log RPC method.

### `proto/server.proto`
Defines the Server that depends on the Logger service.

## Important Notes

1. **Service Dependencies**:
   - The Server depends on the Logger service.
   - Ensure the Logger service is running before starting the Server service.

2. **Port Configuration**:
   - Server Service: 50051
   - Logger Service: 50052

3. **Error Handling**:
   - The Server includes basic error handling for Logger service unavailability, such as retries and fallback mechanisms.
   - The Logger service includes basic error handling for log request failures, ensuring stability.

## Next Steps

1. Add service discovery using tools like Consul or Kubernetes.
2. Implement retries and circuit breaking with libraries like `github.com/sony/gobreaker`.
3. Add distributed tracing between services using OpenTelemetry.
4. Integrate metrics and monitoring with Prometheus and Grafana.
5. Containerize the services using Docker and orchestrate them with Kubernetes.
