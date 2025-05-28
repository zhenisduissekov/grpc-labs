package main

import (
	"context"
	"fmt"
	"log"
	"time"

	greeterpb "step-12_load_balancing/internal/greeter"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

const (
	// Using DNS resolver with multiple server addresses
	serviceAddress = "dns:///example.com"
	server1        = "localhost:50051"
	server2        = "localhost:50052"
)

// dnsResolverBuilder is a custom resolver that returns the list of server addresses
type dnsResolverBuilder struct {
	servers []string
}

func (b *dnsResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	addrs := make([]resolver.Address, len(b.servers))
	for i, s := range b.servers {
		addrs[i] = resolver.Address{Addr: s}
	}
	cc.UpdateState(resolver.State{
		Addresses: addrs,
	})
	return &dnsResolver{cc: cc}, nil
}

func (b *dnsResolverBuilder) Scheme() string { return "dns" }

type dnsResolver struct {
	cc resolver.ClientConn
}

func (r *dnsResolver) ResolveNow(resolver.ResolveNowOptions) {}
func (r *dnsResolver) Close()                                {}

func main() {
	// Register a custom resolver that knows about our servers
	dnsResolverBuilder := &dnsResolverBuilder{
		servers: []string{server1, server2},
	}
	resolver.Register(dnsResolverBuilder)

	// Create a connection to the servers with round-robin load balancing
	log.Printf("Connecting to servers: %s, %s...", server1, server2)
	conn, err := grpc.Dial(
		serviceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := greeterpb.NewGreeterClient(conn)

	// Send multiple requests to see load balancing in action
	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)

		// Add a small delay between requests
		time.Sleep(500 * time.Millisecond)

		// Create a unique name for each request
		name := fmt.Sprintf("World-%d", i+1)
		log.Printf("Sending request %d for name: %s", i+1, name)

		// Make the RPC call
		start := time.Now()
		response, err := client.SayHello(ctx, &greeterpb.HelloRequest{Name: name})
		elapsed := time.Since(start)

		// Handle the response
		if err != nil {
			log.Printf("Error calling SayHello for %s: %v", name, err)
		} else {
			log.Printf("Response %d: %s (took %v)", i+1, response.Message, elapsed)
		}

		cancel()
	}
}
