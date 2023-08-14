package business

import (
	"awesomeProject/database"
	"awesomeProject/models"
	"awesomeProject/payload"
	"awesomeProject/service"
	"awesomeProject/utils"
	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	var requestBody payload.AuthPayload
	if err := c.BodyParser(&requestBody); err != nil {
		return utils.RespondJson(c, fiber.StatusBadRequest, "Failed to parse request body")
	}

	if requestBody.Username == "" {
		return utils.RespondJson(c, fiber.StatusBadRequest, "Username is empty")
	}

	if requestBody.UserEmail == "" {
		return utils.RespondJson(c, fiber.StatusBadRequest, "User email is empty")
	}

	db, err := database.Connect()
	if err != nil {
		return utils.RespondJson(c, fiber.StatusInternalServerError, string(err.Error()))
	}
	var user models.User
	db.Find(&user, "user_email = ? AND username = ?", requestBody.UserEmail, requestBody.Username)

	if err != nil {
		return utils.RespondJson(c, fiber.StatusInternalServerError, "Failed to connect to the database")
	}
	if err != nil {
		return utils.RespondJson(c, fiber.StatusInternalServerError, "Failed to fetch users from the database")
	}

	if user.UserId == 0 {
		return utils.RespondJson(c, fiber.StatusBadRequest, "User not found")
	}

	if user.UserStatus == 0 {
		return utils.RespondJson(c, fiber.StatusUnauthorized, "User is not active")
	}

	token, err := service.GenerateToken(user)
	if err != nil {
		return utils.RespondJson(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	return c.JSON(token)

}

func Logout(c *fiber.Ctx) error {
	err := service.Logout(c)
	if err != nil {
		return err
	}
	return utils.RespondJson(c, fiber.StatusOK, "Logout success")

}
