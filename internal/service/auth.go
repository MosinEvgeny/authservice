package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/MosinEvgeny/authservice/internal/config"
	"github.com/MosinEvgeny/authservice/internal/domain"
	"github.com/MosinEvgeny/authservice/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo          repository.RefreshTokenRepository
	cfg           *config.Config
	emailSender   domain.EmailSender
	tokenDuration time.Duration
}

func NewAuthService(repo repository.RefreshTokenRepository, cfg *config.Config, emailSender domain.EmailSender) domain.Auth {
	return &authService{
		repo:          repo,
		cfg:           cfg,
		emailSender:   emailSender,
		tokenDuration: time.Hour,
	}
}

func (s *authService) GenerateTokens(ctx context.Context, userID string, clientIP string) (*domain.TokenPair, error) {

	accessToken, err := s.generateAccessToken(userID, clientIP)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token %w", err)
	}
	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash refresh token: %w", err)
	}

	err = s.repo.StoreRefreshToken(ctx, userID, string(hashedRefreshToken), accessToken, clientIP)

	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string, clientIP string) (*domain.TokenPair, error) {
	refreshTokenData, err := s.repo.GetRefreshTokenData(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token data: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(refreshTokenData.HashedRefreshToken), []byte(refreshToken)); err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)

	}

	if clientIP != refreshTokenData.ClientIP {
		user, err := s.repo.GetUserByID(ctx, refreshTokenData.UserID)

		if err != nil {
			return nil, fmt.Errorf("failed to get user by ID: %w", err)
		}
		err = s.emailSender.SendEmail(user.Email, "Security Alert", "Your IP address has changed!")

		if err != nil {
			return nil, fmt.Errorf("failed to send email: %w", err)
		}

	}

	accessToken, err := s.generateAccessToken(refreshTokenData.UserID, clientIP)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.generateRefreshToken()

	if err != nil {
		return nil, fmt.Errorf("failed to generate new refresh token: %w", err)
	}
	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(newRefreshToken), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash new refresh token: %w", err)
	}
	err = s.repo.UpdateRefreshToken(ctx, refreshTokenData.UserID, string(hashedRefreshToken), accessToken, clientIP)

	if err != nil {
		return nil, fmt.Errorf("failed to update refresh token: %w", err)
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil

}

func (s *authService) generateAccessToken(userID string, clientIP string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user_id":   userID,
		"exp":       time.Now().Add(s.tokenDuration).Unix(),
		"client_ip": clientIP,
	})
	return token.SignedString([]byte(s.cfg.JWTSecret))

}

func (s *authService) generateRefreshToken() (string, error) {

	randomBytes := make([]byte, 32)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", err

	}
	return base64.StdEncoding.EncodeToString(randomBytes), nil

}
