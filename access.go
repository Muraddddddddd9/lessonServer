package main

import (
	"lesson_server/constants"
	db_core "lesson_server/database"
	"lesson_server/utils"

	"slices"

	"github.com/gofiber/fiber/v2"
)

func Access(db *db_core.DatabaseStruct, status []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies(constants.SessionKey)

		id, err := utils.GetID(session)
		if err != nil {
			utils.DeleteCookie(c)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message":  err,
				"redirect": "/",
			})
		}

		user, err := db.GetOneUser("id = ?", id)
		if err != nil {
			utils.DeleteCookie(c)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message":  err,
				"redirect": "/",
			})
		}

		if slices.Contains(status, user.Status) {
			return c.Next()
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constants.ErrInternalServer,
		})
	}
}
