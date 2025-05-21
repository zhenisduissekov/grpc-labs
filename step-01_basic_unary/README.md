# Basic gRPC Tutorial - Step 01: Unary RPC

This project demonstrates a basic unary gRPC service in Go. It includes a simple "Greeter" service that implements a "SayHello" RPC method.

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

3. Run everything with a single command:
```bash
make run
```
This will:
- Start the gRPC server
- The client will send a request with the name "World" and you should see the greeting message "Hello world" printed in the client's terminal.

If you encounter an error like "program not found or is not executable", try adding the Go bin directory to your PATH:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

4. Verify the service with curl (optional):
```bash
curl -X POST -H "Content-Type: application/json" \
     -d '{"name": "World"}' \
     http://localhost:50051/greeter.SayHello
```
Note: The curl command requires the gRPC HTTP/2 proxy to be running. This is an alternative way to test the service but not recommended for production use.

## Important Notes

1. **Package Structure**:
   - The `go_package` option in `proto/greeter.proto` defines where the generated code will be placed
   - The Go module name in `go.mod` must match the package structure
   - Example: `go_package = "internal/greeter;greeterpb"` will create files in `internal/greeter/`
   - **Caution**: Changing either of these without updating the other can lead to incorrect package paths and compilation errors

2. **Build Commands**:
   - The server and client both use the generated code from `internal/greeter`
   - The server runs on port 50051
   - The client connects to localhost:50051

3. **Dependencies**:
   - Go 1.23 or higher
   - Protocol Buffers compiler (protoc)
   - gRPC Go plugins

## Project Status

This is the first step in a gRPC tutorial series, demonstrating basic unary RPC communication. Future steps will build upon this foundation.

## Important Notes

1. **Package Structure**:
   - The `go_package` option in `proto/greeter.proto` defines where the generated code will be placed
   - The Go module name in `go.mod` must match the package structure
   - Example: `go_package = "internal/greeter;greeterpb"` will create files in `internal/greeter/`
   - **Caution**: Changing either of these without updating the other can lead to incorrect package paths and compilation errors

2. **Makefile Commands**:
   - `make init`: Set up the project with required dependencies
   - `make generate`: Generate protobuf code
   - `make all`: Run generate (equivalent to `make generate`)

3. **Dependencies**:
   - Go 1.23 or higher
   - Protocol Buffers compiler (protoc)
   - gRPC Go plugins

## Project Status

This is the first step in a gRPC tutorial series, demonstrating basic unary RPC communication. Future steps will build upon this foundation.
