package auth

import (
	"awesomeProject/business/user"
	"awesomeProject/database"
	userexception "awesomeProject/exception"
	"awesomeProject/models"
	"awesomeProject/payload/auth"
	"awesomeProject/service"
	"awesomeProject/utils"
	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	var requestBody auth.AuthPayload
	if err := c.BodyParser(&requestBody); err != nil {
		return utils.RespondJson(c, fiber.StatusBadRequest, "Failed to parse request body")
	}

	if requestBody.Username == "" {
		return utils.RespondJson(c, fiber.StatusBadRequest, "Username is empty")
	}

	if requestBody.Password == "" {
		return utils.RespondJson(c, fiber.StatusBadRequest, "Password email is empty")
	}

	db, err := database.Connect()
	if err != nil {
		return utils.RespondJson(c, fiber.StatusInternalServerError, string(err.Error()))
	}
	var user models.User
	db.Find(&user, "username = ?", requestBody.Username)

	if user.Password != utils.EnSha256Hash(requestBody.Password) {
		return utils.RespondJson(c, fiber.StatusBadRequest, "Email or Password is incorrect")
	}

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
		return utils.RespondJson(c, fiber.StatusOK, "Logout failed")
	}
	return utils.RespondJson(c, fiber.StatusOK, "Logout success")

}

func Register(c *fiber.Ctx) error {
	var requestBody auth.RegisterPayload
	if err := c.BodyParser(&requestBody); err != nil {
		return utils.RespondJson(c, fiber.StatusBadRequest, "Failed to parse request body")
	}
	if requestBody.Email == "" {
		return utils.RespondJson(c, fiber.StatusBadRequest, "Email is empty")
	}

	if requestBody.Username == "" {
		return utils.RespondJson(c, fiber.StatusBadRequest, "Username is empty")
	}

	if requestBody.Password == "" {
		return utils.RespondJson(c, fiber.StatusBadRequest, "Password is empty")
	}

	if requestBody.ConfirmPassword == "" {
		return utils.RespondJson(c, fiber.StatusBadRequest, "Confirm Password is empty")
	}

	if requestBody.Password != requestBody.ConfirmPassword {
		return utils.RespondJson(c, fiber.StatusBadRequest, "Confirm Password is not match")
	}

	db, err := database.Connect()
	if err != nil {
		return err
	}

	checkIsUser, err := user.CheckAlreadyUser(requestBody.Username)
	if err != nil {
		return err
	}
	if checkIsUser.Username == requestBody.Username {
		return userexception.AlreadyUser(c)
	}

	tx := db.Begin()

	dataUser := models.User{
		Username:  requestBody.Username,
		UserEmail: requestBody.Email,
		Password:  utils.EnSha256Hash(requestBody.Password),
	}
	tx.Create(&dataUser)
	tx.Commit()

	return utils.RespondJson(c, fiber.StatusCreated, "Register has successfully")
}
