package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"

	greeterpb "step-07_reflection_health/internal/greeter"
)

const (
	address = "localhost:50051"
)

func checkHealth(client grpc_health_v1.HealthClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{
		Service: "greeter.Greeter",
	})

	if err != nil {
		log.Printf("Health check failed: %v", err)
		return
	}

	log.Printf("Service status: %s", resp.Status.String())
}

func callSayHello(client greeterpb.GreeterClient, name string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.SayHello(ctx, &greeterpb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("Greeting: %s at %v", r.GetMessage(), r.GetTimestamp().AsTime().Format(time.RFC3339))
}

func callStreamGreetings(client greeterpb.GreeterClient, name string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := client.StreamGreetings(ctx, &greeterpb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("Could not stream greetings: %v", err)
	}

	for {
		greeting, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while streaming: %v", err)
		}
		log.Printf("Stream greeting: %s at %v", 
			greeting.GetMessage(), 
			greeting.GetTimestamp().AsTime().Format(time.RFC3339Nano),
		)
	}
}

func callChat(client greeterpb.GreeterClient, messages []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.Chat(ctx)
	if err != nil {
		log.Fatalf("Error creating chat stream: %v", err)
	}

	waitc := make(chan struct{})

	// Receive messages
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Error receiving message: %v", err)
			}
			log.Printf("Server says: %s at %v", 
				in.GetMessage(), 
				in.GetTimestamp().AsTime().Format(time.RFC3339Nano),
			)
		}
	}()

	// Send messages
	for _, msg := range messages {
		if err := stream.Send(&greeterpb.HelloRequest{Name: msg}); err != nil {
			log.Fatalf("Failed to send message: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}
	stream.CloseSend()
	<-waitc
}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create clients
	greeterClient := greeterpb.NewGreeterClient(conn)
	healthClient := grpc_health_v1.NewHealthClient(conn)

	// Check service health
	fmt.Println("=== Checking service health ===")
	checkHealth(healthClient)

	// Test unary RPC
	fmt.Println("\n=== Testing SayHello (Unary RPC) ===")
	callSayHello(greeterClient, "World")

	// Test server streaming RPC
	fmt.Println("\n=== Testing StreamGreetings (Server Streaming) ===")
	callStreamGreetings(greeterClient, "Streaming Client")

	// Test bidirectional streaming RPC
	fmt.Println("\n=== Testing Chat (Bidirectional Streaming) ===")
	messages := []string{"Hello!", "How are you?", "This is a test", "Goodbye!"}
	callChat(greeterClient, messages)
}
