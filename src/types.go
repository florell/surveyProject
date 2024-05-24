package types

import "html/template"

type Answer struct {
	AnswerID int
	Text     string
	Value    int
}

type Question struct {
	QuestionID int           `json:"ID"`
	Title      template.HTML `json:"Title"`
	Answers    []Answer      `json:"Answers"`
	MaxValue   int           `json:"MaxPoints"`
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
