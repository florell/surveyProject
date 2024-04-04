package handlers

import (
	"encoding/json"
	"log"
	types "psychward/src"
)

var (
	scaleFESKeys = map[string]map[int]int{
		"Cohesion": {
			1:  1,
			11: 1,
			21: 1,
			31: 1,
			51: 1,
			81: 1,
			41: 0,
			61: 0,
			71: 0,
		},
		"Expressiveness": {
			2:  1,
			12: 1,
			32: 1,
			42: 1,
			52: 1,
			62: 1,
			72: 1,
			82: 1,
			22: 0,
		},
		"Conflict": {
			13: 1,
			33: 1,
			63: 1,
			73: 1,
			83: 1,
			3:  0,
			23: 0,
			43: 0,
			53: 0,
		},
		"Independence": {
			4:  1,
			14: 1,
			44: 1,
			54: 1,
			64: 1,
			24: 0,
			34: 0,
			74: 0,
			84: 0,
		},
		"AchievementOrientation": {
			5:  1,
			15: 1,
			35: 1,
			45: 1,
			65: 1,
			75: 1,
			25: 0,
			55: 0,
			85: 0,
		},
		"IntellectualCulturalOrientation": {
			6:  1,
			26: 1,
			56: 1,
			66: 1,
			86: 1,
			16: 0,
			36: 0,
			76: 0,
		},
		"FocusOnActiveRecreation": {
			17: 1,
			37: 1,
			47: 1,
			67: 1,
			77: 1,
			7:  0,
			27: 0,
			57: 0,
			87: 0,
		},
		"MoralAspects": {
			8:  1,
			28: 1,
			58: 1,
			68: 1,
			78: 1,
			88: 1,
			18: 0,
			38: 0,
			48: 0,
		},
		"Organization": {
			9:  1,
			19: 1,
			39: 1,
			59: 1,
			69: 1,
			89: 1,
			29: 0,
			49: 0,
			79: 0,
		},
		"Control": {
			10: 1,
			30: 1,
			40: 1,
			50: 1,
			60: 1,
			80: 1,
			90: 1,
			20: 0,
			70: 0,
		},
	}
)

// answer: 0 - No, 1 - Yes
func FamilyEnvironmentalScaleHandler(s types.SurveyResults) []byte {
	result := make(map[string]int)
	
	for questionID, answer := range s.Picked {
		for field, keysMap := range scaleFESKeys {
			if answer == keysMap[questionID] {
				result[field]++
			}
		}
	}
	
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	return resultJSON
}
