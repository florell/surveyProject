package handlers

import (
	"encoding/json"
	"log"
	types "psychward/src"
)

type CDTResult struct {
}

func CDTHandler(s *types.SurveyResults) []byte {
	result := struct {
		Result struct {
			Value    int `json:"value"`
			MaxValue int `json:"max_value"`
		} `json:"Количество баллов"`
		Description string `json:"description"`
	}{struct {
		Value    int `json:"value"`
		MaxValue int `json:"max_value"`
	}{Value: s.Picked[398], MaxValue: -1}, ""}

	if s.Picked[398] >= 10 {
		result.Description = "Нарушений оптико-пространственного гнозиса не выявлено, конструктивный праксис сохранен"
	} else {
		result.Description = "Наличие когнитивных нарушений"
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	return resultJSON
}
