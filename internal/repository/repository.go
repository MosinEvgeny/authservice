package repository

import (
	"context"

	"github.com/MosinEvgeny/authservice/internal/domain"
)

type RefreshTokenRepository interface {
	StoreRefreshToken(ctx context.Context, userID, hashedRefreshToken, accessToken, clientIP string) error
	GetRefreshTokenData(ctx context.Context, refreshToken string) (*RefreshTokenData, error)
	UpdateRefreshToken(ctx context.Context, userID string, hashedRefreshToken string, accessToken string, ClientIP string) error
	GetUserByID(ctx context.Context, userId string) (*domain.User, error)
}

type RefreshTokenData struct {
	UserID             string
	HashedRefreshToken string
	AccessToken        string
	ClientIP           string
}
