package main

import (
	"context"
	"log"

	greeterpb "step-02_server_streaming/internal/greeter"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := greeterpb.NewGreeterClient(conn)

	// Unary call
	unaryResponse, err := c.SayHello(context.Background(), &greeterpb.HelloRequest{Name: "World"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Unary response: %s", unaryResponse.Message)

	// Server streaming call
	stream, err := c.StreamGreetings(context.Background(), &greeterpb.HelloRequest{Name: "World"})
	if err != nil {
		log.Fatalf("could not stream: %v", err)
	}

	for {
		response, err := stream.Recv()
		if err != nil {
			log.Printf("Stream ended: %v", err)
			break
		}
		log.Printf("Stream response: %s", response.Message)
	}
}
