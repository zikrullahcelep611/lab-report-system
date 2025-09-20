package contexttimeoutmiddleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TimeoutMiddleware(timeoutValue int) fiber.Handler {
	timeout := time.Duration(timeoutValue) * time.Second

	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), timeout)
		defer cancel()

		c.SetUserContext(ctx)

		done := make(chan error, 1)

		go func() {
			done <- c.Next()
		}()

		select {
		case err := <-done:
			return err
		case <-ctx.Done():
			return c.Status(fiber.StatusGatewayTimeout).JSON(fiber.Map{
				"error": "Request timed out",
			})
		}
	}
}
