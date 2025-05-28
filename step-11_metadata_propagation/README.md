# gRPC Metadata Propagation - Step 11

This project demonstrates metadata propagation between gRPC services. It shows how to pass metadata (like user IDs, authentication tokens, etc.) from a client through one service to another service.

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

## Key Features

- Demonstrates passing metadata between gRPC services
- Shows how to extract and propagate metadata in the server
- Uses context for request-scoped values

## Setup and Usage

1. Initialize the project:
```bash
make init
```

2. Generate protobuf code:
```bash
make generate
```

3. Run the server and client:
```bash
make run
```

## How It Works

1. The client sends a request with metadata (x-user-id)
2. The Greeter service receives the request and forwards the metadata to the Logger service
3. The Logger service extracts and logs the metadata

## Important Notes

- Metadata is passed using `metadata.AppendToOutgoingContext`
- Metadata is extracted using `metadata.FromIncomingContext`
- Always validate and sanitize metadata values

## Project Status

This is step 11 in a gRPC tutorial series, focusing on metadata propagation between services.
