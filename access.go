package main

import (
	"lesson_server/constants"
	db_core "lesson_server/database"
	"lesson_server/utils"

	"slices"

	"github.com/gofiber/fiber/v2"
)

func AccessToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies(constants.AuthKey)

		id, err := utils.VerifyJWT(session)
		if err != nil {
			utils.DeleteCookie(c)
			c.Cookie(&fiber.Cookie{
				Name:     constants.AuthKey,
				Value:    "",
				MaxAge:   -1,
				HTTPOnly: false,
				Secure:   false,
				SameSite: "Lax",
				Path:     "/",
			})
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message":  constants.ErrUserEntry,
				"redirect": "/",
			})
		}

		if id != constants.TokenCheckAuth {
			utils.DeleteCookie(c)
			c.Cookie(&fiber.Cookie{
				Name:     constants.AuthKey,
				Value:    "",
				MaxAge:   -1,
				HTTPOnly: false,
				Secure:   false,
				SameSite: "Lax",
				Path:     "/",
			})
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message":  constants.ErrUserEntry,
				"redirect": "/",
			})
		}

		return c.Next()
	}
}

func Access(db *db_core.DatabaseStruct, status []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies(constants.SessionKey)
		id, err := utils.VerifyJWT(session)

		if err != nil {
			utils.DeleteCookie(c)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message":  err.Error(),
				"redirect": "/",
			})
		}

		user, err := db.GetOneUser("id = ?", id)

		if err != nil {
			utils.DeleteCookie(c)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message":  err.Error(),
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
