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

func AddBimCoin(getAnswer GetAnswer) int {
	var BimCoin int
	var question = constants.Questions

	for _, a := range getAnswer.AnswerUser {
		if question[a.ID].AnswerTrue == a.Answer {
			BimCoin += question[a.ID].Socer
		}
	}

	if BimCoin == 0 {
		return BimCoin
	} else if len(FastTest) < 3 {
		BimCoin += 3 - len(FastTest)
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

func CheckAnswerOnly(c *fiber.Ctx, db *db_core.DatabaseStruct) error {
	session := c.Cookies(constants.SessionKey)
	min, sec, flag := ws.TimeOnly.GetDataTime()
	minStart, secStart, _ := ws.TimeOnlyStart.GetDataTime()
	diffMin, diffSec := calculateTimeDifference(*minStart, *secStart, *min, *sec)
	timeText := fmt.Sprintf("%02d:%02d", diffMin, diffSec)

	pathLogg, methodLogg, ipLogg := c.Path(), c.Method(), c.IP()

	if !*flag || (*min == 0 && *sec == 0) {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusBadRequest, ipLogg,
			fmt.Sprintf("Session: %s, TimeTestNow: %s, TimeTest: %s",
				session,
				fmt.Sprintf("min: %d, sec: %d, flag: %t", *min, *sec, *flag),
				fmt.Sprintf("minStart: %d, secStart: %d", *min, *sec),
			),
			constants.ErrNotSendTest,
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrNotSendTest,
		})
	}

	var getAnswerUser GetAnswer
	if err := c.BodyParser(&getAnswerUser); err != nil {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusInternalServerError, ipLogg,
			fmt.Sprintf("Session: %s, TimeTestNow: %s, TimeTest: %s, Body: %v",
				session,
				fmt.Sprintf("min: %d, sec: %d, flag: %t", *min, *sec, *flag),
				fmt.Sprintf("minStart: %d, secStart: %d", *min, *sec),
				getAnswerUser,
			),
			constants.ErrInternalServer,
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constants.ErrInternalServer,
		})
	}

	if len(getAnswerUser.AnswerUser) != 1 {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusConflict, ipLogg,
			fmt.Sprintf("Session: %s, TimeTestNow: %s, TimeTest: %s, Body: %v",
				session,
				fmt.Sprintf("min: %d, sec: %d, flag: %t", *min, *sec, *flag),
				fmt.Sprintf("minStart: %d, secStart: %d", *min, *sec),
				getAnswerUser,
			),
			constants.ErrNoFullAnswer,
		)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": constants.ErrNoFullAnswer,
		})
	}

	id, err := utils.VerifyJWT(session)
	if err != nil {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusInternalServerError, ipLogg,
			fmt.Sprintf("Session: %s, TimeTestNow: %s, TimeTest: %s, Body: %v",
				session,
				fmt.Sprintf("min: %d, sec: %d, flag: %t", *min, *sec, *flag),
				fmt.Sprintf("minStart: %d, secStart: %d", *min, *sec),
				getAnswerUser,
			),
			err.Error(),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user, err := db.GetOneUser("id = ?", id)
	if err != nil {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusNotFound, ipLogg,
			fmt.Sprintf("Session: %s, TimeTestNow: %s, TimeTest: %s, Body: %v",
				session,
				fmt.Sprintf("min: %d, sec: %d, flag: %t", *min, *sec, *flag),
				fmt.Sprintf("minStart: %d, secStart: %d", *min, *sec),
				getAnswerUser,
			),
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

	bimCoin := AddBimCoin(getAnswerUser)
	if bimCoin > 0 {
		FastTest = append(FastTest, id)
	}

	err = db.UpdateData(db_core.TableUsers, "bim_coin1 = ?, bim_total = bim_total + ? ", "id = ?", bimCoin, bimCoin, id)
	if err != nil {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusInternalServerError, ipLogg,
			fmt.Sprintf("Session: %s, TimeTestNow: %s, TimeTest: %s, ID: %d, BIMCoin: %d, Body: %v",
				session,
				fmt.Sprintf("min: %d, sec: %d, flag: %t", *min, *sec, *flag),
				fmt.Sprintf("minStart: %d, secStart: %d", *min, *sec),
				id, bimCoin,
				getAnswerUser,
			),
			err.Error(),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = db.UpdateData(db_core.TableUsers, "test_first = ?", "id = ?", true, id)
	if err != nil {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusInternalServerError, ipLogg,
			fmt.Sprintf("Session: %s, TimeTestNow: %s, TimeTest: %s, ID: %d, BIMCoin: %d, Body: %v",
				session,
				fmt.Sprintf("min: %d, sec: %d, flag: %t", *min, *sec, *flag),
				fmt.Sprintf("minStart: %d, secStart: %d", *min, *sec),
				id, bimCoin,
				getAnswerUser,
			),
			err.Error(),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = db.UpdateData(db_core.TableUsers, "time_test = ?", "id = ?", timeText, id)
	if err != nil {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusInternalServerError, ipLogg,
			fmt.Sprintf("Session: %s, TimeTestNow: %s, TimeTest: %s, ID: %d, TimeText: %s, Body: %v",
				session,
				fmt.Sprintf("min: %d, sec: %d, flag: %t", *min, *sec, *flag),
				fmt.Sprintf("minStart: %d, secStart: %d", *min, *sec),
				id, timeText,
				getAnswerUser,
			),
			err.Error(),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	constants.NewUsers = true

	utils.LogginAPI(pathLogg, methodLogg, fiber.StatusAccepted, ipLogg,
		fmt.Sprintf("Session: %s, TimeTestNow: %s, TimeTest: %s, ID: %d, TimeText: %s, Body: %v",
			session,
			fmt.Sprintf("min: %d, sec: %d, flag: %t", *min, *sec, *flag),
			fmt.Sprintf("minStart: %d, secStart: %d", *min, *sec),
			id, timeText,
			getAnswerUser,
		),
		constants.SuccGetAnswer,
	)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": constants.SuccGetAnswer,
	})
}
