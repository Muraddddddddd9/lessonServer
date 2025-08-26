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
	Socer      int      `json:"socer"`
}

var Questions = map[string]QuestionStruct{
	"Questions_1": {
		ID:       "Questions_1",
		Question: "Пара объектов «Стены – Стены» общее количество геометрических коллизий в соответствии с матрицей пересечений:",
		Answers: []string{
			"0",
			"2",
			"3",
			"5",
			"6",
		},
		AnswerTrue: "2",
		Socer:      5,
	},
	"Questions_2": {
		ID:       "Questions_2",
		Question: "Пара объектов «Стены – Воздуховоды» общее количество геометрических коллизий в соответствии с матрицей пересечений:",
		Answers: []string{
			"0",
			"2",
			"3",
			"5",
			"6",
		},
		AnswerTrue: "5",
		Socer:      10,
	},
}

var QuestionsTeam = map[string]QuestionStruct{
	"Questions_Team_1": {
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
	"Questions_Team_2": {
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
