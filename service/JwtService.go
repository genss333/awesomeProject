package service

import (
	"awesomeProject/json"
	"awesomeProject/models"
	"awesomeProject/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

var secretKey = []byte(utils.GoDotEnvVariable("JWT_SECRET_KEY"))

var activeTokens = make(map[string]int64)
var refreshTokens = make(map[string]int64)

func GetActiveTokens() map[string]int64 {
	return activeTokens
}

func GenerateToken(user models.User) (json.AuthJson, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims["user_id"] = user.UserId
	accessClaims["username"] = user.Username
	accessClaims["roles"] = user.Role
	accessClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	accessTokenString, err := accessToken.SignedString(secretKey)
	if err != nil {
		return json.AuthJson{}, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["user_id"] = user.UserId
	refreshClaims["username"] = user.Username
	refreshClaims["roles"] = user.Role
	refreshClaims["exp"] = time.Now().Add(time.Minute * 20).Unix()

	refreshTokenString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return json.AuthJson{}, err
	}

	var auth = json.AuthJson{
		UserId:       user.UserId,
		Username:     user.Username,
		Token:        accessTokenString,
		RefreshToken: refreshTokenString,
		TokenExp:     accessClaims["exp"].(int64),
	}

	activeTokens[accessTokenString] = auth.TokenExp
	refreshTokens[refreshTokenString] = refreshClaims["exp"].(int64)

	return auth, nil
}

func RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Get("Authorization")
	if strings.HasPrefix(refreshToken, "Bearer ") {
		refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")
	}

	authToken, err := GetCurrentUserFromToken(refreshToken)
	if err != nil {
		return utils.RespondJson(c, fiber.StatusBadRequest, string(err.Error()))
	}

	if _, exists := refreshTokens[refreshToken]; !exists {
		return utils.RespondJson(c, fiber.StatusBadRequest, "Refresh token is not active")
	}

	refreshTokenExp := GetTokenExpiration(refreshToken)
	currentTime := time.Now().Unix()
	if refreshTokenExp-currentTime <= 900 {
		user := models.User{
			UserId:   authToken.UserId,
			Username: authToken.Username,
		}
		newAccessToken, err := GenerateToken(user)
		if err != nil {
			return utils.RespondJson(c, fiber.StatusBadRequest, "Error while generating new access token")
		}

		RevokeToken(refreshToken)
		activeTokens[newAccessToken.Token] = newAccessToken.TokenExp
		refreshTokens[newAccessToken.RefreshToken] = time.Now().Add(time.Minute * 30).Unix()

		var auth = json.AuthJson{
			UserId:       user.UserId,
			Username:     user.Username,
			Token:        newAccessToken.Token,
			RefreshToken: newAccessToken.RefreshToken,
			TokenExp:     newAccessToken.TokenExp,
		}

		return c.JSON(auth)
	}

	return utils.RespondJson(c, fiber.StatusOK, "Token is Still Active")
}

func GetCurrentUserFromToken(tokenString string) (json.AuthTokenJson, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return json.AuthTokenJson{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return json.AuthTokenJson{}, errors.New("invalid claims format")
	}

	var authJson = json.AuthTokenJson{
		UserId:   0, // Default value in case of error or missing value
		Username: "",
	}

	if userId, ok := claims["user_id"].(float64); ok {
		authJson.UserId = int(userId)
	} else {
		return json.AuthTokenJson{}, errors.New("invalid user_id claim format")
	}

	if username, ok := claims["username"].(string); ok {
		authJson.Username = username
	} else {
		return json.AuthTokenJson{}, errors.New("invalid username claim format")
	}

	if exp, ok := claims["exp"].(float64); ok {
		authJson.TokenExp = int64(exp)
	} else {
		return json.AuthTokenJson{}, errors.New("invalid expiration claim format")
	}

	return authJson, nil
}

func GetTokenExpiration(tokenString string) int64 {
	token, _ := jwt.Parse(tokenString, nil)
	claims := token.Claims.(jwt.MapClaims)
	exp := int64(claims["exp"].(float64))
	return exp
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
