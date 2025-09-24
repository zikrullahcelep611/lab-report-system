package jwtservice

import (
	"time"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/claims"
	customErrors "gitbub.com/zikrullahcelep611/lab-report/backend/models/errors"
	"gitbub.com/zikrullahcelep611/lab-report/backend/models/role"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type JwtService struct {
	jwtSecret []byte
}

func NewJwtService(jwtSecret string) *JwtService {
	return &JwtService{jwtSecret: []byte(jwtSecret)}
}

func (j *JwtService) GetJwtKey() []byte {
	return j.jwtSecret
}

func (j *JwtService) GenerateJwtToken(email string, role role.Role, expirationTime time.Time) (string, error) {
	claims := &claims.Claims{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.jwtSecret)
	if err != nil {
		log.Error().Str("operation", "GenerateJwtToken").Err(err).Msg("Could not create token")
		return "", err
	}

	return tokenString, nil
}

func (j *JwtService) ParseToken(tokenStr string) (*claims.Claims, error) {
	if tokenStr == "" {
		return nil, &customErrors.TokenIsNullError{Message: "Token is required"}
	}

	claims := &claims.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return j.GetJwtKey(), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, &customErrors.InvalidTokenError{Message: "Token is invalid"}
	}

	return claims, nil
}

func (j *JwtService) ParseTokenFromCookie(c *fiber.Ctx) (*claims.Claims, error) {
	tokenStr := c.Cookies("token")
	if tokenStr == "" {
		log.Warn().Str("operation", "ParseTokenFromCookie").Msg("Token cookie not found")
		return nil, &customErrors.TokenIsNullError{Message: "Token not found in cookies"}
	}

	log.Info().Str("operation", "ParseTokenFromCookie").Msg("Token found in cookie")
	return j.ParseToken(tokenStr)
}

/*
	Kullanıcı → Web Sitesi → Çerezde Token Saklanır
                       ↓
	Sonraki İsteklerde → ParseTokenFromCookie (çerezden token al)
                       ↓
                   ParseToken (token'ı çöz ve kontrol et)
                       ↓
                   "Bu kullanıcı kimmiş? Ne yapabilir?"

*/
