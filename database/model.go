package db

type UserStruct struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Password  string  `json:"password"`
	Status    string  `json:"status"`
	BimCoin1  int     `json:"bim_coin1"`
	BimCoin2  float64 `json:"bim_coin2"`
	BimTotal  float64 `json:"bim_total"`
	Team      int     `json:"team"`
	TestFirst bool    `json:"test_first"`
	TimeTest  string  `json:"time_test"`
}

type SuperUserSturct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SendUserStruct struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	BimCoin1 int     `json:"bim_coin1"`
	BimCoin2 float64 `json:"bim_coin2"`
	BimTotal float64 `json:"bim_total"`
	TimeTest string  `json:"time_test"`
	Team     int     `json:"team"`
}

type SettingStruct struct {
	NowStageLesson string `json:"now_stage_lesson"`
	IdPresentation string `json:"id_presentation"`
	TestTeamFirst  string `json:"test_team_first"`
	TestTeamSecond string `json:"test_team_second"`
}

type AnswerTaskTwo struct {
	ConflictElements string `json:"conflict_elements"`
	Visualization1   string `json:"visualization1"`
	Visualization2   string `json:"visualization2"`
	Visualization3   string `json:"visualization3"`
	MatrixFile       string `json:"matrix_file"`
	CollisionsCount  string `json:"collisions_count"`
	Priority         string `json:"priority"`
	RiskAssessment   string `json:"risk_assessment"`
	Consequences     string `json:"consequences"`
	Problems         string `json:"problems"`
	Specialists      string `json:"specialists"`
	Stage            int    `json:"stage"`
	Team             int    `json:"team"`
}
