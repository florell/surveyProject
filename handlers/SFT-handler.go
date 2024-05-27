package handlers

import (
	"encoding/json"
	"log"
	types "psychward/src"
)

type SFTResult struct {
	Result struct {
		Value    int `json:"value"`
		MaxValue int `json:"max_value"`
	} `json:"Количество слов"`
	Description string `json:"description"`
}

func SFTHandler(s *types.SurveyResults) []byte {
	result := SFTResult{
		Result: struct {
			Value    int `json:"value"`
			MaxValue int `json:"max_value"`
		}{0, -1},
		Description: "",
	}
	
	for _, value := range s.Picked {
		result.Result.Value = value
	}
	
	switch {
	case result.Result.Value >= 18:
		result.Description = "Нарушений не обнаружено"
	case result.Result.Value >= 12 && result.Result.Value <= 17:
		result.Description = "Умеренное когнитивное снижение"
	case result.Result.Value <= 11:
		result.Description = "Выраженное когнитивное снижение"
	}
	
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	return resultJSON
}
