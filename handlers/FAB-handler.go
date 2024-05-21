package handlers

import (
	"encoding/json"
	"log"
	types "psychward/src"
)

func FABHandler(s *types.SurveyResults) []byte {
	result := struct {
		Overall struct {
			Value    int `json:"value"`
			MaxValue int `json:"max_value"`
		} `json:"Количество баллов"`
		Description string `json:"description"`
	}{struct {
		Value    int `json:"value"`
		MaxValue int `json:"max_value"`
	}{0, 18}, ""}

	count := 0
	for _, answer := range s.Picked {
		count += answer
	}
	result.Overall.Value = count

	switch {
	case count >= 16 && count <= 18:
		result.Description = "Нормальная лобная функция"
	case count >= 12 && count <= 15:
		result.Description = "Умеренная лобная дисфункция (легкие когнитивные расстройства)"
	case count < 12:
		result.Description = "Выраженная лобная дисфункция (деменция лобного типа)"
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	return resultJSON
}
