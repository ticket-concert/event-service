package middleware

import (
	"github.com/gofiber/fiber/v2"

	"event-service/internal/pkg/errors"
	"event-service/internal/pkg/helpers"
	"event-service/internal/pkg/log"
)

func AllowedRoles(roles ...string) fiber.Handler {
	logger := log.GetLogger()
	roleMap := make(map[string]struct{})
	for _, role := range roles {
		roleMap[role] = struct{}{}
	}

	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("userRole").(string)
		if !ok || role == "" {
			return helpers.RespError(c, logger, errors.ForbiddenError("Unauthorized role!"))
		}

		if _, ok := roleMap[role]; !ok {
			return helpers.RespError(c, logger, errors.ForbiddenError("Unauthorized role!"))
		}

		return c.Next()
	}
}
