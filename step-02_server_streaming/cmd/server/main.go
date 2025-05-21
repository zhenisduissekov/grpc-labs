package main

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"

	greeter "step-02_server_streaming/internal/greeter"

	"google.golang.org/grpc"
)

type server struct {
	greeter.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *greeter.HelloRequest) (*greeter.HelloReply, error) {
	return &greeter.HelloReply{Message: "Hello " + req.Name}, nil
}

func (s *server) StreamGreetings(req *greeter.HelloRequest, stream greeter.Greeter_StreamGreetingsServer) error {
	name := req.Name
	for i := 0; i < 5; i++ {
		if err := stream.Send(&greeter.HelloReply{
			Message: "Hello " + name + " (" + strconv.Itoa(i+1) + "/5)",
		}); err != nil {
			return err
		}
		time.Sleep(1 * time.Second) // Simulate some processing time
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	greeter.RegisterGreeterServer(grpcServer, &server{})
	log.Printf("Server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
