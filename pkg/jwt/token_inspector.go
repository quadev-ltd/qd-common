package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/quadev-ltd/qd-common/pkg/token"
)

// TokenInspectorer is TokenInspector interface
type TokenInspectorer interface {
	GetClaimFromToken(jwtToken *jwt.Token, claimKey string) (interface{}, error)
	GetEmailFromToken(jwtToken *jwt.Token) (*string, error)
	GetExpiryFromToken(jwtToken *jwt.Token) (*time.Time, error)
	GetTypeFromToken(jwtToken *jwt.Token) (*token.TokenType, error)
	GetUserIDFromToken(jwtToken *jwt.Token) (*string, error)
}

// TokenInspector is responsible for inspecting JWT token claims
type TokenInspector struct{}

var _ TokenInspectorer = &TokenInspector{}

func (inspector *TokenInspector) GetClaimFromToken(token *jwt.Token, claimKey string) (interface{}, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("JWT Token claims are not valid")
	}
	return claims[claimKey], nil
}

// GetExpiryFromToken gets the expiry from a JWT token
func (inspector *TokenInspector) GetExpiryFromToken(jwtToken *jwt.Token) (*time.Time, error) {
	expiry, err := inspector.GetClaimFromToken(jwtToken, ExpiryClaim)
	if err != nil {
		return nil, err
	}
	expiryTyped, ok := expiry.(float64)
	if !ok {
		return nil, errors.New("JWT Token expiry is not valid")
	}
	expiryTime := time.Unix(int64(expiryTyped), 0)
	return &expiryTime, nil
}

// GetEmailFromToken gets the email from a JWT token
func (inspector *TokenInspector) GetEmailFromToken(jwtToken *jwt.Token) (*string, error) {
	email, err := inspector.GetClaimFromToken(jwtToken, EmailClaim)
	if err != nil {
		return nil, err
	}
	emailTyped, ok := email.(string)
	if !ok {
		return nil, errors.New("JWT Token email is not of valid type")
	}
	return &emailTyped, nil
}

// GetTypeFromToken gets the type from a JWT token
func (inspector *TokenInspector) GetTypeFromToken(jwtToken *jwt.Token) (*token.TokenType, error) {
	typeValue, err := inspector.GetClaimFromToken(jwtToken, TypeClaim)
	if err != nil {
		return nil, err
	}
	typeValueString, ok := typeValue.(string)
	if !ok {
		return nil, errors.New("JWT Token type is not of valid type")
	}
	typeValueTokenType := token.TokenType(typeValueString)
	return &typeValueTokenType, nil
}

// GetUserIDFromToken gets the userID from a JWT token
func (inspector *TokenInspector) GetUserIDFromToken(
	jwtToken *jwt.Token,
) (*string, error) {
	userIDClaim, err := inspector.GetClaimFromToken(jwtToken, UserIDClaim)
	if err != nil {
		return nil, err
	}
	userID, ok := userIDClaim.(string)
	if !ok {
		return nil, errors.New("JWT Token type is not of valid type")
	}
	return &userID, nil
}
