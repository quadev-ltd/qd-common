package token

// TokenType is the type for the token type
type TokenType string

// Token types
const (
	EmailVerificationTokenType TokenType = "EmailVerificationTokenType"
	ResetPasswordTokenType     TokenType = "ResetPasswordTokenType"
	AccessTokenType            TokenType = "AccessTokenType"
	RefreshTokenType           TokenType = "RefreshTokenType"
)
