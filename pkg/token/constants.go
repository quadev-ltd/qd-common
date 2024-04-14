package token

// Type is the type for the token type
type Type string

// Token types
const (
	EmailVerificationTokenType Type = "EmailVerificationTokenType"
	ResetPasswordTokenType     Type = "ResetPasswordTokenType"
	AccessTokenType            Type = "AccessTokenType"
	RefreshTokenType           Type = "RefreshTokenType"
)
