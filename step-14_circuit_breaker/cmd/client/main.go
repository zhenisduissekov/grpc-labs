package main

import (
	"context"
	"log"
	"time"

	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	greeterpb "step-14_circuit_breaker/internal/greeter"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := greeterpb.NewGreeterClient(conn)

	// Configure circuit breaker with more sensitive settings for demo
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "greeter",
		MaxRequests: 1,                // Only 1 request allowed in half-open state
		Interval:    10 * time.Second, // Reset counts after this time
		Timeout:     5 * time.Second,  // Time in open state
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Trip after 2 failures for demo purposes
			return counts.ConsecutiveFailures >= 2
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			log.Printf("Circuit Breaker '%s' changed from %s to %s\n", name, from, to)
		},
	})

	// First, send enough errors to trip the circuit breaker
	for i := 0; i < 5; i++ {
		log.Printf("\n--- Sending error request %d ---", i+1)
		_, err := cb.Execute(func() (interface{}, error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			_, err := client.SayHello(ctx, &greeterpb.HelloRequest{Name: ""})
			if err != nil {
				log.Printf("gRPC call failed: %v", err)
			}
			return nil, err
		})

		if err != nil {
			log.Printf("Error in circuit breaker execution: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Now send a mix of good and bad requests to see the circuit breaker in action
	for i := 0; i < 10; i++ {
		message := "World"
		if i%2 == 0 {
			message = "" // This will cause an error in the server
		}

		log.Printf("\n--- Making request %d with name: '%s' ---", i+1, message)

		// Wrap the gRPC call with the circuit breaker
		result, err := cb.Execute(func() (interface{}, error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			r, err := client.SayHello(ctx, &greeterpb.HelloRequest{Name: message})
			if err != nil {
				log.Printf("gRPC call failed: %v", err)
				return nil, err
			}
			return r, nil
		})

		if err != nil {
			if err == gobreaker.ErrOpenState {
				log.Printf("Circuit breaker is open, request blocked")
			} else {
				log.Printf("Error in circuit breaker execution: %v", err)
			}
		} else {
			reply, ok := result.(*greeterpb.HelloReply)
			if ok {
				log.Printf("Success: %s", reply.GetMessage())
			}
		}

		time.Sleep(1 * time.Second)
	}
}
