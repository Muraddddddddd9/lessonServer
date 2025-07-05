package api

import (
	"encoding/json"
	"lesson_server/constants"

	"github.com/gofiber/fiber/v2"
)

type SendQuestionsStruct struct {
	ID       string   `json:"id"`
	Question string   `json:"question"`
	Answers  []string `json:"answers"`
	Socer    uint64   `json:"socer"`
}

func GetQuestions(c *fiber.Ctx, qst interface{}) error {
	questionsJSON, err := json.Marshal(qst)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constants.ErrInternalServer,
		})
	}

	var sendQuestions []SendQuestionsStruct
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

func GetQuestionsOnly(c *fiber.Ctx) error {
	return GetQuestions(c, constants.Questions)
}
func GetQuestionsTeam(c *fiber.Ctx) error {
	return GetQuestions(c, constants.QuestionsTeam)
}
