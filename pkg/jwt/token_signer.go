package jwt

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// TokenSignerer signs tokens
type TokenSignerer interface {
	SignToken(email string, expiry time.Time, tokenType TokenType) (*string, error)
}

// TokenSigner signs tokens
type TokenSigner struct {
	rsaPrivateKey *rsa.PrivateKey
}

var _ TokenSignerer = &TokenSigner{}

// NewTokenSigner creates a new JWT signer
func NewTokenSigner(rsaPrivateKey *rsa.PrivateKey) TokenSignerer {
	return &TokenSigner{
		rsaPrivateKey,
	}
}

// SignToken signs a JWT token
func (tokenSigner *TokenSigner) SignToken(email string, expiry time.Time, tokenType TokenType) (*string, error) {
	tokenClaims := jwt.MapClaims{
		EmailClaim:    email,
		ExpiryClaim:   expiry.Unix(),
		IssuedAtClaim: time.Now().Unix(),
		NonceClaim:    uuid.New(),
		TypeClaim:     tokenType,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, tokenClaims)
	tokenString, err := token.SignedString(tokenSigner.rsaPrivateKey)
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}
