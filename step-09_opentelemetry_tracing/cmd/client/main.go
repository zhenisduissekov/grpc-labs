package main

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	greeterpb "step-09_opentelemetry_tracing/internal/greeter"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

const (
	address = "localhost:50051"
)

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

func callSayHello(client greeterpb.GreeterClient, name string) {
	// Create a root span
	ctx, span := otel.Tracer("greeter-client").Start(context.Background(), "callSayHello")
	defer span.End()

	// Add attributes to the span
	span.SetAttributes(attribute.String("client.request.name", name))

	// Call the SayHello RPC
	r, err := client.SayHello(ctx, &greeterpb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}

func main() {
	// Initialize tracer provider
	tp, err := tracerProvider("greeter-client")
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

	// Set up a connection to the server with OpenTelemetry interceptors
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := greeterpb.NewGreeterClient(conn)

	// Make some RPC calls
	callSayHello(client, "world")
	callSayHello(client, "gRPC")
	callSayHello(client, "OpenTelemetry")

	// Give some time for spans to be exported
	time.Sleep(2 * time.Second)

	// Shut down the tracer provider
	if err := tp.Shutdown(context.Background()); err != nil {
		log.Fatalf("Error shutting down tracer provider: %v", err)
	}
}
