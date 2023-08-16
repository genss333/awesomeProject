package middleware

import (
	"awesomeProject/service"
	"awesomeProject/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

var activeTokens = service.GetActiveTokens()
var secretKey = []byte(utils.GoDotEnvVariable("JWT_SECRET_KEY"))

func AuthenticationMiddleware(c *fiber.Ctx) error {
	authorizationHeader := c.Get("Authorization")
	if authorizationHeader == "" {
		return utils.RespondJson(c, fiber.StatusUnauthorized, "Middleware Missing token")
	}

	if !strings.HasPrefix(authorizationHeader, "Bearer ") {
		return utils.RespondJson(c, fiber.StatusUnauthorized, "Middleware Invalid token")
	}

	tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return utils.RespondJson(c, fiber.StatusUnauthorized, "Middleware Invalid token")
		}
		return utils.RespondJson(c, fiber.StatusUnauthorized, "Middleware Missing token")
	}

	if token.Valid {
		tokenExp := int64(token.Claims.(jwt.MapClaims)["exp"].(float64))

		if _, exists := activeTokens[tokenString]; !exists {
			return utils.RespondJson(c, fiber.StatusUnauthorized, "Invalid token")
		}

		currentTime := time.Now().Unix()
		if currentTime > tokenExp {
			return utils.RespondJson(c, fiber.StatusUnauthorized, "Token has expired")
		}

		return c.Next()
	}

	return utils.RespondJson(c, fiber.StatusUnauthorized, "Middleware Invalid token")
}

func AuthorizationMiddleware(allowedRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return utils.RespondJson(c, fiber.StatusUnauthorized, "Missing Authorization header")
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			return utils.RespondJson(c, fiber.StatusUnauthorized, "Invalid token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return utils.RespondJson(c, fiber.StatusUnauthorized, "Invalid token claims")
		}

		rolesInterface, exists := claims["roles"]
		if !exists {
			return utils.RespondJson(c, fiber.StatusForbidden, "Access denied")
		}

		for _, allowedRole := range allowedRoles {
			if rolesInterface == allowedRole {
				return c.Next()
			}
		}

		return utils.RespondJson(c, fiber.StatusForbidden, "Access denied")
	}
}
