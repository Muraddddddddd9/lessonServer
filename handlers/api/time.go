package api

import (
	"lesson_server/constants"
	db_core "lesson_server/database"
	"lesson_server/handlers/ws"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ChangeTimeStruct struct {
	ChageTime string `json:"change_time"`
}

type RedactTimeStuct struct {
	ChageTime string `json:"change_time"`
	NewTime   string `json:"new_time"`
}

func ChangeTime(c *fiber.Ctx, db *db_core.DatabaseStruct) error {
	var changeTime ChangeTimeStruct
	if err := c.BodyParser(&changeTime); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInputValue,
		})
	}

	switch changeTime.ChageTime {
	case "lesson":
		_, _, flag := ws.TimeLesson.GetDataTime()
		ws.TimeLesson.Flag = !*flag
	case constants.TestOneName:
		_, _, flag := ws.TimeOnly.GetDataTime()
		ws.TimeOnly.Flag = !*flag
	case constants.TestTwoName:
		_, _, flag := ws.TimeTeam.GetDataTime()
		ws.TimeTeam.Flag = !*flag
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": constants.SuccChangeTime,
	})
}

func EditTime(c *fiber.Ctx, db *db_core.DatabaseStruct) error {
	var redactTime RedactTimeStuct
	if err := c.BodyParser(&redactTime); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInputValue,
		})
	}

	timeSplit := strings.Split(redactTime.NewTime, ":")
	if len(timeSplit) != 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInputValue,
		})
	}

	minute, err := strconv.Atoi(timeSplit[0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInputValue,
		})
	}

	second, err := strconv.Atoi(timeSplit[1])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInputValue,
		})
	}

	switch redactTime.ChageTime {
	case constants.TestOneName:
		ws.TimeOnly.Flag = false
		ws.TimeOnly.Minute = minute
		ws.TimeOnly.Second = second

		ws.TimeOnlyStart.Minute = minute
		ws.TimeOnlyStart.Second = second
	case constants.TestTwoName:
		ws.TimeTeam.Flag = false
		ws.TimeTeam.Minute = minute
		ws.TimeTeam.Second = second

	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": constants.SuccChangeTime,
	})
}
