package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	greeterpb "step-13_retry_timeout/internal/greeter"
)

const (
	defaultName    = "world"
	defaultServer  = "localhost:50051"
	initialBackoff = 100 * time.Millisecond
	maxBackoff     = 5 * time.Second
	jitter         = 0.2
	maxRetries     = 3
	defaultTimeout = 3 * time.Second
)

func main() {
	// Parse command line flags
	name := flag.String("name", defaultName, "Name to greet")
	serverAddr := flag.String("server", defaultServer, "Server address (host:port)")
	simulateError := flag.Bool("error", false, "Simulate server error")
	delayMs := flag.Int("delay", 0, "Simulate server delay in milliseconds")
	timeout := flag.Duration("timeout", defaultTimeout, "Request timeout")
	flag.Parse()

	// Set up a connection to the server
	conn, err := grpc.Dial(*serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	client := greeterpb.NewGreeterClient(conn)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// Create the request
	req := &greeterpb.HelloRequest{
		Name:          *name,
		SimulateError: *simulateError,
		DelayMs:       int32(*delayMs),
	}

	// Call the server with retry logic
	var resp *greeterpb.HelloReply
	var attempt int
	backoff := initialBackoff

	for attempt = 1; attempt <= maxRetries; attempt++ {
		log.Printf("Attempt %d/%d", attempt, maxRetries)

		// Create a new context for each attempt
		attemptCtx, attemptCancel := context.WithTimeout(ctx, *timeout)

		// Make the RPC call
		start := time.Now()
		resp, err = client.SayHello(attemptCtx, req)
		elapsed := time.Since(start)

		// Clean up the attempt context
		attemptCancel()

		// Check if the call was successful
		if err == nil {
			log.Printf("Success after %d attempts (took %v)", attempt, elapsed)
			fmt.Printf("Response: %s\n", resp.Message)
			return
		}

		// Check if we should retry
		if !shouldRetry(err) {
			log.Printf("Non-retryable error: %v", err)
			break
		}

		// Log the error and wait before retrying
		log.Printf("Attempt %d failed: %v (took %v)", attempt, err, elapsed)

		// Calculate backoff with jitter
		sleepTime := backoff + time.Duration(float64(backoff)*jitter*(2*randomFloat64()-1))
		if sleepTime > maxBackoff {
			sleepTime = maxBackoff
		}

		log.Printf("Waiting %v before retry...", sleepTime)

		// Wait for backoff or until context is done
		select {
		case <-time.After(sleepTime):
		case <-ctx.Done():
			log.Printf("Context done: %v", ctx.Err())
			break
		}

		// Exponential backoff for next attempt
		backoff *= 2
	}

	// If we get here, all retries failed
	log.Fatalf("Failed after %d attempts: %v", attempt-1, err)
}

// shouldRetry determines if an error is retryable
func shouldRetry(err error) bool {
	if err == nil {
		return false
	}

	// Convert to gRPC status
	s, ok := status.FromError(err)
	if !ok {
		// Not a gRPC error, don't retry
		return false
	}

	// Only retry on certain status codes
	switch s.Code() {
	case codes.DeadlineExceeded, codes.Unavailable, codes.ResourceExhausted, codes.Aborted, codes.Internal:
		return true
	default:
		return false
	}
}

// randomFloat64 returns a random float64 between 0 and 1
func randomFloat64() float64 {
	return float64(time.Now().UnixNano()%1000) / 1000.0
}
