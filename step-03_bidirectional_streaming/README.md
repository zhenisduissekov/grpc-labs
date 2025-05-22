# Step 03: Bidirectional Streaming

This example demonstrates bidirectional streaming with gRPC, where both the client and server can send messages to each other.

## Proto Definition

The proto file defines a bidirectional streaming chat service:

```protobuf
syntax = "proto3";

package chat;

service HelloService {
    rpc Chat(stream HelloRequest) returns (stream HelloReply) {}
}

message HelloRequest {
    string message = 1;
}

message HelloReply {
    string message = 1;
}
```

## Usage

### Build and Run

1. Generate proto files:
```bash
make generate
```

2. Run the server:
```bash
make run-server
```

3. Run multiple clients:
```bash
make run-client
```

### Example Flow

1. Start the server:
```bash
Starting gRPC server...
2025/05/21 15:00:32 Server listening at 127.0.0.1:50051
```

2. Start multiple clients:
```bash
Client 1:
Type messages to send (or 'exit' to quit):

Client 2:
Type messages to send (or 'exit' to quit):
```

3. Send messages between clients:
```bash
Client 1:
Type messages to send (or 'exit' to quit):
Hello from Client 1

Client 2:
Type messages to send (or 'exit' to quit):
[Server] Hello from Client 1
```

## Key Features

- Bidirectional streaming between client and server
- Real-time message broadcasting to all connected clients
- Proper connection handling and cleanup
- Error handling and logging
- Multiple client support

## Project Structure

```
step-03_bidirectional_streaming/
├── cmd/
│   ├── client/
│   │   └── main.go
│   └── server/
│       └── main.go
├── internal/
│   └── chat/
│       ├── chat.pb.go
│       └── chat.proto
└── Makefile
```
