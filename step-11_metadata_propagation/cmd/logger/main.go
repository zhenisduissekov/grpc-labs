package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	loggerpb "step-11_metadata_propagation/internal/logger"
)

type server struct {
	loggerpb.UnimplementedLoggerServer
}

func (s *server) Log(ctx context.Context, req *loggerpb.LogRequest) (*loggerpb.LogResponse, error) {
	log.Printf("[%s] %s: %s", req.GetLevel(), req.GetService(), req.GetMessage())
	return &loggerpb.LogResponse{
		Success:   true,
		MessageId: "log-" + req.GetService() + "-" + req.GetLevel(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	loggerpb.RegisterLoggerServer(s, &server{})

	// Enable reflection for testing with grpcurl
	reflection.Register(s)

	log.Println("Logger service listening on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
