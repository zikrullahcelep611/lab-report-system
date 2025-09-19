package user

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app *fiber.App, handler *UserHandler) {
	api := app.Group("/api")

	api.Get("/users/:id", handler.GetUser)
	api.Post("/users", handler.CreateUser)
	api.Put("/users/:id", handler.UpdateUser)
	api.Delete("/users/:id", handler.DeleteUser)
}
