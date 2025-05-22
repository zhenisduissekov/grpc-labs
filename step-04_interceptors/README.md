# Advanced gRPC Tutorial - Step 04: Interceptors and Streaming

This project demonstrates advanced gRPC features in Go, including interceptors, streaming RPCs, and service logging. It builds upon basic unary RPC concepts and adds:

- Unary RPC with interceptors
- Server streaming RPC
- Bidirectional streaming RPC
- Client streaming RPC
- Service logging
- Health checking
- Service reflection

## Project Structure

```
.
├── cmd/              # Command-line applications
│   ├── client/      # gRPC client implementation
│   ├── server/      # gRPC server implementation
│   └── logger/      # Logging service implementation
├── internal/         # Internal packages
│   ├── greeter/     # Generated protobuf code for greeter service
│   └── logger/      # Generated protobuf code for logger service
├── proto/           # Protocol buffer definitions
│   ├── greeter.proto
│   └── logger.proto
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
This will generate Go code from the protobuf definitions in `proto/greeter.proto` and `proto/logger.proto`.

If you encounter an error like "program not found or is not executable", try adding the Go bin directory to your PATH:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

3. Run the services:

First, start the logger service:
```bash
make run-logger
```

Then, start the main server:
```bash
make run-server
```

Finally, run the client:
```bash
make run-client
```

The client will:
1. Send unary RPC requests with metadata (authorization token and user ID)
2. Receive server streaming greetings
3. Engage in bidirectional chat
4. Upload names using client streaming

You should see logs from both the server and logger services as the client makes requests.

## Features

### Interceptors
- Unary interceptor for logging all RPC calls
- Stream interceptor for logging streaming RPCs
- Metadata handling for authentication and user tracking

### Streaming RPCs
- Server streaming: Server sends multiple responses to a single client request
- Bidirectional streaming: Both client and server send multiple messages
- Client streaming: Client sends multiple requests to a single server response

### Service Logging
- Dedicated logger service running on port 50052
- All RPC calls are logged through the logger service
- Metadata is preserved in log messages

### Additional Features
- Health checking endpoint
- Service reflection for discovery
- Proper error handling and context management

## Cleanup

To stop all services:
```bash
lsof -ti:50051,50052 | xargs kill -15
```

This will terminate both the main server (port 50051) and logger service (port 50052).

---
### 🛠️ Built with Help from AI

This project was developed primarily by me while learning and exploring gRPC and Go — with the assistance of OpenAI's ChatGPT and [Windsurf](https://windsurf.ai), a powerful AI coding agent platform. All code has been reviewed and customized for clarity and educational value.
