package business

import (
	"awesomeProject/database"
	"awesomeProject/models"
	"awesomeProject/service"
	"awesomeProject/utils"
	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}
	var users []models.User
	db.Find(&users)
	return c.JSON(users)
}

func GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	db, err := database.Connect()
	if err != nil {
		return err
	}
	var user models.User
	db.Find(&user, id)

	return c.JSON(user)
}

func CreateUser(c *fiber.Ctx) error {
	db, err := database.Connect()
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

	dataUsers := models.User{
		Username:  username,
		UserEmail: userEmail,
	}
	db.Create(&dataUsers)

	dataUserDetails := models.Book{
		UserId:  dataUsers.UserId,
		Address: address,
		Tel:     tel,
		PId:     pId,
	}
	db.Create(&dataUserDetails)

	dataUserImage := models.UserImage{
		UserId: dataUsers.UserId,
		Image:  file,
	}
	db.Create(&dataUserImage)

	return utils.RespondJson(c, fiber.StatusCreated, "User created successfully")
}

func UpdateUser(c *fiber.Ctx) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	authJson, err := service.CurrentUser(c)
	if err != nil {
		return err
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	address := form.Value["address"][0]
	tel := form.Value["tel"][0]
	pId := form.Value["pid"][0]
	image, err := c.FormFile("image")
	file, err := service.UploadFile(image)
	if err != nil {
		return err
	}

	dataUserDetails := models.Book{
		UserId:  uint(authJson.UserId),
		Address: address,
		Tel:     tel,
		PId:     pId,
	}
	db.Model(&dataUserDetails).Where("user_id = ?", authJson.UserId).Updates(&dataUserDetails)

	dataUserImage := models.UserImage{
		UserId: uint(authJson.UserId),
		Image:  file,
	}
	db.Model(&dataUserImage).Where("user_id = ?", authJson.UserId).Updates(&dataUserImage)

	return utils.RespondJson(c, fiber.StatusNoContent, "User updated successfully")
}

func DeleteUser(c *fiber.Ctx) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	id := c.Params("id")

	db.Where("user_id", id).Delete(&models.User{})
	db.Where("user_id", id).Delete(&models.Book{})
	db.Where("user_id", id).Delete(&models.UserImage{})

	return utils.RespondJson(c, fiber.StatusNoContent, "User deleted successfully")
}
