package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "step-11_metadata_propagation/internal/greeter"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	// Prepare metadata
	md := metadata.Pairs(
		"x-user-id", "user123",
		"x-request-id", "req-12345",
	)

	// Create a new context with metadata
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Set a timeout for the RPC
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	// Make the RPC call
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "World"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Response: %s", r.Message)
}
