package jwt

import (
	"crypto/rsa"
	"fmt"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"github.com/quadev-ltd/qd-common/pkg/token"
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
	tokenClaims := jwt.MapClaims{
		IssuedAtClaim: time.Now().Unix(),
		NonceClaim:    uuid.New(),
	}
	for _, claim := range claims {
		switch v := claim.Value.(type) {
		case time.Time:
			tokenClaims[claim.Key] = v.Unix()
		case string:
			tokenClaims[claim.Key] = string(v)
		case token.TokenType:
			tokenClaims[claim.Key] = string(v)
		default:
			return nil, fmt.Errorf("invalid claim value type: %v", reflect.TypeOf(claim.Value))
		}

	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, tokenClaims)
	tokenString, err := token.SignedString(tokenSigner.rsaPrivateKey)
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}
