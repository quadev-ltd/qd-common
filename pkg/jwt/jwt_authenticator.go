package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// JWTAthenticatorer is an interface for JWTAuthenticator
type JWTAthenticatorer interface {
	VerifyToken(token string) (*jwt.Token, error)
	GetEmailFromToken(token *jwt.Token) (*string, error)
	GetExpiryFromToken(token *jwt.Token) (*time.Time, error)
}

// JWTAuthenticator is responsible for generating and verifying JWT tokens
type JWTAuthenticator struct {
	publicKey *rsa.PublicKey
}

var _ JWTAthenticatorer = &JWTAuthenticator{}

// Key constants
const (
	EmailClaim        = "email"
	ExpiryClaim       = "expiry"
	PublicKeyFileName = "public.pem"
	PublicKeyType     = "RSA PUBLIC KEY"
	KeysLocation      = "keys"
)

func createKeysFolderIfNotExists(fileLocation string) error {
	if _, err := os.Stat(fileLocation); os.IsNotExist(err) {
		err := os.Mkdir(fileLocation, 0700)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadPublicKeyFromString(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey.(*rsa.PublicKey), nil
}

// NewJWTAuthenticator creates a new JWT authenticator
func NewJWTAuthenticator(publicKeyString string) (JWTAthenticatorer, error) {
	publicKey, err := loadPublicKeyFromString(publicKeyString)
	if err != nil {
		return nil, err
	}
	return &JWTAuthenticator{
		publicKey,
	}, nil
}

// VerifyToken verifies a JWT token
func (authenticator *JWTAuthenticator) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return authenticator.publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("JWTAuthenticator: JWT Token is not valid")
	}
	expiry, err := authenticator.GetExpiryFromToken(token)
	if err != nil {
		return nil, err
	}
	if expiry.Before(time.Now()) {
		return nil, fmt.Errorf("JWTAuthenticator: JWT Token is expired")
	}
	return token, nil
}

// GetExpiryFromToken gets the expiry from a JWT token
func (authenticator *JWTAuthenticator) GetExpiryFromToken(token *jwt.Token) (*time.Time, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("JWTAuthenticator: JWT Token claims are not valid")
	}
	expiry, ok := claims[ExpiryClaim].(float64)
	if !ok {
		return nil, errors.New("JWTAuthenticator: JWT Token expiry is not valid")
	}
	expiryTime := time.Unix(int64(expiry), 0)
	return &expiryTime, nil
}

// GetEmailFromToken gets the email from a JWT token
func (authenticator *JWTAuthenticator) GetEmailFromToken(token *jwt.Token) (*string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("JWTAuthenticator: JWT Token claims are not valid")
	}
	email, ok := claims[EmailClaim].(string)
	if !ok {
		return nil, errors.New("JWTAuthenticator: JWT Token email is not valid")
	}
	return &email, nil
}
