package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"time"

	"step-04_interceptors/internal/greeter"
	loggerpb "step-04_interceptors/internal/logger"
)

type greeterServer struct {
	greeterpb.UnimplementedGreeterServer
	loggerClient loggerpb.LoggerClient
}

func (s *greeterServer) SayHello(ctx context.Context, req *greeterpb.HelloRequest) (*greeterpb.HelloReply, error) {
	name := req.GetName()
	message := "Hello, " + name + "!"

	// Call logger service
	_, err := s.loggerClient.Log(ctx, &loggerpb.LogRequest{Message: message})
	if err != nil {
		log.Printf("❌ failed to log: %v", err)
	}

	return &greeterpb.HelloReply{Message: message}, nil
}

func (s *greeterServer) StreamGreetings(req *greeterpb.HelloRequest, stream greeterpb.Greeter_StreamGreetingsServer) error {
	name := req.GetName()

	for i := 1; i <= 5; i++ {
		msg := &greeterpb.HelloReply{
			Message: fmt.Sprintf("Hello %s #%d", name, i),
		}

		if err := stream.Send(msg); err != nil {
			return err
		}

		time.Sleep(1 * time.Second) // simulate delay
	}

	return nil
}

func (s *greeterServer) Chat(stream greeterpb.Greeter_ChatServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil // client closed stream
		}
		if err != nil {
			return err
		}

		// Send response
		resp := &greeterpb.HelloReply{
			Message: fmt.Sprintf("👋 Hello, %s!", req.GetName()),
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}

func (s *greeterServer) UploadNames(stream greeterpb.Greeter_UploadNamesServer) error {
	var names []string

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// client finished sending
			reply := &greeterpb.HelloReply{
				Message: fmt.Sprintf("✅ Received %d names: %s", len(names), names),
			}
			return stream.SendAndClose(reply)
		}
		if err != nil {
			return err
		}
		names = append(names, req.GetName())
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggingUnaryInterceptor,
		),
		grpc.ChainStreamInterceptor(
			loggingStreamInterceptor,
		),
	)

	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("❌ Could not connect to logger: %v", err)
	}
	defer conn.Close()

	loggerClient := loggerpb.NewLoggerClient(conn)

	// Register your service first
	greeterpb.RegisterGreeterServer(grpcServer, &greeterServer{loggerClient: loggerClient})

	// Register health check service
	healthServer := health.NewServer()

	// This marks the Greeter service as healthy
	healthServer.SetServingStatus("greeter.Greeter", healthpb.HealthCheckResponse_SERVING)

	// Register the health service with gRPC
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	// Enable reflection
	reflection.Register(grpcServer)

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func loggingUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)

	log.Printf("▶️ Unary call: %s | Duration: %s | Error: %v", info.FullMethod, duration, err)
	return resp, err
}

func loggingStreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	start := time.Now()
	err := handler(srv, ss)
	duration := time.Since(start)

	log.Printf("🔁 Stream call: %s | Duration: %s | Error: %v", info.FullMethod, duration, err)
	return err
}
