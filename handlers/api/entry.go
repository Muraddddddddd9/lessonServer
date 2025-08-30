package api

import (
	"fmt"
	"lesson_server/constants"
	db_core "lesson_server/database"
	"lesson_server/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
)

type EntryStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Team     string `json:"team"`
}

func Entry(c *fiber.Ctx, db *db_core.DatabaseStruct) error {
	var entryData EntryStruct
	if err := c.BodyParser(&entryData); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInputValue,
		})
	}

	username := strings.TrimSpace(entryData.Username)
	password := strings.TrimSpace(entryData.Password)
	team, _ := strconv.Atoi(strings.TrimSpace(entryData.Team))

	var userID int64
	var userStatus string

	user, err := db.GetOneUser("name = ?", username)
	if err != nil {
		newUser := db_core.UserStruct{
			Name:     username,
			Password: password,
			Status:   constants.StudentStatus,
			Team:     team,
		}

		userID, err = db.InsertUser(newUser)
		if err != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": err,
			})
		}
		userStatus = constants.StudentStatus
	} else {
		if user.Password != password {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": constants.ErrUserNotFound,
			})
		} else {
			userID = int64(user.ID)
			userStatus = user.Status
		}
	}

	constants.NewUsers = true
	sessionID := fmt.Sprintf("%v:%v", userID, encryptcookie.GenerateKey())
	utils.AddCookie(c, sessionID, userStatus)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message":  constants.SuccEntry,
		"redirect": "/lesson",
	})
}
