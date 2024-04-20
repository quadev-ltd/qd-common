package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/quadev-ltd/qd-common/pkg/token"
)

// TokenInspectorer is TokenInspector interface
type TokenInspectorer interface {
	GetClaimFromToken(jwtToken *jwt.Token, claimKey string) (interface{}, error)
	GetEmailFromToken(jwtToken *jwt.Token) (*string, error)
	GetExpiryFromToken(jwtToken *jwt.Token) (*time.Time, error)
	GetTypeFromToken(jwtToken *jwt.Token) (*token.Type, error)
	GetUserIDFromToken(jwtToken *jwt.Token) (*string, error)
	GetClaimsFromToken(token *jwt.Token) (*TokenClaims, error)
	GetClaimsFromTokenString(tokenStr string) (*TokenClaims, error)
}

// TokenInspector is responsible for inspecting JWT token claims
type TokenInspector struct{}

var _ TokenInspectorer = &TokenInspector{}

// GetClaimFromToken gets a claim from a JWT token
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
func (inspector *TokenInspector) GetTypeFromToken(jwtToken *jwt.Token) (*token.Type, error) {
	typeValue, err := inspector.GetClaimFromToken(jwtToken, TypeClaim)
	if err != nil {
		return nil, err
	}
	typeValueString, ok := typeValue.(string)
	if !ok {
		return nil, errors.New("JWT Token type is not of valid type")
	}
	typeValueTokenType := token.Type(typeValueString)
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

// GetClaimsFromToken gets the email from a JWT token
func (inspector *TokenInspector) GetClaimsFromToken(token *jwt.Token) (*TokenClaims, error) {
	email, err := inspector.GetEmailFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("Error getting email from token: %v", err)
	}
	tokenType, err := inspector.GetTypeFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("Error getting type claim from token: %v", err)
	}
	expiry, err := inspector.GetExpiryFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("Error getting expiry claim from token: %v", err)
	}
	userID, err := inspector.GetUserIDFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("Error getting expiry claim from token: %v", err)
	}
	return &TokenClaims{
		Email:  *email,
		Type:   *tokenType,
		Expiry: *expiry,
		UserID: *userID,
	}, nil
}

// GetExpiryFromTokenString gets the email from a JWT token ßstringå
func (inspector *TokenInspector) GetClaimsFromTokenString(tokenStr string) (*TokenClaims, error) {
	// Parse the token
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	return inspector.GetClaimsFromToken(token)
}
