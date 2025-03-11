package routes

import (
	"strings"

	"github.com/dunkykorZhik/avito-tech/internal/service"
	"github.com/gofiber/fiber/v2"
)

const userCtx = "Username"
const prefix = "Bearer "

func AuthMiddleware(service service.User) func(next fiber.Handler) fiber.Handler {
	return func(next fiber.Handler) fiber.Handler {
		return func(c *fiber.Ctx) error {
			tokenString, ok := getBearerToken(c)
			if !ok {
				return errorResponse(c, fiber.StatusUnauthorized, "unauthorized")
			}
			username, err := service.ParseToken(tokenString)
			if err != nil {
				return errorResponse(c, fiber.StatusUnauthorized, "unauthorized")
			}
			c.Context().SetUserValue(userCtx, username)

			return next(c)

		}

	}

}

func getBearerToken(c *fiber.Ctx) (string, bool) {
	header := c.Get("Authorization")
	if header == "" {

		return "", false
	}
	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != prefix {

		return "", false
	}
	return parts[1], true

}
