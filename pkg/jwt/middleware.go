package jwt

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/metadata"

	commonToken "github.com/quadev-ltd/qd-common/pkg/token"
)

// ContextKey is the key for the context
type ContextKey string

// Context keys
const (
	ClaimsContextKey ContextKey = "ContextClaimsKey"
	JWTTokenKey      ContextKey = "JWTTokenKey"
)

// TokenClaims is the claims for a JWT token
type TokenClaims struct {
	Email  string
	Type   commonToken.TokenType
	Expiry time.Time
	UserID string
}

// AddAuthorizationMetadataToContext adds the authorization Bearer token to the context
func AddAuthorizationMetadataToContext(ctx context.Context, token string) context.Context {
	existingMD, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		existingMD = metadata.New(map[string]string{})
	}
	newMD := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})
	mergedMD := metadata.Join(existingMD, newMD)
	return metadata.NewOutgoingContext(ctx, mergedMD)
}

// GetLoggerFromContext returns the logger from the context
func GetClaimsFromContext(ctx context.Context) (*TokenClaims, error) {
	if claims, ok := ctx.Value(ClaimsContextKey).(*TokenClaims); ok {
		return claims, nil
	}
	return nil, fmt.Errorf("Claims not found in context")
}
