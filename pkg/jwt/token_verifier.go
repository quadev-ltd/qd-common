package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// TokenVerifierer is an interface for JWTAuthenticator
type TokenVerifierer interface {
	Verify(token string) (*jwt.Token, error)
}

// TokenVerifier is responsible for generating and verifying JWT tokens
type TokenVerifier struct {
	publicKey      *rsa.PublicKey
	tokenInspector TokenInspectorer
}

var _ TokenVerifierer = &TokenVerifier{}

func loadPublicKeyFromString(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey.(*rsa.PublicKey), nil
}

// NewTokenVerifier creates a new JWT authenticator
func NewTokenVerifier(publicKeyString string) (TokenVerifierer, error) {
	publicKey, err := loadPublicKeyFromString(publicKeyString)
	if err != nil {
		return nil, err
	}
	tokenInspector := &TokenInspector{}
	return &TokenVerifier{
		publicKey,
		tokenInspector,
	}, nil
}

// Verify verifies a JWT token
func (authenticator *TokenVerifier) Verify(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return authenticator.publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("Token Verifier: JWT Token is not valid")
	}
	expiry, err := authenticator.tokenInspector.GetExpiryFromToken(token)
	if err != nil {
		return nil, err
	}
	if expiry.Before(time.Now()) {
		return nil, fmt.Errorf("Token Verifier: JWT Token is expired")
	}
	return token, nil
}
