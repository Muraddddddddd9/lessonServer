package file

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	db_core "lesson_server/database"

	"github.com/fumiama/go-docx"
)

func AddParagraphStyle(paragraph *docx.Paragraph, just string, line int) {
	paragraph.Properties = &docx.ParagraphProperties{
		Justification: &docx.Justification{
			Val: just,
		},
		Spacing: &docx.Spacing{
			Line: line,
		},
	}
}

func AddStyle(val *docx.Run, style string) {
	switch style {
	case "main_title":
		val.Size("28").Color("000000").Font("Times New Roman", "Times New Roman", "Times New Roman", "cs").Bold()
	case "block_title":
		val.Size("24").Color("1F4F79").Font("Times New Roman", "Times New Roman", "Times New Roman", "cs").Bold()
	case "block_title_text":
		val.Size("22").Color("000000").Font("Times New Roman", "Times New Roman", "Times New Roman", "cs").Bold()
	case "block_text":
		val.Size("22").Color("000000").Font("Times New Roman", "Times New Roman", "Times New Roman", "cs")
	}
}

func CreateTitle(w *docx.Docx) {
	titleParaOne := w.AddParagraph()
	titleParaTwo := w.AddParagraph()

	AddParagraphStyle(titleParaOne, "center", 350)
	AddParagraphStyle(titleParaTwo, "center", 350)

	titlePartOne := titleParaOne.AddText("Отчёт по выявленным коллизиям")
	titlePartTwo := titleParaTwo.AddText("в цифровой информационной модели здания")
	AddStyle(titlePartOne, "main_title")
	AddStyle(titlePartTwo, "main_title")
}

func AddAnswerText(arr []DataText) {
	for _, v := range arr {
		AddParagraphStyle(v.p, "start", 400)
		if len(v.question) != 0 {
			text := v.p.AddText(v.question).AddTab()
			AddStyle(text, "block_title_text")
		}
		answerText := v.p.AddText(v.answer)
		AddStyle(answerText, "block_text")
	}
}

type DataImg struct {
	p   *docx.Paragraph
	img string
}

type DataText struct {
	p        *docx.Paragraph
	question string
	answer   string
}

func CreateTitleParagraph(w *docx.Docx, text, answer string) {
	blockTitle := w.AddParagraph()
	if len(text) != 0 {
		blockTitleText := blockTitle.AddText(text).AddTab()
		AddStyle(blockTitleText, "block_title")
	}

	if len(text) != 0 {
		blockTitleAnserText := blockTitle.AddText(answer)
		AddStyle(blockTitleAnserText, "block_text")
	}
	AddParagraphStyle(blockTitle, "start", 400)
}

func CreateFirstBlock(w *docx.Docx, answerText, answerImg1, answerImg2, answerImg3 string) {
	CreateTitleParagraph(w, "1. Элементы конфликта:", answerText)

	blockFirstImg := w.AddParagraph()
	blockSecondImg := w.AddParagraph()
	blockThirdImg := w.AddParagraph()

	arrImg := []DataImg{
		{p: blockFirstImg, img: answerImg1},
		{p: blockSecondImg, img: answerImg2},
		{p: blockThirdImg, img: answerImg3},
	}

	for i, v := range arrImg {
		blockText := v.p.AddText(fmt.Sprintf("Визуализация коллизии %d", i+1))
		v.p.AddText("\n")

		if len(v.img) != 0 {
			_, err := v.p.AddInlineDrawingFrom(v.img)
			if err != nil {
				fmt.Printf("Warning: could not add image: %v\n", err)
				v.p.AddText("[Картинка]").Color("FF0000").Italic()
			}
		}

		AddParagraphStyle(v.p, "start", 400)
		AddStyle(blockText, "block_title_text")
	}
}

func CreateSecondBlock(w *docx.Docx, answerImg, answerText1, answerText2 string) {
	CreateTitleParagraph(w, "2. Матрица коллизий:", "")

	imgPara := w.AddParagraph()
	_, err := imgPara.AddInlineDrawingFrom(answerImg)
	if err != nil {
		fmt.Printf("Warning: could not add image: %v\n", err)
		imgPara.AddText("[Картинка]").Color("FF0000").Italic()
	}

	blockFirstText := w.AddParagraph()
	blockSecondText := w.AddParagraph()

	arrAnswer := []DataText{
		{p: blockFirstText, question: "Количество коллизий:", answer: answerText1},
		{p: blockSecondText, question: "Приоритет устранения:", answer: answerText2},
	}

	AddAnswerText(arrAnswer)
}

func CreateThirdBlock(w *docx.Docx, answerText1, answerText2 string) {
	CreateTitleParagraph(w, "3. Оценка производственного риска:", "")

	blockFirstText := w.AddParagraph()
	blockSecondText := w.AddParagraph()

	arrAnswer := []DataText{
		{p: blockFirstText, question: "", answer: answerText1},
		{p: blockSecondText, question: "Последствия:", answer: answerText2},
	}

	AddAnswerText(arrAnswer)
}

func CreateFouthBlock(w *docx.Docx, answerText1, answerText2 string) {
	CreateTitleParagraph(w, "4. Рекомендации для заказчика:", "")

	blockFirstText := w.AddParagraph()
	blockSecondText := w.AddParagraph()

	arrAnswer := []DataText{
		{p: blockFirstText, question: "Выявленные проблемы:", answer: answerText1},
		{p: blockSecondText, question: "В результате проведения экспертизы, рекомендовано направить отчет для устранения выявленных коллизий следующим специалистам:", answer: answerText2},
	}

	AddAnswerText(arrAnswer)
}

func CreateSignBlock(w *docx.Docx, time, team, users string) {
	blockEntry := w.AddParagraph()
	blockEntry1 := w.AddParagraph()
	blockTime := w.AddParagraph()
	blockTeam := w.AddParagraph()
	blockUsers := w.AddParagraph()

	arrAnswer := []DataText{
		{p: blockEntry, question: "", answer: ""},
		{p: blockEntry1, question: "", answer: ""},
		{p: blockTime, question: "Дата составления отчета:", answer: time},
		{p: blockTeam, question: "BIM-отдел: №", answer: team},
		{p: blockUsers, question: "Стажеры отдела:", answer: users},
	}

	AddAnswerText(arrAnswer)
}

func CreateReportTestTwo(team int, db *db_core.DatabaseStruct) error {
	dataAnswer, err := db.GetAnswerTaskTwo(team)
	if err != nil {
		return err
	}

	w := docx.New().WithDefaultTheme().WithA4Page()
	CreateTitle(w)

	CreateFirstBlock(w, dataAnswer.ConflictElements, dataAnswer.Visualization1, dataAnswer.Visualization2, dataAnswer.Visualization3)
	CreateSecondBlock(w, dataAnswer.MatrixFile, dataAnswer.CollisionsCount, dataAnswer.Priority)
	CreateThirdBlock(w, dataAnswer.RiskAssessment, dataAnswer.Consequences)
	CreateFouthBlock(w, dataAnswer.Problems, dataAnswer.Specialists)

	usersTeam, err := db.GetTeamUser(team)
	if err != nil {
		return err
	}

	var userName []string
	for _, v := range usersTeam {
		userName = append(userName, v.Name)
	}

	CreateSignBlock(w, time.Now().Format("01-02-2006"), strconv.Itoa(team), strings.Join(userName, ", "))

	f, err := os.Create(fmt.Sprintf("pr2/file_%d.docx", team))
	if err != nil {
		return err
	}
	_, err = w.WriteTo(f)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}
