package handlers

import (
	"encoding/json"
	"log"
	types "psychward/src"
)

type BDIField struct {
	MaxValue    int    `json:"max_value"`
	Value       int    `json:"value"`
	Description string `json:"description"`
}

type BDIResult struct {
	Overall BDIField `json:"Шкала депрессии Бека"`
	KAP     BDIField `json:"Когнитивно-аффективные проявления"`
	SP      BDIField `json:"Соматические проявления"`
}

func BDIHandler(s *types.SurveyResults) []byte {
	result := BDIResult{}
	
	result.Overall.MaxValue = 63
	result.KAP.MaxValue = 39
	result.SP.MaxValue = 24
	
	for questionID, answer := range s.Picked {
		id := questionID - 90 - 50 - 26
		result.Overall.Value += answer
		if id >= 1 && id <= 13 {
			result.KAP.Value += answer
		}
		if id >= 14 && id <= 21 {
			result.SP.Value += answer
		}
	}
	
	switch {
	case result.Overall.Value >= 0 && result.Overall.Value <= 9:
		result.Overall.Description = "Отсутствие депрессивных симптомов"
	case result.Overall.Value >= 10 && result.Overall.Value <= 15:
		result.Overall.Description = "Легкая депрессия (субдепрессия)"
	case result.Overall.Value >= 16 && result.Overall.Value <= 19:
		result.Overall.Description = "Умеренная депрессия"
	case result.Overall.Value >= 20 && result.Overall.Value <= 29:
		result.Overall.Description = "Выраженная депрессия (средней тяжести)"
	case result.Overall.Value >= 30:
		result.Overall.Description = "Тяжелая депрессия"
	}
	
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	return resultJSON
}
