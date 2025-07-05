package constants

var (
	NewLesson = true
	NewUsers  = true
)

type QuestionStruct struct {
	ID         string   `json:"id"`
	Question   string   `json:"question"`
	Answers    []string `json:"answers"`
	AnswerTrue string   `json:"answer_true"`
	Socer      int   `json:"socer"`
}

var Questions = []QuestionStruct{
	{
		ID:       "Questions_1",
		Question: "Как дела?",
		Answers: []string{
			"Норм",
			"Плохо",
			"А ты как думаешь?",
			"Как так?",
		},
		AnswerTrue: "Норм",
		Socer:      1,
	},
	{
		ID:       "Questions_2",
		Question: "А ты как?",
		Answers: []string{
			"Да так",
			"Как как",
			"Вот так",
			"Окак",
		},
		AnswerTrue: "Окак",
		Socer:      4,
	},
}

var QuestionsTeam = []QuestionStruct{
	{
		ID:       "Questions_Team_1",
		Question: "ОГО?",
		Answers: []string{
			"1",
			"2",
			"ОГООГО",
			"3",
		},
		AnswerTrue: "ОГООГО",
		Socer:      1,
	},
	{
		ID:       "Questions_Team_2",
		Question: "АГА",
		Answers: []string{
			"Ок",
			"Нет",
			"Как",
			"OKAK",
		},
		AnswerTrue: "Как",
		Socer:      4,
	},
}
