package postgres

import (
	"context"
	"database/sql"

	"github.com/MosinEvgeny/authservice/internal/domain"
	"github.com/MosinEvgeny/authservice/internal/repository"

	_ "github.com/lib/pq"
)

type refreshTokenRepository struct {
	db *sql.DB
}

// NewRefreshTokenRepository создает новый репозиторий для refresh токенов.
func NewRefreshTokenRepository(db *sql.DB) repository.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) StoreRefreshToken(ctx context.Context, userID, hashedRefreshToken, accessToken, clientIP string) error {

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO refresh_tokens (user_id, hashed_refresh_token, access_token, client_ip)
		VALUES ($1, $2, $3, $4)
	`, userID, hashedRefreshToken, accessToken, clientIP)

	return err

}

func (r *refreshTokenRepository) GetRefreshTokenData(ctx context.Context, refreshToken string) (*repository.RefreshTokenData, error) {
	var refreshTokenData repository.RefreshTokenData
	err := r.db.QueryRowContext(ctx, `
		SELECT user_id, hashed_refresh_token, access_token, client_ip
		FROM refresh_tokens
		WHERE hashed_refresh_token = $1
	`, refreshToken).Scan(&refreshTokenData.UserID, &refreshTokenData.HashedRefreshToken, &refreshTokenData.AccessToken, &refreshTokenData.ClientIP)

	if err != nil {
		return nil, err

	}

	return &refreshTokenData, nil

}

func (r *refreshTokenRepository) UpdateRefreshToken(ctx context.Context, userID string, hashedRefreshToken string, accessToken string, clientIP string) error {

	_, err := r.db.ExecContext(ctx, `
		UPDATE refresh_tokens
		SET hashed_refresh_token = $1, access_token=$2, client_ip=$3
		WHERE user_id = $4
		
		`, hashedRefreshToken, accessToken, clientIP, userID)

	return err

}

func (r *refreshTokenRepository) GetUserByID(ctx context.Context, userId string) (*domain.User, error) {

	var user domain.User
	err := r.db.QueryRowContext(ctx, `
        SELECT id, email
        FROM users
        WHERE id = $1
    `, userId).Scan(&user.ID, &user.Email)

	if err != nil {
		return nil, err

	}

	return &user, nil

}
