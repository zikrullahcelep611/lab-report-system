package tokenservice

import (
	"context"
	"time"
)

type TokenRepository interface{
	DeleteExpiredTokens(ctx context.Context)
	AddTokenToBlackList(ctx context.Context, tokenStr string, expireTime time.Time) error
	IsTokenBlackListed(ctx context.Context, tokenStr string) bool
}

type TokenService struct{
	tokenRepository TokenRepository
}

func NewTokenService(tokenRepository TokenRepository) *TokenService{
	return &TokenService{tokenRepository: tokenRepository}
}

func (t *TokenService) DeleteExpiredTokens(ctx context.Context){
	t.tokenRepository.DeleteExpiredTokens(ctx)
}

func (t *TokenService) AddTokenToBlackList(ctx context.Context, tokenString string, expireTime time.Time) error{
	return t.tokenRepository.AddTokenToBlackList(ctx, tokenString, expireTime)
}

func (t *TokenService) IsTokenBlackListed(ctx context.Context, tokenString string) bool {
	return t.tokenRepository.IsTokenBlackListed(ctx, tokenString)
}