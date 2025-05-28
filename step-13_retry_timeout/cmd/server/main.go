package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	greeterpb "step-13_retry_timeout/internal/greeter"
)

type server struct {
	greeterpb.UnimplementedGreeterServer
	port string
}

func (s *server) SayHello(ctx context.Context, req *greeterpb.HelloRequest) (*greeterpb.HelloReply, error) {
	log.Printf("Received SayHello request for: %s", req.Name)

	// Simulate processing delay if requested
	if req.DelayMs > 0 {
		time.Sleep(time.Duration(req.DelayMs) * time.Millisecond)
	}

	// Simulate error if requested
	if req.SimulateError {
		log.Printf("Simulating error for request: %s", req.Name)
		return nil, status.Error(codes.Internal, "simulated error as requested")
	}

	// Check if context has been cancelled
	if ctx.Err() == context.Canceled {
		log.Printf("Request cancelled by client: %s", req.Name)
		return nil, status.Error(codes.Canceled, "client cancelled the request")
	}

	// Check if deadline exceeded
	if ctx.Err() == context.DeadlineExceeded {
		log.Printf("Request deadline exceeded for: %s", req.Name)
		return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
	}

	// Return successful response
	return &greeterpb.HelloReply{
		Message:  fmt.Sprintf("Hello %s (from server on port %s)", req.Name, s.port),
		ServerId: s.port,
	}, nil
}

func (s *server) UnaryHello(ctx context.Context, req *greeterpb.HelloRequest) (*greeterpb.HelloReply, error) {
	// For this example, we'll just call SayHello
	return s.SayHello(ctx, req)
}

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}
	log.Printf("Starting server on port %s", port)

	// Create listener
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create gRPC server
	srv := grpc.NewServer()
	greeterpb.RegisterGreeterServer(srv, &server{port: port})

	log.Printf("Server is ready to accept connections on port %s", port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
