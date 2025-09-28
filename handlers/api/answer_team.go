package api

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"lesson_server/constants"
	db_core "lesson_server/database"
	"lesson_server/file"
	"lesson_server/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type JSONData struct {
	ConflictElements string `json:"conflictElements"`
	CollisionsCount  string `json:"collisionsCount"`
	Priority         string `json:"priority"`
	RiskAssessment   string `json:"riskAssessment"`
	Consequences     string `json:"consequences"`
	Problems         string `json:"problems"`
	Specialists      string `json:"specialists"`
}

var (
	filePathAssets = "pr2/assets"
)

func SaveFile(header *multipart.FileHeader, folder string) (string, error) {
	if err := os.MkdirAll(folder, 0755); err != nil {
		return "", fmt.Errorf("не удалось создать папку: %w", err)
	}

	fileName := fmt.Sprintf("%d_%s_%s", time.Now().Unix(), uuid.NewString(), header.Filename)
	filePath := filepath.Join(folder, fileName)

	src, err := header.Open()
	if err != nil {
		return "", fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("не удалось создать файл: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("не удалось сохранить файл: %w", err)
	}

	return filePath, nil
}

func processFormFile(c *fiber.Ctx, fieldName, folder string) (string, error) {
	file, err := c.FormFile(fieldName)
	var filePath string
	if err == nil {
		filePath, err = SaveFile(file, folder)
		if err != nil {
			return "", fmt.Errorf("ошибка сохранения %s: %w", fieldName, err)
		}

	}

	return filePath, nil
}

func CheckAnswerTeam(c *fiber.Ctx, db *db_core.DatabaseStruct) error {
	pathLogg, methodLogg, ipLogg := c.Path(), c.Method(), c.IP()

	session := c.Cookies(constants.SessionKey)
	id, _ := utils.VerifyJWT(session)
	user, err := db.GetOneUser("id = ?", id)
	if err != nil {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusBadRequest, ipLogg,
			fmt.Sprintf("ID: %d, user: %v", id, user),
			constants.ErrUserNotFound,
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrUserNotFound,
		})
	}

	setting, err := db.GetSetting()
	if err != nil {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusBadRequest, ipLogg,
			fmt.Sprintf("setting: %v", setting),
			err.Error(),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var collumnSettingTest string

	switch user.Team {
	case 1:
		if setting.TestTeamFirst == "1" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": constants.ErrAlreadyReplied,
			})
		}
		collumnSettingTest = "test_team_first = ?"
	case 2:
		if setting.TestTeamSecond == "1" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": constants.ErrAlreadyReplied,
			})
		}
		collumnSettingTest = "test_team_second = ?"
	}

	var globalData JSONData
	err = json.Unmarshal([]byte(c.FormValue("globalData")), &globalData)
	if err != nil {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusBadRequest, ipLogg,
			fmt.Sprintf("globalData: %v", globalData),
			constants.ErrInternalServer,
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInternalServer,
		})
	}

	fileFields := []struct {
		fieldName string
		key       string
	}{
		{"visualization1", "visualization1"},
		{"visualization2", "visualization2"},
		{"visualization3", "visualization3"},
		{"matrix_file", "matrix_file"},
	}

	for _, fileField := range fileFields {
		filePath, err := processFormFile(c, fileField.fieldName, filePathAssets)
		if err != nil {
			utils.LogginAPI(pathLogg, methodLogg, fiber.StatusConflict, ipLogg,
				fmt.Sprintf("fileField.fieldName: %v", fileField.fieldName),
				err.Error(),
			)
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		if len(filePath) == 0 {
			filePath = "Нет Визуализации коллизии 3"
		}

		err = db.UpsertAnswerTaskTwo(user.Team, []string{fileField.fieldName}, filePath)
		if err != nil {
			utils.LogginAPI(pathLogg, methodLogg, fiber.StatusConflict, ipLogg,
				fmt.Sprintf("fileField.fieldName: %v", fileField.fieldName),
				err.Error(),
			)
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	}

	keysForUpdateData := []string{
		"conflict_elements",
		"collisions_count",
		"priority",
		"risk_assessment",
		"consequences",
		"problems",
		"specialists",
	}
	dataForUpdateData := []any{
		globalData.ConflictElements,
		globalData.CollisionsCount,
		globalData.Priority,
		globalData.RiskAssessment,
		globalData.Consequences,
		globalData.Problems,
		globalData.Specialists,
	}

	err = db.UpsertAnswerTaskTwo(user.Team, keysForUpdateData, dataForUpdateData...)
	if err != nil {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusConflict, ipLogg,
			fmt.Sprintf("team: %d, keysForUpdateData: %v, dataForUpdateData: %v", user.Team, keysForUpdateData, dataForUpdateData),
			err.Error(),
		)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = db.UpdateData(db_core.TableSetting, collumnSettingTest, "", "1")
	if err != nil {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusConflict, ipLogg,
			fmt.Sprintf("collumnSettingTest: %v", collumnSettingTest),
			err.Error(),
		)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = file.CreateReportTestTwo(user.Team, db)
	if err != nil {
		utils.LogginAPI(pathLogg, methodLogg, fiber.StatusConflict, ipLogg,
			fmt.Sprintf("team: %d", user.Team),
			err.Error(),
		)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	utils.LogginAPI(pathLogg, methodLogg, fiber.StatusAccepted, ipLogg,
		fmt.Sprintf("dataForUpdateData: %v", dataForUpdateData),
		constants.SuccGetAnswer,
	)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": constants.SuccGetAnswer,
	})
}
