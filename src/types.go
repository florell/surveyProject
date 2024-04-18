package types

type Answer struct {
	AnswerID int
	Text     string
	Value    int
}

type Question struct {
	QuestionID int
	Title      string
	Answers    []Answer
}

type Survey struct {
	SurveyID    int
	Title       string
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
