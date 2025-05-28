# gRPC Tutorial - Step 13: Retry and Timeout

This project demonstrates how to implement retry and timeout functionality in gRPC services using Go. It builds upon the basic unary RPC example from Step 1, adding resilience features to handle transient failures.

## Project Structure

```
.
├── cmd/              # Command-line applications
│   ├── client/      # gRPC client implementation with retry logic
│   └── server/      # gRPC server implementation with timeout handling
├── internal/         # Internal packages
│   └── greeter/     # Generated protobuf code
├── proto/           # Protocol buffer definitions
│   └── greeter.proto
├── go.mod           # Go module configuration
├── Makefile         # Build automation
└── README.md        # This file
```

## Features

- **Client-Side Retry Logic**: Implements exponential backoff for transient failures
- **Request Timeout**: Demonstrates setting timeouts on both client and server
- **Circuit Breaking**: Prevents cascading failures with basic circuit breaking
- **Context Propagation**: Proper context handling for cancellation and timeouts

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

3. Start the server:
```bash
make run-server
```

4. In a separate terminal, run the client:
```bash
make run-client
```

The client will demonstrate:
- Successful requests with timeout
- Retry behavior on transient failures
- Circuit breaking after multiple failures

## Key Implementation Details

### Client-Side Retry
```go
// Example retry logic with exponential backoff
for attempt := 0; attempt < maxRetries; attempt++ {
    resp, err := client.SayHello(ctx, req)
    if err == nil {
        return resp, nil
    }
    
    if status.Code(err) == codes.DeadlineExceeded {
        return nil, err // Don't retry on timeout
    }
    
    time.Sleep(time.Duration(math.Pow(2, float64(attempt))) * 100 * time.Millisecond)}
```

### Server-Side Timeout
```go
// Server handler with timeout
func (s *server) SayHello(ctx context.Context, req *greeterpb.HelloRequest) (*greeterpb.HelloReply, error) {
    // Create a new context with timeout
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()
    
    // Process the request
    // ...
    
    return &greeterpb.HelloReply{Message: "Hello " + req.Name}, nil
}
```

## Dependencies

- Go 1.16 or higher
- Protocol Buffers compiler (protoc)
- gRPC Go plugins
- `google.golang.org/grpc`
- `google.golang.org/protobuf`

## Running the Client

You can run the client using the Makefile, which is the recommended approach. Here are the available commands:

### Basic Usage
```bash
# Basic client (connects to default server on localhost:50051)
make run-client name=YourName

# Client with error simulation
make run-client-error name=Test

# Client with delay (in milliseconds)
make run-client-delay name=Test delay=1000

# Client with custom timeout
make run-client-timeout name=Test timeout=2s

# Client with all options
make run-client-all name=Test error=true delay=1000 timeout=5s
```

### Important Notes

1. **Server Must Be Running**
   First, start the server in a separate terminal:
   ```bash
   make run-server
   ```

2. **Resolving Dependencies**
   If you encounter module-related errors, run:
   ```bash
   go mod tidy
   make generate
   ```

3. **Manual Execution (Not Recommended)**
   While possible, manual execution is not recommended due to path resolution issues. If needed:
   ```bash
   cd /path/to/step-13_retry_timeout
   go run ./cmd/client/main.go -name YourName
   ```
   The Makefile handles all path resolution automatically, making it the preferred method.

## Next Steps

This implementation provides a foundation for building resilient gRPC services. Future enhancements could include:
- More sophisticated retry policies
- Distributed tracing integration
- Advanced circuit breaking with metrics
- Load testing and tuning
