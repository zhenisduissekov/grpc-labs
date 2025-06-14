package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	greeterpb "step-15_grafana_dashboards/internal/greeter"
)

type server struct {
	greeterpb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *greeterpb.HelloRequest) (*greeterpb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	if in.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "Name cannot be empty")
	}

	// Simulate some processing time
	time.Sleep(100 * time.Millisecond)

	// Randomly fail 10% of the time for demo purposes
	if time.Now().UnixNano()%10 == 0 {
		return nil, status.Error(codes.Internal, "Random error occurred")
	}

	return &greeterpb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	// Create a metrics registry.
	reg := prometheus.NewRegistry()

	// Create some standard server metrics.
	grpcMetrics := grpc_prometheus.NewServerMetrics()

	// Register the metrics.
	reg.MustRegister(grpcMetrics)

	// Create a HTTP server for prometheus.
	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		Addr:    ":9092",
	}

	// Start your http server for prometheus.
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Unable to start a http server: %v", err)
		}
	}()

	// Create a gRPC Server with the interceptor.
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
	)

	// Register your service.
	service := &server{}
	greeterpb.RegisterGreeterServer(s, service)

	// Register reflection service on gRPC server.
	reflection.Register(s)

	// Initialize all metrics.
	grpcMetrics.InitializeMetrics(s)

	// Start the gRPC server.
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Server started at :50051")
	log.Println("Metrics available at :9092/metrics")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
