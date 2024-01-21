package tls

import (
	"crypto/x509"
	"fmt"
	"os"
)

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
