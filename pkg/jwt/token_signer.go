package jwt

import (
	"crypto/rsa"

	"github.com/golang-jwt/jwt"
)

// TokenSignerer signs tokens
type TokenSignerer interface {
	SignToken(claims ...ClaimPair) (*string, error)
}

// TokenSigner signs tokens
type TokenSigner struct {
	rsaPrivateKey *rsa.PrivateKey
}

var _ TokenSignerer = &TokenSigner{}

// ClaimPair represents a key-value pair for a claim
type ClaimPair struct {
	Key   string
	Value interface{}
}

// NewTokenSigner creates a new JWT signer
func NewTokenSigner(rsaPrivateKey *rsa.PrivateKey) TokenSignerer {
	return &TokenSigner{
		rsaPrivateKey,
	}
}

// SignToken signs a JWT token
func (tokenSigner *TokenSigner) SignToken(claims ...ClaimPair) (*string, error) {
	tokenClaims := jwt.MapClaims{}
	for _, claim := range claims {
		tokenClaims[claim.Key] = claim.Value
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, tokenClaims)
	tokenString, err := token.SignedString(tokenSigner.rsaPrivateKey)
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}
