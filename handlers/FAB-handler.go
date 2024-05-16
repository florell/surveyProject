package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	types "psychward/src"
)

func FABHandler(s *types.SurveyResults) []byte {
	result := map[string]map[string]string{
		"Результат": {
			"value":       "0",
			"max_value":   "18",
			"description": "",
		},
	}
	
	count := 0
	for _, answer := range s.Picked {
		count += answer
	}
	result["Результат"]["value"] = fmt.Sprintf("%d", count)
	
	switch {
	case count >= 16 && count <= 18:
		result["Результат"]["description"] = "Нормальная лобная функция"
	case count >= 12 && count <= 15:
		result["Результат"]["description"] = "Умеренная лобная дисфункция (легкие когнитивные расстройства)"
	case count < 12:
		result["Результат"]["description"] = "Выраженная лобная дисфункция (деменция лобного типа)"
	}
	
	fmt.Println("####")
	fmt.Println(result)
	
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	return resultJSON
}
