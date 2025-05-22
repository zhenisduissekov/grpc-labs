package main

import (
	"context"
	"log"
	"net"

	loggerpb "step-04_interceptors/internal/logger"

	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

type loggerServer struct {
	loggerpb.UnimplementedLoggerServer
}

func (s *loggerServer) Log(ctx context.Context, req *loggerpb.LogRequest) (*loggerpb.LogReply, error) {
	log.Printf("Received log message: %s", req.Message)
	return &loggerpb.LogReply{
		Ok: true,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	loggerpb.RegisterLoggerServer(grpcServer, &loggerServer{})

	log.Printf("Logger service starting on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
