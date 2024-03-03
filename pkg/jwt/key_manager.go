package jwt

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/quadev-ltd/qd-common/pkg/fs"
)

// KeyManagerer handles key generation and retrieval
type KeyManagerer interface {
	GenerateNewKeyPair() error
	GetRSAPrivateKey() *rsa.PrivateKey
	GetRSAPublicKey() *rsa.PublicKey
	GetPublicKey(ctx context.Context) (string, error)
}

// KeyManager is responsible for generating and managing RSA keys
type KeyManager struct {
	fileLocation string
	privateKey   *rsa.PrivateKey
	publicKey    *rsa.PublicKey
	fs           fs.FileSystem
}

var _ KeyManagerer = &KeyManager{}

// NewKeyManager creates a new JWT signer
func NewKeyManager(fileLocation string) (KeyManagerer, error) {
	fs := &fs.OSFileSystem{}
	privateKey, err := loadPrivateKeyFromFile(
		fmt.Sprintf("%s/%s", fileLocation, PrivateKeyFileName),
		fs,
	)
	if err != nil && fs.IsNotExist(err) {
		privateKey, publicKey, err := generateKeyFiles(fileLocation, fs)
		if err != nil {
			return nil, err
		}
		return &KeyManager{
			privateKey:   privateKey,
			publicKey:    publicKey,
			fileLocation: fileLocation,
			fs:           fs,
		}, nil
	} else if err != nil {
		return nil, err
	}
	publicKey, err := loadPublicKeyFromFile(
		fmt.Sprintf("%s/%s", fileLocation, PublicKeyFileName),
		fs,
	)
	if err != nil {
		return nil, err
	}
	return &KeyManager{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func createKeysFolderIfNotExists(fileLocation string, fs fs.FileSystem) error {
	if _, err := fs.Stat(fileLocation); fs.IsNotExist(err) {
		err := fs.Mkdir(fileLocation, 0700)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	publicKey := &privateKey.PublicKey

	return privateKey, publicKey, nil
}

func savePrivateKeyToFile(privateKey *rsa.PrivateKey, filename string, fs fs.FileSystem) error {
	file, err := fs.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  PrivateKeyType,
		Bytes: privateKeyBytes,
	})

	_, err = file.Write(privateKeyPEM)
	return err
}

func savePublicKeyToFile(publicKey *rsa.PublicKey, filename string, fs fs.FileSystem) error {
	file, err := fs.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  PublicKeyType,
		Bytes: publicKeyBytes,
	})

	_, err = file.Write(publicKeyPEM)
	return err
}

func loadPrivateKeyFromFile(filename string, fs fs.FileSystem) (*rsa.PrivateKey, error) {
	privateKeyPEM, err := fs.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privateKeyPEM)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func loadPublicKeyFromFile(filename string, fs fs.FileSystem) (*rsa.PublicKey, error) {
	publicKeyPEM, err := fs.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(publicKeyPEM)
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey.(*rsa.PublicKey), nil
}

func generateKeyFiles(fileLocation string, fs fs.FileSystem) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	createKeysFolderIfNotExists(fileLocation, fs)
	privateKey, publicKey, err := generateKeyPair()
	if err != nil {
		return nil, nil, err
	}
	err = savePrivateKeyToFile(
		privateKey,
		fmt.Sprintf("%s/%s", fileLocation, PrivateKeyFileName),
		fs,
	)
	if err != nil {
		return nil, nil, err
	}
	err = savePublicKeyToFile(
		publicKey,
		fmt.Sprintf("%s/%s", fileLocation, PublicKeyFileName),
		fs,
	)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, publicKey, nil
}

// GenerateNewKeyPair generates a new key pair
func (keyManager *KeyManager) GenerateNewKeyPair() error {
	privateKey, publicKey, err := generateKeyFiles(keyManager.fileLocation, keyManager.fs)
	if err != nil {
		return err
	}
	keyManager.privateKey = privateKey
	keyManager.publicKey = publicKey
	return nil
}

// GetRSAPrivateKey gets the RSA private key
func (keyManager *KeyManager) GetRSAPrivateKey() *rsa.PrivateKey {
	return keyManager.privateKey
}

// GetRSAPublicKey gets the RSA public key
func (keyManager *KeyManager) GetRSAPublicKey() *rsa.PublicKey {
	return keyManager.publicKey
}

// GetPublicKey gets the public key
func (keyManager *KeyManager) GetPublicKey(ctx context.Context) (string, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(keyManager.publicKey)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal public key: %v", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  PublicKeyType,
		Bytes: publicKeyBytes,
	})
	return string(publicKeyPEM), nil
}
