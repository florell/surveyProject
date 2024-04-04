package main

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
	SurveyID  int
	Title     string
	Questions []Question
}

type SurveyResults struct {
	SurveyID  int
	PatientID int
	Picked    map[int]int
}

type FamilyEnvironmentalScaleResponse struct {
	Cohesion                        int // Сплоченность
	Expressiveness                  int // Экспрессивность
	Conflict                        int // Конфликт
	Independence                    int // Независимость
	AchievementOrientation          int // Ориентация на достижения
	IntellectualCulturalOrientation int // Интеллектуально-культурная ориентация
	FocusOnActiveRecreation         int // Ориентация на активный отдых
	MoralAspects                    int // Морально-нравственные аспекты
	Organization                    int // Организация
	Control                         int // Контроль
}

type WaysOfCopingQuestionnaireResponse struct {
	Confrontation        int // Конфронтация
	Distancing           int // Дистанцирование
	SelfControl          int // Самоконтроль
	SearchSocialSupport  int // Поиск соц. поддержки
	TakingResponsibility int // Принятие ответственности
	EscapeAvoidance      int // Бегство-избегание
	PlanningTheSolution  int // Планирование решения
	PositiveRevaluation  int // Положительная переоценка
}
