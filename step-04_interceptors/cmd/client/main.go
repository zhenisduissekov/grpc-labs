package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc/metadata"

	greeterpb "step-04_interceptors/internal/greeter"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer func() {
		errC := conn.Close()
		if errC != nil {
			log.Printf("Failed to close connection: %v", errC)
		}
	}()

	client := greeterpb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx,
		"authorization", "secret-token-123",
		"x-user-id", "5f11a5a9-1b9e-4a07-8c7e-9b8c7a2b2cd5",
	)
	for _, name := range []string{"Zhenis", "John", "Doe"} {
		resp, err := client.SayHello(ctx, &greeterpb.HelloRequest{Name: name})
		if err != nil {
			log.Fatalf("error calling SayHello: %v", err)
		}
		log.Printf("Server responded: %s", resp.GetMessage())
	}

	log.Println("gRPC client finished successfully")
}
