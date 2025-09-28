package main

import (
	"fmt"
	"lesson_server/constants"
	db_core "lesson_server/database"
	"os"

	"log"
)

func InitTeacher(db *db_core.DatabaseStruct) error {
	_, err := db.GetOneUser("status = ?", constants.TeacherStatus)
	if err == nil {
		log.Print(constants.ErrTeacherAlreadyExist)
		return nil
	}

	newUser := db_core.UserStruct{
		Name:      "Учитель",
		Password:  "BIM_LOCAL123",
		Status:    constants.TeacherStatus,
		BimCoin1:  0,
		BimCoin2:  0,
		Team:      -1,
		TestFirst: false,
		BimTotal:  0,
	}
	_, err = db.InsertUser(newUser)
	if err != nil {
		return err
	}

	log.Print(constants.SuccTeacherExist)
	return nil
}

func InitSetting(db *db_core.DatabaseStruct) error {
	_, err := db.GetSetting()

	var filePathAssets = "pr2"
	if err := os.MkdirAll(filePathAssets, 0755); err != nil {
		fmt.Printf("Ошибка создания папки %s: %v\n", filePathAssets, err)
	} else {
		fmt.Printf("Папка создана: %s\n", filePathAssets)
	}

	if err == nil {
		log.Print(constants.ErrSettingAlreadyExist)
		return nil
	}

	newSetting := db_core.SettingStruct{
		NowStageLesson: "",
		IdPresentation: "",
		TestTeamFirst:  "",
		TestTeamSecond: "",
	}
	err = db.InsertSetting(newSetting)
	if err != nil {
		return err
	}

	log.Print(constants.SuccSettingExist)
	return nil
}
