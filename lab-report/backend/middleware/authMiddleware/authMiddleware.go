package authmiddleware

import (
	"context"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/claims"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type TokenService interface {
	IsTokenBlackListed(ctx context.Context, tokenString string) bool
}

type JwtService interface {
	GetJwtKey() []byte
}

type AuthMiddleware struct {
	jwtService   JwtService
	tokenService TokenService
}

func NewAuthMiddleware(tokenService TokenService, jwtService JwtService) *AuthMiddleware {
	return &AuthMiddleware{tokenService: tokenService, jwtService: jwtService}
}

func (auth *AuthMiddleware) Authenticate(c *fiber.Ctx) error {
	tokenString := c.Cookies("token")
	if tokenString == "" {
		log.Warn().Str("operation", "Authenticate").Msg("Token is blacklisted")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token is invalid",
		})
	}

	claims := &claims.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return auth.jwtService.GetJwtKey(), nil
	})

	if err != nil {
		log.Error().Err(err).Str("operation", "Authenticate").Msg("Token parsing failed")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	if !token.Valid {
		log.Warn().Str("operation", "Authenticate").Msg("Invalid token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token is not valid",
		})
	}

	c.Locals("user", claims)

	log.Info().Str("operation", "Authenticate").Str("email", claims.Email).Msg("User authenticated successfully")

	return c.Next()
}
