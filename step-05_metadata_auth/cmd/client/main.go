package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "step-05_metadata_auth/internal/greeter"
)

const (
	defaultName = "world"
	token      = "my-secret-token"
)

func main() {
	// Set up a connection to the server.
	addr := flag.String("addr", "localhost:50051", "the address to connect to")
	name := flag.String("name", defaultName, "Name to greet")
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new Greeter client
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Test unauthenticated call (should work)
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Printf("could not greet: %v", err)
	} else {
		log.Printf("Greeting (unauthenticated): %s", r.GetMessage())
	}

	// Test authenticated call without token (should fail)
	r, err = c.SecureGreeting(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Printf("Secure greeting without token failed (expected): %v", err)
	} else {
		log.Printf("Secure greeting (without token): %s", r.GetMessage())
	}

	// Create a new context with metadata for authentication
	md := metadata.Pairs("authorization", "bearer "+token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	// Test authenticated call with token (should work)
	r, err = c.SecureGreeting(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Printf("Secure greeting with token failed: %v", err)
	} else {
		log.Printf("Secure greeting (with token): %s", r.GetMessage())
	}
}
