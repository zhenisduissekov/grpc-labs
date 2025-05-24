# gRPC Tutorial - Step 07: Server Reflection and Health Checking

This project demonstrates how to enable gRPC server reflection and health checking in a gRPC service. These features are essential for service discovery, debugging, and monitoring in production environments.

## Project Structure

```
.
├── cmd/              # Command-line applications
│   ├── client/      # gRPC client implementation
│   └── server/      # gRPC server implementation
├── internal/         # Internal packages
│   └── greeter/     # Generated protobuf code
├── proto/           # Protocol buffer definitions
│   └── greeter.proto
├── go.mod           # Go module configuration
├── Makefile         # Build automation
└── README.md        # This file
```

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

3. Run the server:
```bash
make run-server
```

4. In a separate terminal, you can use `grpcurl` to explore the service:
```bash
# List all services
grpcurl -plaintext localhost:50051 list

# Check service health
grpcurl -plaintext -d '{"service":"greeter.Greeter"}' localhost:50051 grpc.health.v1.Health/Check

# Call the SayHello method
grpcurl -plaintext -d '{"name":"World"}' localhost:50051 greeter.Greeter/SayHello
```

## Key Features

### Server Reflection
Enables tools like `grpcurl` and `grpc_cli` to dynamically discover and call services without requiring pre-generated client code.

### Health Checking
Implements the standard gRPC health checking protocol, allowing load balancers and orchestration systems to monitor service health.

## Important Notes

1. **Package Structure**:
   - The `go_package` option in `proto/greeter.proto` defines where the generated code will be placed
   - The Go module name in `go.mod` must match the package structure
   - Example: `go_package = "internal/greeter;greeterpb"` will create files in `internal/greeter/`
   - **Caution**: Changing either of these without updating the other can lead to incorrect package paths and compilation errors

2. **Makefile Commands**:
   - `make init`: Set up the project with required dependencies
   - `make generate`: Generate protobuf code
   - `make run-server`: Start the gRPC server with reflection and health checking enabled
   - `make run-client`: Run the example client (to be implemented)

3. **Dependencies**:
   - Go 1.23 or higher
   - Protocol Buffers compiler (protoc)
   - gRPC Go plugins
   - grpcurl (for testing)

## Project Status

This is the seventh step in a gRPC tutorial series, demonstrating server reflection and health checking. The service includes:
- Unary RPC (SayHello)
- Server streaming RPC (StreamGreetings)
- Bidirectional streaming RPC (Chat)
- Server reflection for service discovery
- Standard health checking protocol
