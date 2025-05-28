# gRPC Load Balancing - Step 12: Service Discovery and Load Balancing

This project demonstrates how to implement client-side load balancing in gRPC using the built-in round-robin policy and service discovery. It shows how to distribute client requests across multiple server instances for better scalability and high availability.

## Project Structure

```
.
├── cmd/              # Command-line applications
│   ├── client/      # gRPC client with load balancing
│   └── server/      # gRPC server implementation
├── internal/         # Internal packages
│   └── greeter/     # Generated protobuf code
├── proto/           # Protocol buffer definitions
│   └── greeter.proto
├── go.mod           # Go module configuration
├── go.sum           # Dependency checksums
├── Makefile         # Build automation
└── README.md        # This file
```

## Key Features

- Demonstrates client-side load balancing with round-robin policy
- Supports multiple service discovery mechanisms (DNS, static list)
- Shows how to configure gRPC client for load balancing
- Includes example of running multiple server instances
- Demonstrates request distribution across available servers

## Setup and Usage

### Prerequisites
- Go 1.16 or higher
- Protocol Buffers compiler (protoc)
- gRPC Go plugins

### 1. Initialize the project
```bash
make init
```
This will:
- Set up the Go module
- Install required protobuf and gRPC plugins
- Clean up any existing generated code

### 2. Generate protobuf code
```bash
make generate
```
This will generate the necessary Go code from the protobuf definitions.

### 3. Run multiple server instances
Open multiple terminal windows and run the following commands in different terminals to start multiple server instances:

```bash
# Terminal 1 - First server instance
SERVER_PORT=50051 make run-server

# Terminal 2 - Second server instance
SERVER_PORT=50052 make run-server
```

### 4. Run the client with load balancing
In a new terminal, run the client:
```bash
make run-client
```

The client is configured to use round-robin load balancing across all available server instances.

## How It Works

### Client-Side Load Balancing

The client is configured to use gRPC's built-in round-robin load balancing policy:

```go
conn, err := grpc.Dial(
    serviceAddress,
    grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
    grpc.WithInsecure(),
    grpc.WithBlock(),
)
```

### Service Discovery

This example demonstrates two service discovery approaches:

1. **DNS-based discovery**:
   ```go
   // Format: dns:///host1:port1,host2:port2
   const serviceAddress = "dns:///localhost:50051,localhost:50052"
   ```

2. **Static list of addresses**:
   ```go
   // Can be dynamically populated from a service registry
   const serviceAddress = "static:///localhost:50051,localhost:50052"
   ```

### Load Balancing Policies

gRPC supports several load balancing policies:

1. **round_robin**: Distributes requests sequentially across all available servers
2. **pick_first**: Always uses the first available server (default)
3. **grpclb**: For use with external load balancers

## Important Notes

### Performance Considerations
- Client-side load balancing reduces the need for a centralized load balancer
- The round-robin policy is simple but doesn't account for server load or health
- For production, consider more sophisticated load balancing strategies

### Health Checking
- gRPC supports health checking via the health check protocol
- Implement health checks to ensure traffic is only sent to healthy instances

### Service Discovery
- In production, use a service registry like etcd, Consul, or Kubernetes services
- The resolver can be customized to work with different service discovery backends

### Error Handling
- The client will automatically retry failed requests on other available servers
- Implement proper error handling and circuit breakers

## Example Output

When running the example, you should see output similar to:

```
# Server 1
🚀 Starting the load-balanced gRPC server on :50051...

# Server 2
🚀 Starting the load-balanced gRPC server on :50052...

# Client
🚀 Running the gRPC client with load balancing...
Response from server: Hello World (from :50051)
Response from server: Hello World (from :50052)
Response from server: Hello World (from :50051)
Response from server: Hello World (from :50052)
```

## Next Steps

- Implement health checking for server instances
- Add metrics and monitoring for load balancing
- Explore more advanced load balancing strategies
- Integrate with service mesh solutions like Istio or Linkerd
- Add circuit breaking and retry policies

## Troubleshooting

- If you see connection errors, ensure all server instances are running
- Verify that the service address format is correct
- Check that the load balancing policy is properly configured in the client
- Ensure proper error handling for cases when servers become unavailable