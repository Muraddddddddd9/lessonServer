package api

import (
	"lesson_server/constants"
	db_core "lesson_server/database"
	"lesson_server/utils"

	"github.com/gofiber/fiber/v2"
)

func Exit(c *fiber.Ctx, db *db_core.DatabaseStruct) error {
	pathLogg, methodLogg, ipLogg := c.Path(), c.Method(), c.IP()

	session := c.Cookies(constants.SessionKey)
	if session == "" {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusBadRequest, ipLogg, session, constants.ErrUserExit)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":  constants.ErrUserExit,
			"redirect": "/",
		})
	}

	utils.DeleteCookie(c)

	utils.LogginAPI(pathLogg, methodLogg, fiber.StatusAccepted, ipLogg, session, constants.SuccExit)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message":  constants.SuccExit,
		"redirect": "/",
	})
}
