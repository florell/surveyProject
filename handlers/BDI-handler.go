package handlers

import (
	"encoding/json"
	"log"
	types "psychward/src"
)

type BDIField struct {
	MaxValue int `json:"max_value"`
	Value    int `json:"value"`
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
	
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	return resultJSON
}
