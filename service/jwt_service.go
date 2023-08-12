package service

import (
	"awesomeProject/json"
	"awesomeProject/models"
	"awesomeProject/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

type JWTService struct {
	SecretKey []byte
}

var secretKey = JWTService{
	SecretKey: []byte("8Zz5tw0Ion3XPZZfN0NOml3z9FMultiwordR9fp6ryDIoGRM8STEPHA6iHsc0fb"),
}

func GenerateToken(user models.User) (json.AuthJson, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.UserId
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(secretKey.SecretKey)
	if err != nil {
		return json.AuthJson{}, err
	}

	var auth = json.AuthJson{
		User:     user,
		Token:    tokenString,
		TokenExp: claims["exp"].(int64),
	}

	return auth, nil
}

func GetCurrentUserFromToken(tokenString string) (json.AuthTokenJson, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey.SecretKey, nil
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

func CurrentUser(c *fiber.Ctx) error {
	var token = c.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}
	authJson, err := GetCurrentUserFromToken(token)
	if err != nil {
		utils.RespondWithError(c, fiber.StatusBadRequest, "Invalid token")
		return err
	}

	return c.JSON(authJson)
}
