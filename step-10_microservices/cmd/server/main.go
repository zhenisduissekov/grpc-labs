package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	loggerpb "step-10_microservices/internal/logger"
	serverpb "step-10_microservices/internal/server"
)

type server struct {
	serverpb.UnimplementedServerServer
	loggerClient loggerpb.LoggerClient
}

func (s *server) SayHello(ctx context.Context, req *serverpb.HelloRequest) (*serverpb.HelloReply, error) {
	// Debug log the incoming request
	log.Printf("Received SayHello request with name: %q", req.GetName())
	
	// Log the request using the Logger service
	_, err := s.loggerClient.Log(context.Background(), &loggerpb.LogRequest{
		Message: "Received hello request for: " + req.GetName(),
		Service: "server",
		Level:   "INFO",
	})
	if err != nil {
		log.Printf("Failed to log: %v", err)
	}
	log.Printf("Received message from client: %s", req.GetName())

	return &serverpb.HelloReply{
		Message: "Hello, " + req.GetName() + "!",
	}, nil
}

func main() {
	// Set up a connection to the Logger service
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect to logger: %v", err)
	}
	defer conn.Close()

	// Create server instance with logger client
	srv := &server{
		loggerClient: loggerpb.NewLoggerClient(conn),
	}

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	serverpb.RegisterServerServer(s, srv)

	// Enable reflection for testing with grpcurl
	reflection.Register(s)

	log.Println("Server service listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
