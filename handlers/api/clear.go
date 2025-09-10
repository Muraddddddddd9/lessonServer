package api

import (
	"bytes"
	"fmt"
	"lesson_server/constants"
	db_core "lesson_server/database"
	"lesson_server/handlers/ws"
	"lesson_server/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/olekukonko/tablewriter"
)

func getUsersTableString(students []db_core.UserStruct) string {
	var buf bytes.Buffer

	table := tablewriter.NewWriter(&buf)
	table.Header([]string{"ID", "Name", "Password", "Status", "BimCoin", "Team", "TestFirst", "TimeTest"})
	for _, student := range students {
		table.Append([]any{
			student.ID,
			student.Name,
			student.Password,
			student.Status,
			student.BimCoin,
			student.Team,
			student.TestFirst,
			student.TimeTest,
		})
	}

	table.Render()
	return buf.String()
}

func ClearData(c *fiber.Ctx, db *db_core.DatabaseStruct) error {
	pathLogg, methodLogg, ipLogg := c.Path(), c.Method(), c.IP()

	students, _ := db.GetDataUsers()

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
	FastTest = []int{}

	utils.LogginAPI(pathLogg, methodLogg, fiber.StatusAccepted, ipLogg,
		fmt.Sprintf("\n%s\n", getUsersTableString(students)),
		constants.SuccClearData,
	)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": constants.SuccClearData,
	})
}

func ClearStageLesson(c *fiber.Ctx, db *db_core.DatabaseStruct) error {
	err := db.UpdateData(db_core.TableSetting, "now_stage_lesson = ?", "", "0")
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": constants.ErrClearStageLesson,
		})
	}

	constants.NewLesson = true

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": constants.SuccCleatStageLesson,
	})
}

type IDStudetn struct {
	ID string `json:"id"`
}

func DeleteStudent(c *fiber.Ctx, db *db_core.DatabaseStruct) error {
	var id IDStudetn

	if err := c.BodyParser(&id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInputValue,
		})
	}

	err := db.DeleteUserByID(id.ID)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	constants.NewUsers = true

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": constants.SuccStudentDeleteById,
	})
}
