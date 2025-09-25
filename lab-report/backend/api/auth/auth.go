package auth

import (
	"context"

	"time"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/auth"
	"gitbub.com/zikrullahcelep611/lab-report/backend/models/role"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type AuthService interface {
	Login(ctx context.Context, email string, password string) (auth.Login, error)
	Logout(ctx context.Context, tokenString string) error
}

type JwtService interface {
	GenerateJwtToken(email string, role role.Role, expirationTime time.Time) (string, error)
}

type AuthHandler struct {
	authService AuthService
	jwtService  JwtService
}

func NewAuthController(authService AuthService, jwtService JwtService) *AuthHandler {
	return &AuthHandler{authService: authService, jwtService: jwtService}
}

func (a *AuthHandler) Login(c *fiber.Ctx) error {
	ctx := c.Context()
	var creds auth.Login
	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	log.Warn().Str("claimden gelen password", creds.Password).Str("claimden gelen email", creds.Email).Msg("test ")

	user, err := a.authService.Login(ctx, creds.Email, creds.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	expirationTime := time.Now().Add(time.Hour * 24)
	tokenString, err := a.jwtService.GenerateJwtToken(user.Email, user.Role, expirationTime)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create token",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Log in successfully",
	})
}

func (a *AuthHandler) Logout(c *fiber.Ctx) error {
	ctx := c.Context()
	tokenString := c.Cookies("token")
	if tokenString == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token not found",
		})
	}

	err := a.authService.Logout(ctx, tokenString)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to logout",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Unix(0, 0),
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}
