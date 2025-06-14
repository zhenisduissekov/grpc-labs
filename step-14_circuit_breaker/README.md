# gRPC Circuit Breaker - Step 14

This step demonstrates how to implement a Circuit Breaker pattern for gRPC services using the `github.com/sony/gobreaker` package. The Circuit Breaker helps prevent cascading failures by stopping calls to a failing service.

## Project Structure

```
.
├── cmd/              # Command-line applications
│   ├── client/      # gRPC client with circuit breaker
│   └── server/      # gRPC server implementation
├── internal/         # Internal packages
│   └── greeter/     # Generated protobuf code
├── proto/           # Protocol buffer definitions
│   └── greeter.proto
├── go.mod           # Go module configuration
├── Makefile         # Build automation
└── README.md        # This file
```

## Key Components

- **Circuit Breaker**: Wraps gRPC client calls to prevent cascading failures
- **gRPC Service**: Simple Greeter service with a SayHello method
- **Resilience**: Automatic retries and failure handling

## Setup and Usage

1. Initialize the project:
```bash
make init
```
This command will:
- Clean up existing Go module files and generated code
- Initialize a new Go module
- Install required protobuf and gRPC plugins
- Display the paths of installed protobuf tools

2. Generate protobuf code:
```bash
make generate
```
This will generate Go code from the protobuf definition in `proto/greeter.proto`.

If you encounter an error like "program not found or is not executable", try adding the Go bin directory to your PATH:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

3. Run the server and client:
```bash
# Start the server in one terminal
make run-server

# In another terminal, run the client with circuit breaker
make run-client
```

## Important Notes

1. **Circuit Breaker Configuration**:
   - The circuit breaker is configured with specific thresholds for failures
   - It will open the circuit after a certain number of consecutive failures
   - Half-open state allows testing if the service has recovered
   - Timeout settings control when to consider a request as failed

2. **Package Structure**:
   - The `go_package` option in `proto/greeter.proto` defines where the generated code will be placed
   - The Go module name in `go.mod` must match the package structure
   - **Caution**: Changing package paths without updating imports can lead to compilation errors

3. **Makefile Commands**:
   - `make init`: Set up the project with required dependencies
   - `make generate`: Generate protobuf code
   - `make run-server`: Start the gRPC server
   - `make run-client`: Run the client with circuit breaker
   - `make all`: Run all necessary setup commands

## Dependencies

- Go 1.23 or higher
- Protocol Buffers compiler (protoc)
- gRPC Go plugins
- `github.com/sony/gobreaker` for circuit breaker implementation

## What is Being Tested

This step demonstrates the Circuit Breaker pattern in a gRPC service. The implementation uses `github.com/sony/gobreaker` to wrap gRPC client calls, providing resilience against failing services.

### Key Components Being Tested:

1. **Circuit Breaker States**:
   - **Closed**: Normal operation, requests pass through
   - **Open**: Requests fail immediately without calling the service
   - **Half-Open**: Limited requests allowed to test if service has recovered

2. **Failure Handling**:
   - Automatic detection of failing services
   - Fallback behavior when service is unavailable
   - Automatic recovery when service becomes available again

## Expected Results

When you run the client (`make run-client`), you should see:

1. **Initial Requests**:
   ```
   Greeting: Hello World
   Error calling SayHello: rpc error: code = InvalidArgument desc = Name cannot be empty
   Greeting: Hello World
   Error calling SayHello: rpc error: code = InvalidArgument desc = Name cannot be empty
   ```

2. **Circuit Breaker Trips** (after 3 consecutive failures):
   ```
   Circuit Breaker 'greeter' changed from closed to open
   Error calling SayHello: circuit breaker is open
   ```

3. **After Timeout** (5 seconds in Open state):
   ```
   Circuit Breaker 'greeter' changed from open to half-open
   ```

4. **Recovery** (if service starts working):
   ```
   Circuit Breaker 'greeter' changed from half-open to closed
   Greeting: Hello World
   ```

## Why This Matters

1. **Resilience**: Prevents cascading failures in distributed systems
2. **Fail Fast**: Quickly rejects calls to failing services
3. **Self-Healing**: Automatically recovers when services become available
4. **Graceful Degradation**: Provides fallback behavior when services are unavailable

## Configuration

The circuit breaker is configured with these parameters:
- `MaxRequests: 3`: Number of requests allowed in half-open state
- `Interval: 10s`: Time after which the failure count is reset
- `Timeout: 5s`: Time in open state before moving to half-open
- `ReadyToTrip`: Trips after 3 consecutive failures

## Project Status

This implementation demonstrates a production-grade circuit breaker pattern for gRPC services. The circuit breaker helps maintain system stability by preventing cascading failures and providing automatic recovery when services become available again.
