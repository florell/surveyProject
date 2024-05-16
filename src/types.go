package types

import "html/template"

type Answer struct {
	AnswerID int
	Text     string
	Value    int
}

type Question struct {
	QuestionID int
	Title      template.HTML
	Answers    []Answer
}

type Survey struct {
	SurveyID    int
	Title       template.HTML
	Description string
	Questions   []Question
}

type SurveyResults struct {
	SurveyID  int
	PatientID int
	Age       int
	Sex       string // 'Man' or 'Woman'
	Picked    map[int]int
}
