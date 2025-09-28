package api

import (
	"fmt"
	"lesson_server/constants"
	db_core "lesson_server/database"
	"lesson_server/utils"

	"github.com/gofiber/fiber/v2"
)

type ScoreReq struct {
	Score float64 `json:"score"`
}

func ScoreTeam(c *fiber.Ctx, db *db_core.DatabaseStruct) error {
	pathLogg, methodLogg, ipLogg := c.Path(), c.Method(), c.IP()

	session := c.Cookies(constants.SessionKey)
	id, err := utils.VerifyJWT(session)
	if err != nil {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusInternalServerError, ipLogg,
			fmt.Sprintf("session: %s, id: %v", session, id),
			err.Error(),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user, err := db.GetOneUser("id = ?", id)
	if err != nil {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusInternalServerError, ipLogg,
			fmt.Sprintf("user: %v, id: %v", user, id),
			err.Error(),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var teamSend int

	if user.Team == 1 {
		teamSend = 2
	} else {
		teamSend = 1
	}

	user, err = db.GetOneUser("team = ? LIMIT 1", teamSend)
	if err != nil {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusInternalServerError, ipLogg,
			fmt.Sprintf("user: %v, id: %v", user, teamSend),
			err.Error(),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if user.BimCoin2 > 0 {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusConflict, ipLogg,
			fmt.Sprintf("user: %v, id: %v", user, teamSend),
			constants.ErrScoreTeam,
		)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": constants.ErrScoreTeam,
		})
	}

	var score ScoreReq
	if err := c.BodyParser(&score); err != nil {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusBadRequest, ipLogg,
			fmt.Sprintf("score: %v, teamSend: %d", score, teamSend),
			constants.ErrInputValue,
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInputValue,
		})
	}

	err = db.UpdateData(db_core.TableUsers, "bim_coin2 = ?, bim_total = bim_total + ? ", "team = ?", score.Score, score.Score, teamSend)
	if err != nil {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusInternalServerError, ipLogg,
			fmt.Sprintf("score: %v, teamSend: %d", score, teamSend),
			err.Error(),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	constants.NewUsers = true

	utils.LogginAPI(pathLogg, methodLogg, fiber.StatusAccepted, ipLogg,
		fmt.Sprintf("score: %v, teamSend: %d", score, teamSend),
		constants.SuccScoreTeam,
	)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": constants.SuccScoreTeam,
	})
}
