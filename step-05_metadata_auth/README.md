# gRPC Metadata-Based Authentication

This project demonstrates how to implement metadata-based authentication in gRPC using interceptors. It builds upon the basic unary RPC example and adds authentication using metadata headers.

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

## Features

- Basic unary RPC with `SayHello`
- Secure RPC with `SecureGreeting` that requires authentication
- Server-side interceptor for metadata validation
- Client-side interceptor for adding authentication tokens

## Setup and Usage

1. Initialize the project:
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

4. In a separate terminal, run the client:
```bash
make run-client
```

## Authentication Flow

1. The client adds an authentication token to the metadata:
   ```go
   md := metadata.Pairs("authorization", "bearer my-secret-token")
   ctx := metadata.NewOutgoingContext(context.Background(), md)
   ```

2. The server intercepts the request and validates the token:
   ```go
   md, ok := metadata.FromIncomingContext(ctx)
   if !ok {
       return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
   }
   
   // Validate token from metadata
   if !isValidToken(md["authorization"]) {
       return nil, status.Error(codes.Unauthenticated, "invalid token")
   }
   ```

## Testing

1. The client will demonstrate two scenarios:
   - A successful unauthenticated call to `SayHello`
   - A successful authenticated call to `SecureGreeting`
   - A failed unauthenticated call to `SecureGreeting`

## Important Notes

- This is a basic example for demonstration purposes.
- In production, use proper token validation and secure token storage.
- Always use TLS in production environments.
- The token in this example is hardcoded for demonstration only.
