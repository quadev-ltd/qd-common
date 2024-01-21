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
func CreateTLSListener(grpcServerAddress, certFile, keyFile string) (net.Listener, error) {
	certPool, err := CreateCACertificatePool()
	if err != nil {
		return nil, fmt.Errorf("Failed to create CA certificate pool: %v", err)
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("Could not load server key pair: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
	}

	listener, err := tls.Listen("tcp", grpcServerAddress, tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to listen: %v", err)
	}

	return listener, nil
}

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

func CreateGRPCConnection(grpcServerAddress string) (*grpc.ClientConn, error) {
	tlsConfig, err := CreateTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("Could not create CA certificate pool: %v", err)
	}
	creds := credentials.NewTLS(tlsConfig)
	connection, err := grpc.Dial(grpcServerAddress, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("Could not connect to server: %v", err)
	}
	return connection, nil
}
