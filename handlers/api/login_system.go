package api

import (
	"fmt"
	"lesson_server/constants"
	db_core "lesson_server/database"
	"lesson_server/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type EntrySystem struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginSystem(c *fiber.Ctx, db *db_core.DatabaseStruct) error {
	var entrySystem EntrySystem
	if err := c.BodyParser(&entrySystem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInputValue,
		})
	}

	entrySystem.Username = strings.TrimSpace(entrySystem.Username)
	entrySystem.Password = strings.TrimSpace(entrySystem.Password)

	user, err := db.GetSuperUser(entrySystem.Username)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrUserNotFound,
		})
	}

	if user.Password != entrySystem.Password {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrUserNotFound,
		})
	}

	value, err := utils.GenerateJWT(fmt.Sprintf("%d", constants.TokenCheckAuth))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constants.ErrInternalServer,
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     constants.AuthKey,
		Value:    value,
		Expires:  time.Now().Add(48 * time.Hour),
		HTTPOnly: false,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": constants.SuccEntry,
	})
}
