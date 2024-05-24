package handlers

import (
	"encoding/json"
	"log"
	types "psychward/src"
)

type MMSEResult struct {
	Result struct {
		Value    int `json:"value"`
		MaxValue int `json:"max_value"`
	} `json:"Количество слов"`
	Description string `json:"description"`
}

func MMSEHandler(s *types.SurveyResults) []byte {
	result := MMSEResult{
		Result: struct {
			Value    int `json:"value"`
			MaxValue int `json:"max_value"`
		}{0, 30},
		Description: "",
	}

	count := 0
	for _, answer := range s.Picked {
		count += answer
	}
	result.Result.Value = count

	switch {
	case count >= 28 && count <= 30:
		result.Description = "Нет нарушений когнитивных функций"
	case count >= 24 && count <= 27:
		result.Description = "Преддементные когнитивные нарушения"
	case count >= 20 && count <= 23:
		result.Description = "Деменция легкой степени выраженности"
	case count >= 11 && count <= 19:
		result.Description = "Деменция умеренной степени выраженности"
	case count <= 10:
		result.Description = "Тяжелая деменция"
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	return resultJSON
}
