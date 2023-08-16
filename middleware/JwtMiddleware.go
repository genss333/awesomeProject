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
		token, _ := jwt.Parse(tokenString, nil)
		claims, _ := token.Claims.(jwt.MapClaims)
		roles, _ := claims["roles"].([]interface{})

		for _, allowedRole := range allowedRoles {
			for _, userRole := range roles {
				if allowedRole == userRole {
					return c.Next()
				}
			}
		}

		return utils.RespondJson(c, fiber.StatusUnauthorized, "Unauthorized")
	}
}
