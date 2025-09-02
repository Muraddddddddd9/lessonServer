package api

import (
	"fmt"
	"lesson_server/constants"
	db_core "lesson_server/database"
	"lesson_server/handlers/ws"
	"lesson_server/utils"

	"github.com/gofiber/fiber/v2"
)

type AnswerStruct struct {
	ID     string `json:"id"`
	Answer string `json:"answer"`
}

type GetAnswer struct {
	AnswerUser []AnswerStruct `json:"answer"`
}

var FastTest []int

func AddBimCoin(getAnswer GetAnswer, test string) int {
	var BimCoin int
	var question map[string]constants.QuestionStruct

	if test == constants.TestOneName {
		question = constants.Questions
	} else if test == constants.TestTwoName {
		question = constants.QuestionsTeam
	}

	for _, a := range getAnswer.AnswerUser {
		if question[a.ID].AnswerTrue == a.Answer {
			BimCoin += question[a.ID].Socer
		}
	}

	if BimCoin == 0 {
		return BimCoin
	} else {
		BimCoin += 10 - len(FastTest)
	}

	return BimCoin
}

func CheckAnswer(c *fiber.Ctx, db *db_core.DatabaseStruct, test string, time *ws.TimeData) error {
	session := c.Cookies(constants.SessionKey)
	min, sec, flag := time.GetDataTime()
	timeText := fmt.Sprintf("%02d:%02d", *min, *sec)

	if !*flag || (*min == 0 && *sec == 0) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrNotSendTest,
		})
	}

	var getAnswerUser GetAnswer
	if err := c.BodyParser(&getAnswerUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constants.ErrInternalServer,
		})
	}

	if len(getAnswerUser.AnswerUser) != len(constants.Questions) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": constants.ErrNoFullAnser,
		})
	}

	id, err := utils.GetID(session)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	switch test {
	case constants.TestOneName:
		user, err := db.GetOneUser("id = ?", id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": constants.ErrUserNotFound,
			})
		}

		if user.TestFirst {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": constants.ErrAlreadyReplied,
			})
		}

		bimCoin := AddBimCoin(getAnswerUser, test)

		err = db.UpdateData(db_core.TableUsers, "bim_coin = bim_coin + ?", "id = ?", bimCoin, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		err = db.UpdateData(db_core.TableUsers, "test_first = ?", "id = ?", true, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		err = db.UpdateData(db_core.TableUsers, "time_test = ?", "id = ?", timeText, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		FastTest = append(FastTest, id)
	case constants.TestTwoName:
		user, err := db.GetOneUser("id = ?", id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": constants.ErrUserNotFound,
			})
		}

		setting, err := db.GetSetting()
		if err != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		if (user.Team == 1 && setting.TestTeamFirst) || (user.Team == 2 && setting.TestTeamSecond) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": constants.ErrAlreadyReplied,
			})
		}

		bimCoin := AddBimCoin(getAnswerUser, test)

		err = db.UpdateData(db_core.TableUsers, "bim_coin = bim_coin + ?", "team = ?", bimCoin, user.Team)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		var colume string
		switch user.Team {
		case 1:
			colume = "test_team_first = ?"
		case 2:
			colume = "test_team_second = ?"
		}

		err = db.UpdateData(db_core.TableSetting, colume, "", true)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	}

	constants.NewUsers = true

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": constants.SuccGetAnswer,
	})
}

func CheckAnswerOnly(c *fiber.Ctx, db *db_core.DatabaseStruct) error {
	return CheckAnswer(c, db, constants.TestOneName, ws.TimeOnly)
}

func CheckAnswerTeam(c *fiber.Ctx, db *db_core.DatabaseStruct) error {
	return CheckAnswer(c, db, constants.TestTwoName, ws.TimeTeam)
}
