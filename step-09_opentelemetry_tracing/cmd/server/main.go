package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	greeterpb "step-09_opentelemetry_tracing/internal/greeter"
)

type server struct {
	greeterpb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *greeterpb.HelloRequest) (*greeterpb.HelloReply, error) {
	// Get the current span from the context
	_, span := otel.Tracer("greeter").Start(ctx, "SayHello")
	defer span.End()

	// Add attributes to the span
	span.SetAttributes(attribute.String("request.name", in.Name))

	// Simulate some work
	time.Sleep(100 * time.Millisecond)

	return &greeterpb.HelloReply{Message: "Hello " + in.Name}, nil
}

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the OTLP HTTP exporter that will send spans to the provided url.
func tracerProvider(serviceName string) (*sdktrace.TracerProvider, error) {
	// Create the OTLP HTTP exporter
	exp, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint("localhost:14318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", "demo"),
		)),
	)

	return tp, nil
}

func main() {
	// Initialize tracer provider
	tp, err := tracerProvider("greeter-service")
	if err != nil {
		log.Fatalf("Failed to create tracer provider: %v", err)
	}

	// Register our TracerProvider as the global
	otel.SetTracerProvider(tp)

	// Set up context propagation
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Set up gRPC server with OpenTelemetry interceptors
	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	// Register the Greeter server
	greeterpb.RegisterGreeterServer(s, &server{})

	// Enable reflection for tools like grpcurl
	reflection.Register(s)

	// Start listening on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Server listening at %v", lis.Addr())

	// Set up channel to handle shutdown signals
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stopCh
	log.Println("Shutting down server...")

	// Gracefully stop the server
	s.GracefulStop()

	// Shut down the tracer provider
	if err := tp.Shutdown(context.Background()); err != nil {
		log.Fatalf("Error shutting down tracer provider: %v", err)
	}

	log.Println("Server stopped")
}
