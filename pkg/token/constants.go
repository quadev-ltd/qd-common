package token

// Type is the type for the token type
type Type string

// Token types
const (
	EmailVerificationTokenType Type = "EmailVerificationTokenType"
	ResetPasswordTokenType     Type = "ResetPasswordTokenType"
	AuthTokenType              Type = "AuthTokenType"
	RefreshTokenType           Type = "RefreshTokenType"
	AllTokenType               Type = "AllTokenType"
)
