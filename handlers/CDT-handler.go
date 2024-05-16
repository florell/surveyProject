package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	types "psychward/src"
)

func CDTHandler(s *types.SurveyResults) []byte {
	result := map[string]map[string]string{
		"Тест рисования часов": {
			"value":       fmt.Sprintf("%d", s.Picked[398]),
			"max_value":   "10",
			"description": "",
		},
	}
	
	if s.Picked[398] != 10 {
		result["Тест рисования часов"]["description"] = "Наличие выраженных нарушений памяти"
	} else {
		result["Тест рисования часов"]["description"] = "Нарушений оптико-пространственного гнозиса не выявлено, конструктивный праксис сохранен"
	}
	
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	return resultJSON
}
