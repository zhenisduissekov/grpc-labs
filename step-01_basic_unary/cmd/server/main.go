package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	greeterpb "step-01_basic_unary/internal/greeter"
)

type server struct {
	greeterpb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *greeterpb.HelloRequest) (*greeterpb.HelloReply, error) {
	return &greeterpb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greeterpb.RegisterGreeterServer(s, &server{})
	reflection.Register(s)

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
