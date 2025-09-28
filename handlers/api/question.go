package api

import (
	"encoding/json"
	"fmt"
	"lesson_server/constants"
	db_core "lesson_server/database"
	"lesson_server/utils"

	"github.com/gofiber/fiber/v2"
)

type SendQuestionsStruct struct {
	ID       string   `json:"id"`
	Question string   `json:"question"`
	Answers  []string `json:"answers"`
	Socer    uint64   `json:"socer"`
}

func GetQuestionsOnly(c *fiber.Ctx, db *db_core.DatabaseStruct) error {
	groupUser := c.Cookies(constants.SessionKey)
	id, err := utils.VerifyJWT(groupUser)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user, err := db.GetOneUser("id = ?", id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	keyQuestion := fmt.Sprintf("Questions_%d", user.Team)
	question := constants.Questions[keyQuestion]
	
	questionsJSON, err := json.Marshal(question)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constants.ErrInternalServer,
		})
	}

	var sendQuestions SendQuestionsStruct
	err = json.Unmarshal(questionsJSON, &sendQuestions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constants.ErrInternalServer,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"questions": sendQuestions,
	})
}
