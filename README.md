
# 🧠 gRPC Mastery: Step-by-Step Guide (Steps 1–17)

This README documents each advanced gRPC feature implemented step-by-step, ideal for learning and evolving production-grade microservices in Go.


---

## ✅ Step 1: Basic Unary RPC – `SayHello`

- A simple request/response method.
- Sends `HelloRequest`, receives `HelloReply`.

```proto
rpc SayHello(HelloRequest) returns (HelloReply);
```

⸻

## ✅ Step 2: Server Streaming – StreamGreetings
- The client sends one request.
- The server responds with a stream of `HelloReply`.

```proto
rpc StreamGreetings(HelloRequest) returns (stream HelloReply);
```

Use case: live updates, chat history, or data feeds.

⸻

✅ Step 3: Bidirectional Streaming – Chat
- Both the client and server stream messages.
- Enables real-time chat or event-driven applications.

```proto
rpc Chat(stream HelloRequest) returns (stream HelloReply);
```

⸻

✅ Step 4: Interceptors (Middleware)
- Add logic around RPCs: logging, auth, tracing, etc.

```go
grpc.NewServer(
  grpc.UnaryInterceptor(loggingUnaryInterceptor),
)
```

Use cases:
- 🔐 Auth
- 📝 Logging
- 📊 Metrics
- 🔁 Retry

⸻

✅ Step 5: Metadata-Based Auth (Token)
	•	Send headers from client (e.g. authorization, x-user-id)
	•	Read headers in interceptor using:

metadata.FromIncomingContext(ctx)

Reject unauthorized requests with:

status.Error(codes.Unauthenticated, "invalid token")

⸻

✅ Step 6: TLS with Self-Signed Certs
	•	Secure the gRPC channel using TLS
	•	Generate certs with SAN using openssl
	•	Server uses:

credentials.NewServerTLSFromFile("cert/server.crt", "cert/server.key")

Client uses:

credentials.NewClientTLSFromFile("cert/server.crt", "")

⸻

✅ Step 7: Reflection + Health Checking
	•	Enable server reflection:

reflection.Register(grpcServer)

	•	Add health check endpoint:

healthpb.RegisterHealthServer(grpcServer, health.NewServer())

This enables tools like [`grpcurl`](https://github.com/fullstorydev/grpcurl), Kubernetes, and Envoy to monitor the service state.

⸻

✅ Step 8: Prometheus Metrics
	•	Use go-grpc-prometheus to expose metrics like:
	•	RPC count
	•	Latency
	•	Errors

grpc_prometheus.Register(grpcServer)
http.Handle("/metrics", promhttp.Handler())

Scrape from Prometheus at :9090/metrics.

⸻

✅ Step 9: OpenTelemetry Tracing with Jaeger
	•	Set up tracer with otel
	•	Export spans to Jaeger:

go.opentelemetry.io/otel/exporters/jaeger

	•	Wrap gRPC with interceptors:

otelgrpc.UnaryServerInterceptor()

View traces in Jaeger UI (http://localhost:16686).

⸻

✅ Step 10: Multi-Service (gRPC → gRPC)
	•	GreeterService calls LoggerService via gRPC
	•	Both are independent microservices
	•	Simulates microservice communication

Use grpc.Dial() to connect inside Greeter and send a request.

⸻

✅ Step 11: Metadata Propagation
	•	Pass metadata (e.g. x-user-id) across services:

metadata.AppendToOutgoingContext(ctx, "x-user-id", "zhenis123")

	•	LoggerService reads it from incoming context:

metadata.FromIncomingContext(ctx)

⸻

✅ Step 12: Load Balancing + Service Discovery
	•	Use gRPC’s built-in round_robin LB:

grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`)

	•	Use DNS (e.g. dns:///host1:60051,host2:60052) or Kubernetes service discovery.

⸻

✅ Step 13: Retry + Timeout (Optional Extension)
	•	Use context timeouts:

ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	•	Add retry logic around calls (e.g., with backoff or heimdall)

⸻

✅ Step 14: Circuit Breaker (Resilience)
	•	Use github.com/sony/gobreaker
	•	Wrap gRPC calls:

breaker.Execute(func() (interface{}, error) {
  return client.Call(...)
})

Avoids cascading failures by stopping bad calls.

⸻

✅ Step 15: Grafana Dashboards
	•	Visualize:
	•	gRPC latency, count, error rate (Prometheus)
	•	Traces via Jaeger or Tempo
	•	Build dashboards with:
	•	Panels for grpc_server_handled_total
	•	Heatmaps for latency buckets

⸻

✅ Step 16: Async gRPC + Kafka/NATS Integration
	•	GreeterService returns immediately but:
	•	Sends Kafka message to events.greeted
	•	LoggerService or analytics consume async

Use the [`segmentio/kafka-go`](https://github.com/segmentio/kafka-go) library to decouple services and improve scalability.

⸻

✅ Step 17: Versioned APIs + Canary Deployments
	•	Use versioned packages:

package greeter.v1;
package greeter.v2;

	•	Run GreeterV1 and GreeterV2 side by side
	•	Use Envoy/Istio to route traffic:

weighted_clusters:
  - name: v1  weight: 90
  - name: v2  weight: 10

Safe rollout of breaking changes.

⸻

👋 Next Steps?
- ✅ Step 18: Rate Limiting per user/service
- ✅ Step 19: gRPC-Gateway (REST ↔ gRPC)
- ✅ Step 20: CI/CD for gRPC with Docker + K8s
- ✅ Step 21: Request validation (with proto annotations)

Feel free to explore these additional steps to enhance your gRPC services.

⸻

✅ Who is this for?
	•	👨‍💻 Go developers building microservices
	•	🧪 Backend engineers working with gRPC APIs
	•	📈 DevOps / Platform teams adding observability and resilience
	•	🧠 Anyone who wants to build, secure, scale, and monitor gRPC services like a pro

⸻

🚀 Getting Started

git clone https://github.com/zhenisduissekov/grpc-labs.git
cd grpc-labs/step-01_basic_unary
# Run the following commands as described in the local README.md:
# make init
# make generate
# make run-server
cd ../step-02_server_streaming
# Refer to the local README.md for detailed instructions, including `make init`, `make generate`, and running the service.
...

---
### 🛠️ Built with Help from AI

This project was developed while learning and exploring gRPC and Go, with the assistance of AI tools like [OpenAI's ChatGPT](https://openai.com/chatgpt) and [Windsurf](https://windsurf.ai). All code has been reviewed and customized to ensure clarity and educational value.
