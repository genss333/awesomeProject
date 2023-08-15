package userexception

import (
	"awesomeProject/utils"
	"github.com/gofiber/fiber/v2"
)

func NotFoundUser(c *fiber.Ctx) error {
	return utils.RespondJson(c, fiber.StatusNotFound, "Not found user")
}

func AlreadyUser(c *fiber.Ctx) error {
	return utils.RespondJson(c, fiber.StatusBadRequest, "Username is already")
}

func EmptyUsername(c *fiber.Ctx) error {
	return utils.RespondJson(c, fiber.StatusBadRequest, "Username is empty")
}

func EmptyEmail(c *fiber.Ctx) error {
	return utils.RespondJson(c, fiber.StatusBadRequest, "Email is empty")
}

func EmptyPassword(c *fiber.Ctx) error {
	return utils.RespondJson(c, fiber.StatusBadRequest, "Password is empty")
}

func EmptyAddress(c *fiber.Ctx) error {
	return utils.RespondJson(c, fiber.StatusBadRequest, "Address is empty")
}

func EmptyTel(c *fiber.Ctx) error {
	return utils.RespondJson(c, fiber.StatusBadRequest, "Tel is empty")
}

func EmptyPid(c *fiber.Ctx) error {
	return utils.RespondJson(c, fiber.StatusBadRequest, "Pid is empty")
}

func EmptyImage(c *fiber.Ctx) error {
	return utils.RespondJson(c, fiber.StatusBadRequest, "Image is empty")
}
