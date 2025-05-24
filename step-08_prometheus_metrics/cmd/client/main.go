package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	greeterpb "step-08_prometheus_metrics/internal/greeter"
)

func main() {
	// Set up a connection to the server
	conn, err := grpc.Dial("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a client
	c := greeterpb.NewGreeterClient(conn)

	// Test Unary RPC
	testSayHello(c)
	// Test Server Streaming
	testStreamGreetings(c)
	// Test Bidirectional Streaming
	testChat(c)
}

func testSayHello(c greeterpb.GreeterClient) {
	log.Println("--- Testing SayHello ---")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &greeterpb.HelloRequest{Name: "World"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Response: %s (at %v)", r.GetMessage(), r.GetTimestamp().AsTime())
}

func testStreamGreetings(c greeterpb.GreeterClient) {
	log.Println("\n--- Testing StreamGreetings ---")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stream, err := c.StreamGreetings(ctx, &greeterpb.HelloRequest{Name: "Streaming Client"})
	if err != nil {
		log.Fatalf("error creating stream: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("Stream ended: %v", err)
			break
		}
		log.Printf("Received: %s (at %v)", msg.GetMessage(), msg.GetTimestamp().AsTime())
	}
}

func testChat(c greeterpb.GreeterClient) {
	log.Println("\n--- Testing Chat ---")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := c.Chat(ctx)
	if err != nil {
		log.Fatalf("error creating chat: %v", err)
	}

	// Send a few messages
	messages := []string{"Hello", "How are you?", "Bye!"}
	for _, msg := range messages {
		log.Printf("Sending: %s", msg)
		if err := stream.Send(&greeterpb.HelloRequest{Name: msg}); err != nil {
			log.Fatalf("Failed to send: %v", err)
		}

		// Receive response
		reply, err := stream.Recv()
		if err != nil {
			log.Fatalf("Failed to receive: %v", err)
		}
		log.Printf("Received: %s (at %v)", reply.GetMessage(), reply.GetTimestamp().AsTime())
	}

	// Close the send direction
	if err := stream.CloseSend(); err != nil {
		log.Fatalf("Failed to close send: %v", err)
	}
}
