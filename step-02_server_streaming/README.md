# gRPC Server Streaming Demo

This project demonstrates server-side streaming in gRPC, where the server sends multiple responses to a single client request. This is particularly useful for scenarios like:
- Real-time updates
- Large data streams
- Event notifications

## Features

- Server-side streaming implementation (StreamGreetings RPC)
- Example client-server communication
- Makefile for easy setup and running

## Setup

1. Install dependencies:
   ```bash
   make init
   ```

2. Generate protobuf code:
   ```bash
   make generate
   ```

3. Run the server:
   ```bash
   make run-server
   ```

4. Run the client:
   ```bash
   make run-client
   ```

## Project Structure

- `cmd/server/`: Server implementation with streaming RPC
- `cmd/client/`: Client implementation that receives streamed responses
- `proto/`: Protocol buffer definitions with streaming RPC
- `internal/greeter/`: Generated protobuf code

## Protocol Details

The server streaming RPC is defined in `proto/greeter.proto`:
```protobuf
service Greeter {
    rpc StreamGreetings(StreamGreetingsRequest) returns (stream StreamGreetingsResponse) {}
}

message StreamGreetingsRequest {
    string name = 1;
    int32 count = 2;
}

message StreamGreetingsResponse {
    string greeting = 1;
}
```

## Usage Example

The server sends streaming responses with progress indicators. For example, if the client requests 5 responses:

```
Running gRPC client...
2025/05/21 14:59:12 Stream response: Hello World (1/5)
2025/05/21 14:59:13 Stream response: Hello World (2/5)
2025/05/21 14:59:14 Stream response: Hello World (3/5)
2025/05/21 14:59:15 Stream response: Hello World (4/5)
2025/05/21 14:59:16 Stream response: Hello World (5/5)
```

## Important Notes

1. **Streaming Behavior**:
   - The server maintains an open connection until all responses are sent
   - The client receives responses asynchronously
   - Both sides can cancel the stream at any time

2. **Makefile Commands**:
   - `make init`: Set up the project with required dependencies
   - `make generate`: Generate protobuf code
   - `make all`: Run generate (equivalent to `make generate`)

3. **Dependencies**:
   - Go 1.23 or higher
   - Protocol Buffers compiler (protoc)
   - gRPC Go plugins

## Project Status

This is the second step in a gRPC tutorial series, demonstrating server-side streaming capabilities. For the basic unary RPC implementation, see step-01_basic_unary.