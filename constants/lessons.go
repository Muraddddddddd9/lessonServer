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
		Question: "«Стены – Стены» количество геометрических коллизий в соответствии с матрицей пересечений:",
		Answers: []string{
			"0",
			"2",
			"3",
			"5",
			"6",
		},
		AnswerTrue: "2",
		Socer:      10,
	},
	"Questions_2": {
		ID:       "Questions_2",
		Question: "«Стены – Воздуховоды» количество геометрических коллизий в соответствии с матрицей пересечений:",
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