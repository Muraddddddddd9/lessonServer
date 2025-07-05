package api

import (
	"lesson_server/constants"
	db_core "lesson_server/database"
	"lesson_server/handlers/ws"

	"github.com/gofiber/fiber/v2"
)

func ClearData(c *fiber.Ctx, db *db_core.DatabaseStruct) error {
	err := db.DeleteStudent()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resetSetting := map[string]any{
		"now_stage_lesson = ?": "0",
		"test_team_first = ?":  false,
		"test_team_second = ?": false,
	}

	for v := range resetSetting {
		err := db.UpdateData(db_core.TableSetting, v, "", resetSetting[v])
		if err != nil {
			continue
		}
	}

	constants.NewLesson = true
	constants.NewUsers = true

	minutCls, secondCls, flagCls := ws.TimeLesson.GetDataTime()
	*minutCls = 45
	*secondCls = 0
	*flagCls = false

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": constants.SuccClearData,
	})
}
