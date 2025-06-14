package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	greeterpb "step-15_grafana_dashboards/internal/greeter"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_client_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_client_request_duration_seconds",
			Help:    "Duration of gRPC requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)
)

func init() {
	// Register metrics with the global prometheus registry
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestDuration)
}

func main() {
	// Start HTTP server for Prometheus metrics
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Starting metrics server on :9093/metrics")
		log.Fatal(http.ListenAndServe(":9093", nil))
	}()

	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := greeterpb.NewGreeterClient(conn)

	// Channel to listen for interrupt signal to terminate.
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Run requests in a loop until interrupted
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	names := []string{"Alice", "Bob", "Charlie", "David", "Eve", ""}
	nameIndex := 0

	for {
		select {
		case <-done:
			log.Println("Shutting down client...")
			return
		case <-ticker.C:
			name := names[nameIndex%len(names)]
			nameIndex++

			// Record start time
			start := time.Now()

			// Make the gRPC call
			_, err := c.SayHello(context.Background(), &greeterpb.HelloRequest{Name: name})

			// Calculate duration
			duration := time.Since(start).Seconds()

			// Record metrics
			status := "success"
			if err != nil {
				status = "error"
			}

			requestCounter.WithLabelValues("SayHello", status).Inc()
			requestDuration.WithLabelValues("SayHello", status).Observe(duration)

			if err != nil {
				log.Printf("Error calling SayHello: %v", err)
			} else {
				log.Printf("Request sent for name: %s (took %.3fs)", name, duration)
			}
		}
	}
}
