package business

import (
	"awesomeProject/database"
	"awesomeProject/models"
	"awesomeProject/service"
	"awesomeProject/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	var users []models.User
	db.Preload("Books").Preload("UserImages").Find(&users)
	fmt.Println(users)

	return c.JSON(users)
}

func GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	db, err := database.Connect()
	if err != nil {
		return err
	}

	var user models.User
	db.Preload("Books").Preload("UserImages").Find(&user, id)

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
	password := form.Value["password"][0]
	address := form.Value["address"][0]
	tel := form.Value["tel"][0]
	pId := form.Value["pid"][0]
	image, err := c.FormFile("image")
	if err != nil {
		return err
	}

	if username == "" {
		return utils.RespondJson(c, fiber.StatusBadRequest, "Username is required")
	}
	if userEmail == "" {
		return utils.RespondJson(c, fiber.StatusBadRequest, "Email is required")
	}
	if password == "" {
		return utils.RespondJson(c, fiber.StatusBadRequest, "Password is required")
	}
	if address == "" {
		return utils.RespondJson(c, fiber.StatusBadRequest, "Address is required")
	}
	if tel == "" {
		return utils.RespondJson(c, fiber.StatusBadRequest, "Tel is required")
	}
	if pId == "" {
		return utils.RespondJson(c, fiber.StatusBadRequest, "Pid is required")
	}
	if image == nil {
		return utils.RespondJson(c, fiber.StatusBadRequest, "Image is required")
	}

	checkIsUser, err := CheckAlreadyUser(username)
	if err != nil {
		return utils.RespondJson(c, fiber.StatusBadRequest, err.Error())
	}
	if checkIsUser.Username == username {
		return utils.RespondJson(c, fiber.StatusBadRequest, "User already exists")
	}

	tx := db.Begin()

	dataUser := models.User{
		Username:  username,
		UserEmail: userEmail,
		Password:  utils.EnSha256Hash(password),
	}
	tx.Create(&dataUser)

	dataBook := models.Book{
		Address: address,
		Tel:     tel,
		PId:     pId,
		UserID:  dataUser.UserId,
	}
	tx.Create(&dataBook)

	file, err := service.UploadFile(image)
	if err != nil {
		tx.Rollback()
		return err
	}
	dataUserImage := models.UserImage{
		Image:  file,
		UserID: dataUser.UserId,
	}
	tx.Create(&dataUserImage)

	tx.Commit()

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
	if err != nil {
		return err
	}

	tx := db.Begin()

	tx.Model(&models.Book{}).Where("user_id = ?", authJson.UserId).
		Updates(models.Book{
			Address: address,
			Tel:     tel,
			PId:     pId,
		})

	if image != nil {
		file, err := service.UploadFile(image)
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Model(&models.UserImage{}).Where("user_id = ?", authJson.UserId).
			Updates(models.UserImage{
				Image: file,
			})
	}

	tx.Commit()

	return utils.RespondJson(c, fiber.StatusNoContent, "User updated successfully")
}

func DeleteUser(c *fiber.Ctx) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	id := c.Params("id")

	tx := db.Begin()

	tx.Find(&models.User{}, "user_id = ?", id)
	fmt.Println(tx.RowsAffected)
	if tx.RowsAffected == 0 {
		return utils.RespondJson(c, fiber.StatusBadRequest, "User not found")
	}
	tx.Delete(&models.User{}, id)

	tx.Commit()

	return utils.RespondJson(c, fiber.StatusNoContent, "User deleted successfully")
}

func CheckAlreadyUser(username string) (models.User, error) {
	db, err := database.Connect()
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	db.Find(&user, "username = ?", username)

	return user, nil

}
