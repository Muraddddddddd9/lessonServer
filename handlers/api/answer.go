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

	switch test {
	case constants.TestOneName:
		question = constants.Questions
	case constants.TestTwoName:
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

func calculateTimeDifference(minuteStart, secondStart, minuteEnd, secondEnd int) (int, int) {
	startTotal := minuteStart*60 + secondStart
	endTotal := minuteEnd*60 + secondEnd

	diffSeconds := endTotal - startTotal
	if diffSeconds < 0 {
		diffSeconds = -diffSeconds
	}

	return diffSeconds / 60, diffSeconds % 60
}

func CheckAnswer(c *fiber.Ctx, db *db_core.DatabaseStruct, test string, time *ws.TimeData) error {
	session := c.Cookies(constants.SessionKey)
	min, sec, flag := time.GetDataTime()
	minStart, secStart, _ := ws.TimeOnlyStart.GetDataTime()
	diffMin, diffSec := calculateTimeDifference(*minStart, *secStart, *min, *sec)
	timeText := fmt.Sprintf("%02d:%02d", diffMin, diffSec)

	pathLogg, methodLogg, ipLogg := c.Path(), c.Method(), c.IP()

	if !*flag || (*min == 0 && *sec == 0) {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusBadRequest, ipLogg,
			struct {
				session     string
				timeTestNow string
				timeTest    string
			}{
				session:     session,
				timeTestNow: fmt.Sprintf("min: %d, sec: %d, flag: %t", *min, *sec, *flag),
				timeTest:    fmt.Sprintf("minStart: %d, secStart: %d", *min, *sec),
			},
			constants.ErrNotSendTest,
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrNotSendTest,
		})
	}

	var getAnswerUser GetAnswer
	if err := c.BodyParser(&getAnswerUser); err != nil {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusInternalServerError, ipLogg,
			struct {
				session     string
				timeTestNow string
				timeTest    string
				body        any
			}{
				session:     session,
				timeTestNow: fmt.Sprintf("min: %d, sec: %d, flag: %t", *min, *sec, *flag),
				timeTest:    fmt.Sprintf("minStart: %d, secStart: %d", *min, *sec),
				body:        getAnswerUser,
			},
			constants.ErrInternalServer,
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constants.ErrInternalServer,
		})
	}

	if len(getAnswerUser.AnswerUser) != len(constants.Questions) {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusConflict, ipLogg,
			struct {
				session     string
				timeTestNow string
				timeTest    string
				body        any
			}{
				session:     session,
				timeTestNow: fmt.Sprintf("min: %d, sec: %d, flag: %t", *min, *sec, *flag),
				timeTest:    fmt.Sprintf("minStart: %d, secStart: %d", *min, *sec),
				body:        getAnswerUser,
			},
			constants.ErrNoFullAnser,
		)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": constants.ErrNoFullAnser,
		})
	}

	id, err := utils.GetID(session)
	if err != nil {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusInternalServerError, ipLogg,
			struct {
				session     string
				timeTestNow string
				timeTest    string
				body        any
			}{
				session:     session,
				timeTestNow: fmt.Sprintf("min: %d, sec: %d, flag: %t", *min, *sec, *flag),
				timeTest:    fmt.Sprintf("minStart: %d, secStart: %d", *min, *sec),
				body:        getAnswerUser,
			},
			err.Error(),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	switch test {
	case constants.TestOneName:
		user, err := db.GetOneUser("id = ?", id)
		if err != nil {
			utils.LogginAPI(pathLogg, methodLogg, fiber.StatusNotFound, ipLogg,
				struct {
					session     string
					timeTestNow string
					timeTest    string
					body        any
				}{
					session:     session,
					timeTestNow: fmt.Sprintf("min: %d, sec: %d, flag: %t", *min, *sec, *flag),
					timeTest:    fmt.Sprintf("minStart: %d, secStart: %d", *min, *sec),
					body:        getAnswerUser,
				},
				constants.ErrUserNotFound,
			)
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
			utils.LogginAPI(pathLogg, methodLogg, fiber.StatusInternalServerError, ipLogg,
				struct {
					session     string
					timeTestNow string
					timeTest    string
					id          int
					bim_coin    int
					body        any
				}{
					session:     session,
					timeTestNow: fmt.Sprintf("min: %d, sec: %d, flag: %t", *min, *sec, *flag),
					timeTest:    fmt.Sprintf("minStart: %d, secStart: %d", *min, *sec),
					id:          id,
					bim_coin:    bimCoin,
					body:        getAnswerUser,
				},
				err.Error(),
			)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		err = db.UpdateData(db_core.TableUsers, "test_first = ?", "id = ?", true, id)
		if err != nil {
			utils.LogginAPI(pathLogg, methodLogg, fiber.StatusInternalServerError, ipLogg,
				struct {
					session     string
					timeTestNow string
					timeTest    string
					id          int
					bim_coin    int
					body        any
				}{
					session:     session,
					timeTestNow: fmt.Sprintf("min: %d, sec: %d, flag: %t", *min, *sec, *flag),
					timeTest:    fmt.Sprintf("minStart: %d, secStart: %d", *min, *sec),
					id:          id,
					bim_coin:    bimCoin,
					body:        getAnswerUser,
				},
				err.Error(),
			)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		err = db.UpdateData(db_core.TableUsers, "time_test = ?", "id = ?", timeText, id)
		if err != nil {
			utils.LogginAPI(pathLogg, methodLogg, fiber.StatusInternalServerError, ipLogg,
				struct {
					session     string
					timeTestNow string
					timeTest    string
					id          int
					timeText    string
					body        any
				}{
					session:     session,
					timeTestNow: fmt.Sprintf("min: %d, sec: %d, flag: %t", *min, *sec, *flag),
					timeTest:    fmt.Sprintf("minStart: %d, secStart: %d", *min, *sec),
					id:          id,
					timeText:    timeText,
					body:        getAnswerUser,
				},
				err.Error(),
			)
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

	utils.LogginAPI(pathLogg, methodLogg, fiber.StatusAccepted, ipLogg,
		struct {
			session     string
			timeTestNow string
			timeTest    string
			id          int
			timeText    string
			body        any
		}{
			session:     session,
			timeTestNow: fmt.Sprintf("min: %d, sec: %d, flag: %t", *min, *sec, *flag),
			timeTest:    fmt.Sprintf("minStart: %d, secStart: %d", *min, *sec),
			id:          id,
			timeText:    timeText,
			body:        getAnswerUser,
		},
		constants.SuccGetAnswer,
	)
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
