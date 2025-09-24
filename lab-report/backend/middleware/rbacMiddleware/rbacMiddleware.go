package rbacmiddleware

import (
	"gitbub.com/zikrullahcelep611/lab-report/backend/models/claims"
	"gitbub.com/zikrullahcelep611/lab-report/backend/models/role"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func RequireRole(allowedRoles ...role.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims, ok := c.Locals("user").(*claims.Claims)
		if !ok {
			log.Warn().Str("operation", "RequireRole").Msg("User claims not found in context")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authentication required",
			})
		}

		userRole := userClaims.Role
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				log.Info().
					Str("operation", "RequireRole").
					Str("user_email", userClaims.Email).
					Str("user_role", string(userRole)).
					Str("endpoint", c.Path()).
					Msg("Access granted")
				return c.Next()
			}
		}

		// Yetkisiz erişim
		log.Warn().
			Str("operation", "RequireRole").
			Str("user_email", userClaims.Email).
			Str("user_role", string(userRole)).
			Str("endpoint", c.Path()).
			Str("method", c.Method()).
			Msg("Access denied - insufficient permissions")

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "Insufficient permissions",
			"message": "You don't have the required role to access this resource",
		})
	}
}

// Özel role kontrolleri için helper fonksiyonlar
func RequireSuperAdmin() fiber.Handler {
	return RequireRole(role.RoleSuperAdmin)
}

func RequireAdminOrSuperAdmin() fiber.Handler {
	return RequireRole(role.RoleAdmin, role.RoleSuperAdmin)
}

func RequireAnyRole() fiber.Handler {
	return RequireRole(
		role.RoleSuperAdmin,
		role.RoleAdmin,
		role.RoleTechnician,
	)
}
