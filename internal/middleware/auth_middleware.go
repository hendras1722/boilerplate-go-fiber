package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/username/project-name/config"
	"github.com/username/project-name/domain/dto"
)

func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return dto.ErrorResponse(c, "Unauthorized", "Missing authorization header", fiber.StatusUnauthorized)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return dto.ErrorResponse(c, "Unauthorized", "Invalid authorization format", fiber.StatusUnauthorized)
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return dto.ErrorResponse(c, "Unauthorized", "Invalid or expired token", fiber.StatusUnauthorized)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return dto.ErrorResponse(c, "Unauthorized", "Invalid token claims", fiber.StatusUnauthorized)
		}

		// Store user info in context locals
		c.Locals("user_id", claims["sub"])
		c.Locals("user_email", claims["email"])

		return c.Next()
	}
}
