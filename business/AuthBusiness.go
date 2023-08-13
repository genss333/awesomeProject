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
		utils.RespondWithError(c, fiber.StatusBadRequest, "Failed to parse request body")
	}

	if requestBody.Username == "" {
		utils.RespondWithError(c, fiber.StatusBadRequest, "Username is empty")
	}

	if requestBody.UserEmail == "" {
		utils.RespondWithError(c, fiber.StatusBadRequest, "User email is empty")
	}

	rows, err := database.Select("Users", map[string]interface{}{
		"username":   requestBody.Username,
		"user_email": requestBody.UserEmail,
	})
	if err != nil {
		utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to connect to the database")
		return err
	}
	if err != nil {
		utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to fetch users from the database")
		return err
	}

	var user models.User
	for rows.Next() {
		err := rows.Scan(&user.UserId, &user.Username, &user.UserEmail, &user.UserStatus)
		if err != nil {
			utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to read user data")
			return err
		}
	}

	if user.UserId == 0 {
		utils.RespondWithError(c, fiber.StatusBadRequest, "User not found")
		return err
	}

	if user.UserStatus == 0 {
		utils.RespondWithError(c, fiber.StatusUnauthorized, "User is not active")
		return err
	}

	token, err := service.GenerateToken(user)
	if err != nil {
		utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to generate token")
		return err
	}

	return c.JSON(token)

}

func Logout(c *fiber.Ctx) error {
	err := service.Logout(c)
	if err != nil {
		return err
	}
	return utils.RespondWithSuccess(c, fiber.StatusOK, "Logout success")

}
