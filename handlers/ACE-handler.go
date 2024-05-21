package handlers

import (
	"encoding/json"
	"log"
	types "psychward/src"
)

type ACEField struct {
	Value    int `json:"value"`
	MaxValue int `json:"max_value"`
}

type ACEResult struct {
	Overall     ACEField `json:"Общий счет ACE-III"`
	Attention   ACEField `json:"Внимание"`
	Memory      ACEField `json:"Память"`
	Fluency     ACEField `json:"Беглость"`
	Language    ACEField `json:"Язык"`
	VSO         ACEField `json:"Зрительно-пространственная ориентация"`
	Description string   `json:"description"`
}

func ACEHandler(s *types.SurveyResults) []byte {
	result := ACEResult{
		Overall:     ACEField{0, 100},
		Attention:   ACEField{0, 18},
		Memory:      ACEField{0, 26},
		Fluency:     ACEField{0, 14},
		Language:    ACEField{0, 26},
		VSO:         ACEField{0, 16},
		Description: "",
	}

	minId := 100000
	for key, _ := range s.Picked {
		if key < minId {
			minId = key
		}
	}

	for key, value := range s.Picked {
		result.Overall.Value += value

		id := key%minId + 1
		switch id {
		case 1, 2, 3, 4:
			result.Attention.Value += value
		case 5, 8, 9, 23, 24:
			result.Memory.Value += value
		case 6, 7:
			result.Fluency.Value += value
		case 10, 11, 12, 13, 14, 15, 16, 17:
			result.Language.Value += value
		case 18, 19, 20, 21, 22:
			result.VSO.Value += value
		}
	}

	if result.Overall.Value >= 88 {
		result.Description = "Норма"
	} else {
		result.Description = "Не норма"
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	return resultJSON
}
