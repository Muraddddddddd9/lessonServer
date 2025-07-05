package api

import (
	"lesson_server/constants"
	db_core "lesson_server/database"

	"github.com/gofiber/fiber/v2"
)

type RedactPresentationSturct struct {
	ID string `json:"id"`
}

func EditPresentation(c *fiber.Ctx, db *db_core.DatabaseStruct) error {
	var dataPresentation RedactPresentationSturct
	if err := c.BodyParser(&dataPresentation); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInputValue,
		})
	}

	err := db.UpdateData(db_core.TableSetting, "id_presentation = ?", "", dataPresentation.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": constants.SuccUpdateIdPresentation,
	})
}

func GetPresentation(c *fiber.Ctx, db *db_core.DatabaseStruct) error {
	setting, err := db.GetSetting()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"id": setting.IdPresentation,
	})
}
