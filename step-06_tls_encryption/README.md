# gRPC with TLS Encryption

This project demonstrates how to secure gRPC communication using TLS with self-signed certificates. It includes scripts to generate the necessary certificates and shows how to configure both server and client to use TLS.

## Project Structure

```
.
├── certs/            # Directory containing TLS certificates
├── cmd/              # Command-line applications
│   ├── client/      # gRPC client implementation with TLS
│   └── server/      # gRPC server implementation with TLS
├── internal/         # Internal packages
│   └── greeter/     # Generated protobuf code
├── proto/           # Protocol buffer definitions
│   └── greeter.proto
├── go.mod           # Go module configuration
├── Makefile         # Build automation
└── README.md        # This file
```

## Features

- Secure gRPC communication using TLS 1.3
- Self-signed certificate generation script
- Both server and client authentication
- Secure by default configuration

## Setup and Usage

1. Generate TLS certificates (requires OpenSSL):
   ```bash
   make certs
   ```
   This will create the following files in the `certs` directory:
   - `ca.crt` - Certificate Authority certificate
   - `server.crt` - Server certificate
   - `server.key` - Server private key
   - `client.crt` - Client certificate (for mTLS, if needed)
   - `client.key` - Client private key (for mTLS, if needed)

2. Initialize the project and generate code:
   ```bash
   make init
   make generate
   ```

3. Run the server with TLS:
   ```bash
   make run-server
   ```

4. In a separate terminal, run the client with TLS:
   ```bash
   make run-client
   ```

## Certificate Generation

The `make certs` command generates:
- A self-signed Certificate Authority (CA)
- A server certificate signed by the CA
- A client certificate signed by the CA (for mTLS)
- All necessary private keys

## Security Notes

- In production, use certificates from a trusted Certificate Authority (CA)
- Keep private keys secure and never commit them to version control
- The generated certificates use a 1-year expiration for demonstration purposes
- The server's hostname is set to `localhost` in the certificate's Subject Alternative Name (SAN)

## Verifying the Connection

You can verify the TLS connection using `openssl`:

```bash
echo | openssl s_client -connect localhost:50051 -showcerts
```

Or using `grpcurl`:

```bash
grpcurl -insecure localhost:50051 list
grpcurl -insecure -d '{"name": "World"}' localhost:50051 greeter.Greeter/SayHello
```

## Cleaning Up

To remove generated certificates and clean up:

```bash
make clean-certs
```
