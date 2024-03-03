package jwt

// Key constants
const (
	EmailClaim         = "email"
	ExpiryClaim        = "expiry"
	IssuedAtClaim      = "iat"
	NonceClaim         = "nonce"
	TypeClaim          = "type"
	PublicKeyFileName  = "public.pem"
	PrivateKeyFileName = "private.pem"
	PublicKeyType      = "RSA PUBLIC KEY"
	PrivateKeyType     = "RSA PRIVATE KEY"
)

// TokenType is the type for the token type
type TokenType string

// Token types
const (
	EmailVerificationTokenType TokenType = "EmailVerificationTokenType"
	ResetPasswordTokenType     TokenType = "ResetPasswordTokenType"
	AccessTokenType            TokenType = "AccessTokenType"
	RefreshTokenType           TokenType = "RefreshTokenType"
)
