package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	pb "step-06_tls_encryption/internal/greeter"
)

const (
	port     = ":50051"
	certFile = "certs/server.crt"
	keyFile  = "certs/server.key"
	caFile   = "certs/ca.crt"
)

// server is used to implement greeter.GreeterServer
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements greeter.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load server key pair: %v", err)
	}

	// Load certificate of the CA who signed client's certificate
	pemClientCA, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load client CA: %v", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	// Create the TLS credentials
	creds, err := loadTLSCredentials()
	if err != nil {
		log.Fatalf("could not load TLS keys: %s", err)
	}

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create an array of gRPC server options with the credentials
	s := grpc.NewServer(grpc.Creds(creds))

	// Register the Greeter service on the server
	pb.RegisterGreeterServer(s, &server{})

	// Enable server reflection
	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())

	// Start the server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
