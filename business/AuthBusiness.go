package business

import (
	"awesomeProject/database"
	"awesomeProject/models"
	"awesomeProject/payload"
	"awesomeProject/service"
	"awesomeProject/utils"
	"database/sql"
	"fmt"
	fiber "github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	db, err := database.Connect()
	if err != nil {
		utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to connect to the database")
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			utils.RespondWithError(c, fiber.StatusInternalServerError, string(err.Error()))
		}
	}(db)

	var requestBody payload.AuthPayload
	if err := c.BodyParser(&requestBody); err != nil {
		utils.RespondWithError(c, fiber.StatusBadRequest, "Failed to parse request body")
		return err
	}

	if requestBody.Username == "" {
		utils.RespondWithError(c, fiber.StatusBadRequest, "Username is empty")
		return err
	}

	if requestBody.UserEmail == "" {
		utils.RespondWithError(c, fiber.StatusBadRequest, "User email is empty")
		return err
	}

	rows, err := db.Query("SELECT * FROM Users WHERE username = ? AND user_email=?", requestBody.Username, requestBody.UserEmail)
	if err != nil {
		utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to fetch users from the database")
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(rows)

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
