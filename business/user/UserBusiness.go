package user

import (
	"awesomeProject/database"
	userexception "awesomeProject/exception"
	"awesomeProject/models"
	"awesomeProject/service"
	"awesomeProject/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetUsers(c *fiber.Ctx) error {

	db, err := database.Connect()
	if err != nil {
		return err
	}

	offset, err := strconv.Atoi(c.Params("offset"))
	if err != nil {
		return utils.RespondJson(c, fiber.StatusBadRequest, string(err.Error()))
	}

	limit, err := strconv.Atoi(c.Params("limit"))
	if err != nil {
		return utils.RespondJson(c, fiber.StatusBadRequest, string(err.Error()))
	}

	var users []models.User
	err = db.Preload("Books").Preload("UserImages").Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return err
	}

	return c.JSON(users)
}

func GetUserById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.RespondJson(c, fiber.StatusBadRequest, string(err.Error()))
	}
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
		return userexception.EmptyUsername(c)
	}
	if userEmail == "" {
		return userexception.EmptyEmail(c)
	}
	if password == "" {
		return userexception.EmptyPassword(c)
	}
	if address == "" {
		return userexception.EmptyAddress(c)
	}
	if tel == "" {
		return userexception.EmptyTel(c)
	}
	if pId == "" {
		return userexception.EmptyPid(c)
	}
	if image == nil {
		return userexception.EmptyImage(c)
	}

	checkIsUser, err := CheckAlreadyUser(username)
	if err != nil {
		return err
	}
	if checkIsUser.Username == username {
		return userexception.AlreadyUser(c)
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

	if address == "" {
		return userexception.EmptyAddress(c)
	}
	if tel == "" {
		return userexception.EmptyTel(c)
	}
	if pId == "" {
		return userexception.EmptyPid(c)
	}
	if image == nil {
		return userexception.EmptyImage(c)
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
	if tx.RowsAffected == 0 {
		return userexception.NotFoundUser(c)
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
