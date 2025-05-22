package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"google.golang.org/grpc"
	pb "step-03_bidirectional_streaming/internal/chat"
)

func main() {
	// Connect to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		errC := conn.Close()
		if errC != nil {
			log.Printf("Failed to close connection: %v", errC)
		}
	}()

	// Create client
	client := pb.NewHelloServiceClient(conn)

	// Create bidirectional stream
	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	// Handle incoming messages in a goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			reply, err := stream.Recv()
			if err != nil {
				if err.Error() != "EOF" {
					log.Printf("Error receiving message: %v", err)
				}
				break
			}
			log.Printf("\nReceived: %s\n", reply.Message)
		}
		fmt.Printf("\nStream closed. Exiting...\n")
	}()

	// Read user input and send messages
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Type messages to send (or 'exit' to quit):")

	for scanner.Scan() {
		message := scanner.Text()
		if message == "exit" {
			break
		}

		if err := stream.Send(&pb.HelloRequest{
			Message: message,
		}); err != nil {
			log.Fatalf("Failed to send message: %v", err)
		}
	}

	// Close the stream
	if err := stream.CloseSend(); err != nil {
		log.Fatalf("Failed to close stream: %v", err)
	}

	wg.Wait()
}
