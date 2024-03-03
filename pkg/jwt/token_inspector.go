package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

// TokenInspector inspects tokens
type TokenInspectorer interface {
	GetEmailFromToken(token *jwt.Token) (*string, error)
	GetExpiryFromToken(token *jwt.Token) (*time.Time, error)
}

type TokenInspector struct{}

var _ TokenInspectorer = &TokenInspector{}

// GetExpiryFromToken gets the expiry from a JWT token
func (inspector *TokenInspector) GetExpiryFromToken(token *jwt.Token) (*time.Time, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("JWT Token claims are not valid")
	}
	expiry, ok := claims[ExpiryClaim].(float64)
	if !ok {
		return nil, errors.New("JWT Token expiry is not valid")
	}
	expiryTime := time.Unix(int64(expiry), 0)
	return &expiryTime, nil
}

// GetEmailFromToken gets the email from a JWT token
func (inspector *TokenInspector) GetEmailFromToken(token *jwt.Token) (*string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("JWT Token claims are not valid")
	}
	email, ok := claims[EmailClaim].(string)
	if !ok {
		return nil, errors.New("JWT Token email is not valid")
	}
	return &email, nil
}
