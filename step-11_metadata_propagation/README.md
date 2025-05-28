# gRPC Metadata Propagation - Step 11

This project demonstrates metadata propagation between gRPC services. It shows how to pass metadata (like user IDs, authentication tokens, etc.) from a client through one service to another service.

## Project Structure

```
.
├── cmd/              # Command-line applications
│   ├── client/      # gRPC client implementation
│   └── server/      # gRPC server implementation
├── internal/         # Internal packages
│   ├── greeter/     # Generated protobuf code for Greeter service
│   └── logger/      # Generated protobuf code for Logger service
├── proto/           # Protocol buffer definitions
│   ├── greeter.proto
│   └── logger.proto
├── go.mod           # Go module configuration
├── go.sum           # Dependency checksums
├── Makefile         # Build automation
└── README.md        # This file
```

## Key Features

- Demonstrates bidirectional metadata propagation between gRPC services
- Shows how to extract, modify, and forward metadata in middleware
- Implements context-based request scoping with metadata
- Includes example of metadata validation and sanitization
- Demonstrates practical use cases for metadata in microservices

## Setup and Usage

1. Initialize the project:
```bash
make init
```
This will:
- Clean up existing Go module files and generated code
- Initialize a new Go module
- Install required protobuf and gRPC plugins
- Display the paths of installed protobuf tools

2. Generate protobuf code:
```bash
make generate
```
This will generate Go code from the protobuf definitions in the `proto/` directory.

3. Run the server and client:
```bash
make run
```
This will:
- Start the Greeter service (listening on port 50051)
- Start the Logger service (listening on port 50052)
- Run the client which will:
  - Create a context with metadata (x-user-id and x-request-id)
  - Call the Greeter service
  - The Greeter service will forward the call to the Logger service
  - Both services will log the received metadata

## How It Works

### Client-Side Metadata
1. The client creates metadata using `metadata.New`:
   ```go
   md := metadata.Pairs(
       "x-user-id", "user-123",
       "x-request-id", uuid.New().String(),
   )
   ctx := metadata.NewOutgoingContext(context.Background(), md)
   ```

### Server-Side Metadata Handling
1. The Greeter server extracts incoming metadata:
   ```go
   md, ok := metadata.FromIncomingContext(ctx)
   if ok {
       // Process metadata
   }
   ```

2. The Greeter service forwards metadata to the Logger service:
   ```go
   // Create a new context with the incoming metadata
   ctx = metadata.NewOutgoingContext(ctx, md)
   
   // Add or modify metadata before forwarding
   ctx = metadata.AppendToOutgoingContext(ctx, "x-forwarded-by", "greeter-service")
   
   // Call the Logger service
   response, err := loggerClient.Log(ctx, &loggerpb.LogRequest{...})
   ```

### Metadata Propagation Best Practices
1. **Always Validate Metadata**: Check for required metadata fields and validate their values
2. **Be Careful with Sensitive Data**: Don't log or forward sensitive metadata
3. **Use Context for Request Scoping**: Always pass metadata through context
4. **Add Request Tracing**: Include request IDs for distributed tracing
5. **Document Your Metadata**: Maintain clear documentation of all metadata fields

## Important Notes

1. **Metadata Format**
   - Keys are automatically converted to lowercase
   - Values are always strings (convert other types as needed)
   - Multiple values for the same key are supported

2. **Performance Considerations**
   - Keep metadata size small (preferably < 8KB)
   - Be mindful of the number of metadata fields
   - Consider using binary metadata for large values

3. **Security Considerations**
   - Never trust metadata from untrusted sources
   - Validate and sanitize all metadata values
   - Consider encrypting sensitive metadata

4. **Error Handling**
   - Always check if metadata exists before accessing it
   - Provide meaningful error messages for missing or invalid metadata
   - Log metadata-related errors appropriately

## Example Output

When you run the example, you should see output similar to:

```
[Greeter] Received request with metadata: map[x-request-id:[req-123] x-user-id:[user-123]]
[Logger] Received log request with metadata: map[x-forwarded-by:[greeter-service] x-request-id:[req-123] x-user-id:[user-123]]
[Client] Received response: Hello, World!
```

## Project Status

This is step 11 in a comprehensive gRPC tutorial series, focusing on metadata propagation between services. It builds upon the previous steps by adding cross-service context propagation using gRPC metadata.

## Next Steps

- Add authentication and authorization using metadata
- Implement distributed tracing with OpenTelemetry
- Add request validation middleware
- Explore more advanced metadata patterns like deadline propagation
