.PHONY: all init generate certs clean-certs clean run-server run-client

# Go parameters
GOCMD = go
GOBIN = $(shell go env GOPATH)/bin
PROTOC_GEN_GO = $(GOBIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC = $(GOBIN)/protoc-gen-go-grpc
CERT_DIR = certs

# Initialize the project
init:
	rm -rf go.mod go.sum
	go mod init step-06_tls_encryption
	go mod tidy
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate protobuf code
generate:
	mkdir -p internal/greeter
	PATH="$(shell go env GOPATH)/bin:$$PATH" protoc --go_out=. --go-grpc_out=. proto/greeter.proto
	go mod tidy

# Generate TLS certificates
certs:
	mkdir -p $(CERT_DIR)
	# Generate CA private key and certificate
	openssl req -x509 -newkey rsa:4096 -days 365 -nodes \
	  -keyout $(CERT_DIR)/ca.key -out $(CERT_DIR)/ca.crt \
	  -subj "/CN=Test CA" -addext "basicConstraints=critical,CA:TRUE"
	
	echo "[req]\ndistinguished_name=req_distinguished_name\nreq_extensions=req_ext\n[req_distinguished_name]\n[req_ext]\nsubjectAltName=@alt_names\n[alt_names]\nDNS.1=localhost\nIP.1=127.0.0.1" > $(CERT_DIR)/server-ext.cnf

	# Generate server private key and certificate signing request (CSR)
	openssl req -newkey rsa:4096 -nodes \
	  -keyout $(CERT_DIR)/server.key -out $(CERT_DIR)/server.csr \
	  -subj "/CN=localhost" \
	  -config $(CERT_DIR)/server-ext.cnf
	
	# Sign server certificate with CA
	openssl x509 -req -in $(CERT_DIR)/server.csr \
	  -CA $(CERT_DIR)/ca.crt -CAkey $(CERT_DIR)/ca.key -CAcreateserial \
	  -out $(CERT_DIR)/server.crt -days 365 \
	  -extfile $(CERT_DIR)/server-ext.cnf -extensions req_ext
	
	# Generate client private key and certificate signing request (CSR)
	openssl req -newkey rsa:4096 -nodes \
	  -keyout $(CERT_DIR)/client.key -out $(CERT_DIR)/client.csr \
	  -subj "/CN=client"
	
	# Sign client certificate with CA
	openssl x509 -req -in $(CERT_DIR)/client.csr \
	  -CA $(CERT_DIR)/ca.crt -CAkey $(CERT_DIR)/ca.key -CAcreateserial \
	  -out $(CERT_DIR)/client.crt -days 365
	
	# Clean up certificate signing requests and CA serial
	rm -f $(CERT_DIR)/*.csr $(CERT_DIR)/*.srl $(CERT_DIR)/*.cnf
	
	@echo "✅ Certificates generated in $(CERT_DIR)/"

# Clean generated files
clean:
	rm -f internal/greeter/*.pb.go

# Clean certificates
clean-certs:
	rm -rf $(CERT_DIR)

# Run the server
run-server:
	@echo "Starting gRPC server with TLS..."
	@go run cmd/server/main.go

# Run the client
run-client:
	@echo "Running gRPC client with TLS..."
	@go run cmd/client/main.go
