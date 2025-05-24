# gRPC Microservices - Step 10: Service-to-Service Communication

This project demonstrates gRPC service-to-service communication between two microservices: `ServerService` and `LoggerService`. It simulates a real-world microservice architecture where one service calls another.

## Project Structure

```
.
├── cmd/                  # Command-line applications
│   ├── greeter/         # Greeter service (calls LoggerService)
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

### 1. Server Service (port 50051)
- Provides a simple greeting functionality
- Calls the Logger service for each request
- Implements the `Server` service defined in `proto/server.proto`

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

# Terminal 2 - Start Greeter Service
make run-greeter

# Terminal 3 - Run the client
make run-client
```

## Protocol Buffer Definitions

### `proto/logger.proto`
Defines the Logger service with a simple Log RPC method.

### `proto/greeter.proto`
Defines the Greeter service that depends on the Logger service.

## Important Notes

1. **Service Dependencies**:
   - The Greeter service depends on the Logger service
   - Logger service must be running before starting the Greeter service

2. **Port Configuration**:
   - Greeter Service: 50051
   - Logger Service: 50052

3. **Error Handling**:
   - The Greeter service includes basic error handling for Logger service unavailability

## Next Steps

1. Add service discovery
2. Implement retries and circuit breaking
3. Add tracing between services
4. Add metrics and monitoring
5. Containerize the services
