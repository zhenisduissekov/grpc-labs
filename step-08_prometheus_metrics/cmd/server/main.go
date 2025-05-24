package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	greeterpb "step-08_prometheus_metrics/internal/greeter"
)

type server struct {
	greeterpb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *greeterpb.HelloRequest) (*greeterpb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &greeterpb.HelloReply{
		Message:   fmt.Sprintf("Hello %s", in.GetName()),
		Timestamp: timestamppb.Now(),
	}, nil
}

func (s *server) StreamGreetings(in *greeterpb.HelloRequest, stream greeterpb.Greeter_StreamGreetingsServer) error {
	for i := 0; i < 3; i++ {
		if err := stream.Send(&greeterpb.HelloReply{
			Message:   fmt.Sprintf("Hello %s #%d", in.GetName(), i+1),
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
			Message:   "You said: " + in.GetName(),
			Timestamp: timestamppb.Now(),
		}

		if err := stream.Send(reply); err != nil {
			return err
		}
	}
}

func main() {
	// Create gRPC server with Prometheus interceptors
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)

	// Register service
	greeterpb.RegisterGreeterServer(s, &server{})

	// Enable reflection for debugging
	reflection.Register(s)

	// Enable Prometheus metrics
	grpc_prometheus.EnableHandlingTimeHistogram(
		grpc_prometheus.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
	)
	grpc_prometheus.Register(s)

	// Start metrics server in a separate goroutine
	go func() {
		metricsMux := http.NewServeMux()
		metricsMux.Handle("/metrics", promhttp.Handler())
		metricsServer := &http.Server{
			Addr:    ":9090",
			Handler: metricsMux,
		}

		log.Printf("Starting metrics server on http://localhost:9090/metrics")
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
