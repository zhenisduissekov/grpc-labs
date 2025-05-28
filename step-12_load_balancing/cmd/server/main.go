package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	greeterpb "step-12_load_balancing/internal/greeter"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	greeterpb.UnimplementedGreeterServer
	port string
}

func (s *server) SayHello(ctx context.Context, in *greeterpb.HelloRequest) (*greeterpb.HelloReply, error) {
	log.Printf("Received request from client for name: %s", in.Name)
	resp := &greeterpb.HelloReply{
		Message: fmt.Sprintf("Hello %s (from server on port %s)", in.Name, s.port),
	}
	log.Printf("Sending response: %s", resp.Message)
	return resp, nil
}

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	// Create listener
	addr := ":" + port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	srv := &server{port: port}
	log.Printf("Server started on %s", addr)

	s := grpc.NewServer()
	greeterpb.RegisterGreeterServer(s, srv)
	reflection.Register(s)

	log.Printf("Server %s listening at %v", os.Args[0], lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
