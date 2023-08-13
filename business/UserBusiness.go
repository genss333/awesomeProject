package business

import (
	"awesomeProject/database"
	"awesomeProject/models"
	CreateUserPayload "awesomeProject/payload"
	"awesomeProject/service"
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

func CreateUser(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	username := form.Value["username"][0]
	userEmail := form.Value["email"][0]
	address := form.Value["address"][0]
	tel := form.Value["tel"][0]
	pId := form.Value["pid"][0]
	image, err := c.FormFile("image")
	file, err := service.UploadFile(image)
	if err != nil {
		return err
	}

	payload := CreateUserPayload.CreateUserPayload{
		Username:  username,
		UserEmail: userEmail,
		Address:   address,
		Tel:       tel,
		PId:       pId,
		Image:     file,
	}

	dataUsers := map[string]interface{}{
		"username":   payload.Username,
		"user_email": payload.UserEmail,
	}

	errUser := database.Insert("Users", dataUsers)
	if errUser != nil {
		utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to insert user data")
		return errUser
	}

	userIdRow, err := database.CustomQuery("SELECT user_id FROM Users ORDER BY user_id DESC LIMIT 1")
	if err != nil {
		utils.RespondWithError(c, fiber.StatusInternalServerError, string(err.Error()))
		return err
	}
	var userId string
	if userIdRow.Next() {
		err := userIdRow.Scan(&userId)
		if err != nil {
			utils.RespondWithError(c, fiber.StatusInternalServerError, string(err.Error()))
			return err
		}
	}

	errBook := database.Insert("Book", map[string]interface{}{
		"address": payload.Address,
		"tel":     payload.Tel,
		"pid":     payload.PId,
	})
	if errBook != nil {
		utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to insert book data")
		return errBook
	}

	errImage := database.Insert("User_Image", map[string]interface{}{
		"image":   payload.Image,
		"user_id": userId,
	})
	if errImage != nil {
		utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to insert image data")
		return errImage
	}

	return utils.RespondWithSuccess(c, fiber.StatusCreated, "User created successfully")
}

func UpdateUser(c *fiber.Ctx) error {
	authJson, err := service.CurrentUser(c)
	if err != nil {
		return err
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	username := form.Value["username"][0]
	userEmail := form.Value["email"][0]
	address := form.Value["address"][0]
	tel := form.Value["tel"][0]
	pId := form.Value["pid"][0]
	image, err := c.FormFile("image")
	file, err := service.UploadFile(image)
	if err != nil {
		return err
	}

	payload := CreateUserPayload.CreateUserPayload{
		Username:  username,
		UserEmail: userEmail,
		Address:   address,
		Tel:       tel,
		PId:       pId,
		Image:     file,
	}

	dataUsers := map[string]interface{}{
		"username":   payload.Username,
		"user_email": payload.UserEmail,
	}

	errUser := database.Update("Users", dataUsers, map[string]interface{}{"user_id": authJson.UserId})
	if errUser != nil {
		utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to update user data")
		return errUser
	}

	errBook := database.Update("Book", map[string]interface{}{
		"address": payload.Address,
		"tel":     payload.Tel,
		"pid":     payload.PId,
	}, map[string]interface{}{"user_id": authJson.UserId})
	if errBook != nil {
		utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to update book data")
		return errBook
	}

	errImage := database.Update("User_Image", map[string]interface{}{
		"image": payload.Image,
	}, map[string]interface{}{"user_id": authJson.UserId})
	if errImage != nil {
		utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to update image data")
		return errImage
	}

	return utils.RespondWithSuccess(c, fiber.StatusNoContent, "User updated successfully")
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	err := database.Delete("Users", map[string]interface{}{"user_id": id})
	if err != nil {
		utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to delete user data")
		return err
	}

	errBook := database.Delete("Book", map[string]interface{}{"user_id": id})
	if errBook != nil {
		utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to delete book data")
		return errBook
	}

	errImage := database.Delete("User_Image", map[string]interface{}{"user_id": id})
	if errImage != nil {
		utils.RespondWithError(c, fiber.StatusInternalServerError, "Failed to delete image data")
		return errImage
	}

	return utils.RespondWithSuccess(c, fiber.StatusNoContent, "User deleted successfully")
}
