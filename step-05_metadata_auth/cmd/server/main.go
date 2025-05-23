package main

import (
	"context"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "step-05_metadata_auth/internal/greeter"
)

const (
	port = ":50051"
	token = "my-secret-token"
)

// server is used to implement greeter.GreeterServer
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements unary RPC without authentication
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

// SecureGreeting implements unary RPC with token authentication
func (s *server) SecureGreeting(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	// Extract metadata from context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	// Get authorization header
	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	// Validate token
	if !strings.HasPrefix(authHeader[0], "bearer ") || 
	   !isValidToken(strings.TrimPrefix(authHeader[0], "bearer ")) {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	log.Printf("Secure greeting for: %v", in.GetName())
	return &pb.HelloReply{Message: "Secure hello " + in.GetName()}, nil
}

// isValidToken validates the provided token
func isValidToken(token string) bool {
	// In a real application, you would validate against a database or auth service
	// This is a simplified example
	return token == "my-secret-token"
}

// unaryInterceptor is a server interceptor for authentication
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Skip auth for SayHello method
	if info.FullMethod == "/greeter.Greeter/SayHello" {
		return handler(ctx, req)
	}

	// For other methods, require authentication
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	if !strings.HasPrefix(authHeader[0], "bearer ") || 
	   !isValidToken(strings.TrimPrefix(authHeader[0], "bearer ")) {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	// Call the handler if token is valid
	return handler(ctx, req)
}

func main() {
	// Create a TCP listener
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a gRPC server with the interceptor
	s := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor),
	)

	// Register the Greeter service
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())

	// Start serving
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
