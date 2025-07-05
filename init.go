package main

import (
	"lesson_server/constants"
	db_core "lesson_server/database"

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
		BimCoin:   0,
		Team:      -1,
		TestFirst: false,
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
	if err == nil {
		log.Print(constants.ErrSettingAlreadyExist)
		return nil
	}

	newSetting := db_core.SettingStruct{
		NowStageLesson: "",
		IdPresentation: "",
		TestTeamFirst:  false,
		TestTeamSecond: false,
	}
	err = db.InsertSetting(newSetting)
	if err != nil {
		return err
	}

	log.Print(constants.SuccSettingExist)
	return nil
}
