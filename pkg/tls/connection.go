package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// CreateCACertificatePool creates a certificate pool from the CA certificate
func CreateCACertificatePool() (*x509.CertPool, error) {
	// TODO: Set domain info in the config file
	const caCertFile = "certs/ca.pem"
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("Could not read ca certificate: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, fmt.Errorf("Failed to append ca certs: %v", err)
	}

	return certPool, nil
}

// CreateTLSListener creates a TLS listener
func CreateTLSListener(grpcServerAddress, certFilePath, keyFilePath string, tlsEnabled bool) (net.Listener, error) {
	var err error
	var listener net.Listener
	if tlsEnabled {
		certPool, err := CreateCACertificatePool()
		if err != nil {
			return nil, fmt.Errorf("Failed to create CA certificate pool: %v", err)
		}

		cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
		if err != nil {
			return nil, fmt.Errorf("Could not load server key pair: %v", err)
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientCAs:    certPool,
		}
		listener, err = tls.Listen("tcp", grpcServerAddress, tlsConfig)
	} else {
		listener, err = net.Listen("tcp", grpcServerAddress)
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to listen: %v", err)
	}

	return listener, nil
}

// CreateTLSConfig creates a TLS config
func CreateTLSConfig() (*tls.Config, error) {
	caCertPool, err := CreateCACertificatePool()
	if err != nil {
		return nil, fmt.Errorf("Could not create CA certificate pool: %v", err)
	}
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	return tlsConfig, nil
}

// CreateGRPCConnection creates a gRPC connection
func CreateGRPCConnection(grpcServerAddress string, tlsEnabled bool) (*grpc.ClientConn, error) {
	tlsConfig, err := CreateTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("Could not create CA certificate pool: %v", err)
	}
	creds := credentials.NewTLS(tlsConfig)
	var connection *grpc.ClientConn
	if tlsEnabled {
		connection, err = grpc.Dial(grpcServerAddress, grpc.WithTransportCredentials(creds))
	} else {
		connection, err = grpc.Dial("qd.authentication.api:9090", grpc.WithInsecure())
	}
	if err != nil {
		return nil, fmt.Errorf("Could not connect to server: %v", err)
	}
	return connection, nil
}
