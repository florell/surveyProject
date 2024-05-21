package handlers

import (
	"encoding/json"
	"log"
	types "psychward/src"
)

type PFTResult struct {
	Result struct {
		Value    int
		MaxValue int
	} `json:"Количество слов"`
	Description string `json:"description"`
}

func PFTHandler(s *types.SurveyResults) []byte {
	result := SFTResult{
		Result: struct {
			Value    int
			MaxValue int
		}{0, -1},
		Description: "",
	}
	
	for _, value := range s.Picked {
		result.Result.Value = value
	}
	
	switch {
	case result.Result.Value >= 12:
		result.Description = "норма"
	case result.Result.Value >= 9 && result.Result.Value <= 11:
		result.Description = "умеренное когнитивное снижение"
	case result.Result.Value <= 8:
		result.Description = "выраженное когнитивное снижение"
	}
	
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	return resultJSON
}
