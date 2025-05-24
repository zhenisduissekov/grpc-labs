package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	
	greeterpb "step-07_reflection_health/internal/greeter"
)

type server struct {
	greeterpb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *greeterpb.HelloRequest) (*greeterpb.HelloReply, error) {
	return &greeterpb.HelloReply{
		Message:   "Hello " + in.Name,
		Timestamp: timestamppb.Now(),
	}, nil
}

func (s *server) StreamGreetings(in *greeterpb.HelloRequest, stream greeterpb.Greeter_StreamGreetingsServer) error {
	for i := 0; i < 5; i++ {
		if err := stream.Send(&greeterpb.HelloReply{
			Message:   "Hello " + in.Name + " #" + fmt.Sprint(i+1),
			Timestamp: timestamppb.Now(),
		}); err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

func (s *server) Chat(stream greeterpb.Greeter_ChatServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			return err
		}

		reply := &greeterpb.HelloReply{
			Message:   "You said: " + in.Name,
			Timestamp: timestamppb.Now(),
		}

		if err := stream.Send(reply); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	
	// Register the Greeter service
	greeterpb.RegisterGreeterServer(s, &server{})
	
	// Register reflection service on gRPC server
	reflection.Register(s)
	
	// Register health check service
	healthServer := health.NewServer()
	healthServer.SetServingStatus("greeter.Greeter", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
