package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	pb "step-03_bidirectional_streaming/internal/chat"
	"sync"
)

type helloServer struct {
	pb.UnimplementedHelloServiceServer
	clients map[pb.HelloService_ChatServer]bool
	mu      sync.Mutex
}

func newHelloServer() *helloServer {
	return &helloServer{
		clients: make(map[pb.HelloService_ChatServer]bool),
	}
}

func (s *helloServer) Chat(stream pb.HelloService_ChatServer) error {
	// Register the client
	s.mu.Lock()
	s.clients[stream] = true
	log.Printf("Client connected. Total clients: %d", len(s.clients))
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		delete(s.clients, stream)
		log.Printf("Client disconnected. Total clients: %d", len(s.clients))
		s.mu.Unlock()
	}()

	for {
		// Read incoming message from client
		req, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return err
		}

		// Log received message
		log.Printf("Received message: %s", req.Message)

		// Broadcast message to all connected clients
		message := "[Server] " + req.Message

		// Send message to all clients except the sender
		s.mu.Lock()
		for client := range s.clients {
			if client != stream {
				if err := client.Send(&pb.HelloReply{
					Message: message,
				}); err != nil {
					log.Printf("Failed to send message to client: %v", err)
				} else {
					log.Printf("Message sent to client: %s", message)
				}
			}
		}
		s.mu.Unlock()
	}
}

func main() {
	// Create TCP listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register our service
	pb.RegisterHelloServiceServer(grpcServer, newHelloServer())

	// Start serving
	log.Printf("Server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
