package api

import (
	"lesson_server/constants"
	db_core "lesson_server/database"
	"lesson_server/utils"

	"github.com/gofiber/fiber/v2"
)

func Exit(c *fiber.Ctx, db *db_core.DatabaseStruct) error {
	session := c.Cookies(constants.SessionKey)
	if session == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":  constants.ErrUserExit,
			"redirect": "/",
		})
	}

	utils.DeleteCookie(c)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message":  constants.SuccExit,
		"redirect": "/",
	})
}
