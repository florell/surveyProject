package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	types "psychward/src"
)

func MMSEHandler(s *types.SurveyResults) []byte {
	result := map[string]map[string]string{
		"Результат": {
			"value":       "0",
			"max_value":   "30",
			"description": "",
		},
	}
	
	count := 0
	for _, answer := range s.Picked {
		count += answer
	}
	result["Результат"]["value"] = fmt.Sprintf("%d", count)
	
	switch {
	case count >= 28 && count <= 30:
		result["Результат"]["description"] = "Нет нарушений когнитивных функций"
	case count >= 24 && count <= 27:
		result["Результат"]["description"] = "Преддементные когнитивные нарушения"
	case count >= 20 && count <= 23:
		result["Результат"]["description"] = "Деменция легкой степени выраженности"
	case count >= 11 && count <= 19:
		result["Результат"]["description"] = "Деменция умеренной степени выраженности"
	case count <= 10:
		result["Результат"]["description"] = "Тяжелая деменция"
	}
	
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	return resultJSON
}
