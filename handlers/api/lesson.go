package api

import (
	"fmt"
	"lesson_server/constants"
	db_core "lesson_server/database"
	"lesson_server/handlers/ws"

	"github.com/gofiber/fiber/v2"
)

type ActionLesson struct {
	Action int `json:"action"`
}

func LessonChange(c *fiber.Ctx, db *db_core.DatabaseStruct) error {
	var action ActionLesson
	if err := c.BodyParser(&action); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInputValue,
		})
	}

	err := db.UpdateData(db_core.TableSetting, "now_stage_lesson = ?", "", fmt.Sprintf("%d", action.Action))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	_, _, flagOnly := ws.TimeOnly.GetDataTime()
	if *flagOnly {
		*flagOnly = false
	}

	_, _, flagTeam := ws.TimeTeam.GetDataTime()
	if *flagTeam {
		*flagTeam = false
	}

	constants.NewLesson = true

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": constants.SuccUpdateLessonStage,
	})
}
