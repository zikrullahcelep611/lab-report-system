package auth

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(app *fiber.App, handler *AuthHandler) {
	api := app.Group("/api")

	api.Post("/login", handler.Login)
	api.Post("/logout", handler.Logout)
}
