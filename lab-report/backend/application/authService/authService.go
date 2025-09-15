package authservice

import (
	"context"
	"time"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/auth"
	"gitbub.com/zikrullahcelep611/lab-report/backend/models/claims"
	customErrors "gitbub.com/zikrullahcelep611/lab-report/backend/models/errors"
	"github.com/rs/zerolog/log"
)

type LoginRepository interface {
	Login(ctx context.Context, email string, password string) (auth.Login, error)
}

type TokenRepository interface {
	AddTokenToBlackList(ctx context.Context, tokenStr string, expireTime time.Time) error
	IsTokenBlackListed(ctx context.Context, tokenStr string) bool
}

type JwtService interface {
	ParseToken(tokenStr string) (*claims.Claims, error)
}

type AuthService struct {
	loginRepository LoginRepository
	tokenRepository TokenRepository
	jwtService      JwtService
}

func NewAuthService(loginRepository LoginRepository, tokenRepository TokenRepository, jwtService JwtService) *AuthService {
	return &AuthService{
		loginRepository: loginRepository,
		tokenRepository: tokenRepository,
		jwtService:      jwtService,
	}
}

func (a *AuthService) Login(ctx context.Context, email string, password string) (auth.Login, error) {
	return a.loginRepository.Login(ctx, email, password)
}

func (a *AuthService) Logout(ctx context.Context, tokenString string) error {
	log.Info().Str("operation", "Logout").Msg("User logout attempt")

	claims, err := a.jwtService.ParseToken(tokenString)
	if err != nil {
		log.Error().Str("operation", "Logout").Err(err).Msg("Invalid token for logout")
		return &customErrors.InvalidTokenError{Message: "Invalid token"}
	}

	err = a.tokenRepository.AddTokenToBlackList(ctx, tokenString, claims.ExpiresAt.Time)
	if err != nil {
		log.Error().Str("operation", "Logout").Err(err).Msg("Failed to add token to blacklist")
		return err
	}

	log.Info().Str("operation", "Logout").Str("email", claims.Email).Msg("User logut successfully")
	return nil
}
