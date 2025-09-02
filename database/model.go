package db

type UserStruct struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Status    string `json:"status"`
	BimCoin   int    `json:"bim_coin"`
	Team      int    `json:"team"`
	TestFirst bool   `json:"test_first"`
	TimeTest  string `json:"time_test"`
}

type SendUserStruct struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	BimCoin  int    `json:"bim_coin"`
	TimeTest string `json:"time_test"`
	Team     int    `json:"team"`
}

type SettingStruct struct {
	NowStageLesson string `json:"now_stage_lesson"`
	IdPresentation string `json:"id_presentation"`
	TestTeamFirst  bool   `json:"test_team_first"`
	TestTeamSecond bool   `json:"test_team_second"`
}
