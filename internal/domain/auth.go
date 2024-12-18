package domain

import "context"

type Auth interface {
	GenerateTokens(ctx context.Context, userID string, clientIP string) (*TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string, clientIP string) (*TokenPair, error)
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type User struct {
	ID    string
	Email string
}

type EmailSender interface {
	SendEmail(to, subject, body string) error
}
