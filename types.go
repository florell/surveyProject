package main

type Answer struct {
	ID    int
	Text  string
	Value string
}

type Question struct {
	ID      int
	Title   string
	Answers []Answer
}

type Survey struct {
	ID        int
	Title     string
	Questions []Question
}

type SurveyWithAnswers struct {
	ID      int
	Title   string
	Answers []int
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
