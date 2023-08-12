package business

import (
	"awesomeProject/database"
	"awesomeProject/models"
	"awesomeProject/utils"
	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	rows, err := database.Select("Users", nil)

	if err != nil {
		utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to connect to the database")
		return err
	}

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UserId, &user.Username, &user.UserEmail, &user.UserStatus)
		if err != nil {
			utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to read user data")
			return err
		}
		users = append(users, user)
	}

	return c.JSON(users)
}

func GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")

	row, err := database.Select("Users", map[string]interface{}{"user_id": id})
	if err != nil {
		utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to connect to the database")
		return err
	}
	var user models.User
	if row.Next() {
		err := row.Scan(&user.UserId, &user.Username, &user.UserEmail, &user.UserStatus)
		if err != nil {
			utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to read user data")
			return err
		}
	} else {
		utils.RespondWithError(c, fiber.StatusBadRequest, "User not found")
		return nil
	}

	return c.JSON(user)
}
