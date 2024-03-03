package jwt

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
		tokenString, err := keySigner.SignToken(email, expiry, AccessTokenType)
		if err != nil {
			t.Fatal(err)
		}
		token, err := tokenVerifier.Verify(*tokenString)
		assert.NoError(t, err)
		emailClaim, err := tokenInspector.GetEmailFromToken(token)
		assert.NoError(t, err)
		assert.Equal(t, *emailClaim, email)
		expiryClaim, err := tokenInspector.GetExpiryFromToken(token)
		assert.NoError(t, err)
		assert.Equal(t, (*expiryClaim).Day(), expiry.Day())
		assert.Equal(t, (*expiryClaim).Hour(), expiry.Hour())
		typeClaim, err := tokenInspector.GetTypeFromToken(token)
		assert.NoError(t, err)
		assert.Equal(t, string(AccessTokenType), *typeClaim)
	})

}
