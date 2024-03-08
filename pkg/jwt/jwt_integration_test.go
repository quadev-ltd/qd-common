package jwt

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/quadev-ltd/qd-common/pkg/token"
)

func TestTokenInspector(t *testing.T) {
	keyManager, err := NewKeyManager("./keys")
	if err != nil {
		t.Fatal(err)
	}
	keySigner := NewTokenSigner(keyManager.GetRSAPrivateKey())
	tokenInspector := &TokenInspector{}
	publicKey, err := keyManager.GetPublicKey(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	tokenVerifier, err := NewTokenVerifier(publicKey)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Test_Claims", func(t *testing.T) {
		email := "test@email.com"
		expiry := time.Now().Add(2 * time.Hour)
		userID := primitive.NewObjectID().Hex()
		claims := []ClaimPair{
			{EmailClaim, email},
			{ExpiryClaim, expiry},
			{TypeClaim, token.AccessTokenType},
			{UserIDClaim, userID},
		}
		tokenString, err := keySigner.SignToken(claims...)
		if err != nil {
			t.Fatal(err)
		}
		jwtToken, err := tokenVerifier.Verify(*tokenString)

		assert.NoError(t, err)
		emailClaim, err := tokenInspector.GetEmailFromToken(jwtToken)
		assert.NoError(t, err)
		assert.Equal(t, *emailClaim, email)
		expiryClaim, err := tokenInspector.GetExpiryFromToken(jwtToken)
		assert.NoError(t, err)
		assert.Equal(t, (*expiryClaim).Day(), expiry.Day())
		assert.Equal(t, (*expiryClaim).Hour(), expiry.Hour())
		typeClaim, err := tokenInspector.GetTypeFromToken(jwtToken)
		assert.NoError(t, err)
		assert.Equal(t, token.AccessTokenType, *typeClaim)
		userIDClaim, err := tokenInspector.GetUserIDFromToken(jwtToken)
		assert.NoError(t, err)
		assert.Equal(t, userID, *userIDClaim)
	})

}
