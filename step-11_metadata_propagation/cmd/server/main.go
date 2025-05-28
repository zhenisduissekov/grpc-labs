package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	pb "step-11_metadata_propagation/internal/greeter"
	loggerpb "step-11_metadata_propagation/internal/logger"
)

type server struct {
	pb.UnimplementedGreeterServer
	loggerClient loggerpb.LoggerClient
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	// Extract incoming metadata
	md, _ := metadata.FromIncomingContext(ctx)
	log.Printf("Received request with metadata: %v", md)

	// Create a new context with the incoming metadata
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Call Logger service (in a real app, this would be a separate service)
	log.Printf("Forwarding metadata to Logger service: %v", md)

	// Send metadata to Logger service
	_, err := s.loggerClient.Log(ctx, &loggerpb.LogRequest{
		Message: "Forwarded metadata",
		Service: "GreeterService",
		Level:   "INFO",
	})
	if err != nil {
		log.Printf("Failed to log metadata: %v", err)
	}

	return &pb.HelloReply{
		Message: "Hello " + req.Name + ", your metadata has been processed",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a connection to the Logger service
	conn, err := grpc.Dial(":50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to Logger service: %v", err)
	}
	defer conn.Close()

	loggerClient := loggerpb.NewLoggerClient(conn)

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{loggerClient: loggerClient})

	// Register reflection service on gRPC server
	reflection.Register(s)

	log.Println("Server started on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
