package jwt

import (
	"context"

	"google.golang.org/grpc/metadata"
)

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
