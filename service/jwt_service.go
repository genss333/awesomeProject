package service

import (
	"awesomeProject/json"
	"awesomeProject/models"
	"awesomeProject/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"strings"
	"time"
)

var secretKey = []byte(utils.GoDotEnvVariable("JWT_SECRET_KEY"))

var activeTokens = make(map[string]int64)

func JWTMiddleware(c *fiber.Ctx) error {
	authorizationHeader := c.Get("Authorization")
	if authorizationHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing authorization header",
			"status":  fiber.StatusUnauthorized,
			"path":    c.Path(),
		})
	}

	if !strings.HasPrefix(authorizationHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token format",
			"status":  fiber.StatusUnauthorized,
			"path":    c.Path(),
		})
	}

	tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token signature",
				"status":  fiber.StatusUnauthorized,
				"path":    c.Path(),
			})
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
			"status":  fiber.StatusUnauthorized,
			"path":    c.Path(),
		})
	}

	if token.Valid {
		tokenExp := int64(token.Claims.(jwt.MapClaims)["exp"].(float64))

		if _, exists := activeTokens[tokenString]; !exists {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token is no longer active",
				"status":  fiber.StatusUnauthorized,
				"path":    c.Path(),
			})
		}

		currentTime := time.Now().Unix()
		if currentTime > tokenExp {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token has expired",
				"status":  fiber.StatusUnauthorized,
				"path":    c.Path(),
			})
		}

		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "Invalid token",
		"status":  fiber.StatusUnauthorized,
		"path":    c.Path(),
	})
}

func GenerateToken(user models.User) (json.AuthJson, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.UserId
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 3).Unix()

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return json.AuthJson{}, err
	}

	var auth = json.AuthJson{
		User:     user,
		Token:    tokenString,
		TokenExp: claims["exp"].(int64),
	}

	activeTokens[tokenString] = auth.TokenExp

	return auth, nil
}

func GetCurrentUserFromToken(tokenString string) (json.AuthTokenJson, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return json.AuthTokenJson{}, err
	}

	claims := token.Claims.(jwt.MapClaims)

	var authJson = json.AuthTokenJson{
		UserId:   int(claims["user_id"].(float64)),
		Username: claims["username"].(string),
		TokenExp: int64(claims["exp"].(float64)), // Convert to int64
	}

	return authJson, nil
}

func CurrentUser(c *fiber.Ctx) (json.AuthTokenJson, error) {
	var token = c.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}
	authJson, err := GetCurrentUserFromToken(token)
	if err != nil {
		return json.AuthTokenJson{}, utils.RespondJson(c, fiber.StatusBadRequest, "Invalid token")
	}

	return authJson, nil
}

func Logout(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}

	RevokeToken(token)
	return c.SendStatus(fiber.StatusOK)
}

func RevokeToken(token string) {
	delete(activeTokens, token)
}

func LogRequests(c *fiber.Ctx) error {
	log.Println(c.Method(), c.Path())
	return c.Next()
}
