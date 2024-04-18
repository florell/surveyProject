package handlers

import (
	"encoding/json"
	"log"
	types "psychward/src"
)

type ITTField struct {
	MaxValue int `json:"max_value"`
	Value    int `json:"value"`
}

type ITTResult struct {
	SituationalAnxietySTS  ITTField `json:"Ситуационная тревожность (СТ-С)"`
	EmotionalDiscomfortEDS ITTField `json:"Эмоциональный дискомфорт (ЭД-С)"`
	AsthenicComponentASTS  ITTField `json:"Астенический компонент (АСТ-С)"`
	PhobicComponentFOBS    ITTField `json:"Фобический компонент (ФОБ-С)"`
	AnxietyAssessmentOPS   ITTField `json:"Тревожная оценка перспективы (ОП-С)"`
	SocialProtectionSZS    ITTField `json:"Социальная защита (СЗ-С)"`
	PersonalAnxietySTL     ITTField `json:"Личностная тревожность (СТ-Л)"`
	EmotionalDiscomfortEDL ITTField `json:"Эмоциональный дискомфорт (ЭД-Л)"`
	AsthenicComponentASTL  ITTField `json:"Астенический компонент (АСТ-Л)"`
	PhobicComponentFOBL    ITTField `json:"Фобический компонент (ФОБ-Л)"`
	AnxietyAssessmentOPL   ITTField `json:"Тревожная оценка перспективы (ОП-Л)"`
	SocialProtectionSZL    ITTField `json:"Социальная защита (СЗ-Л)"`
}

func getGeneralScore(value, age int, sex string) int {
	if sex == "male" || age > 16 {
		switch {
		case value <= 6:
			return 1
		case value >= 7 && value <= 8:
			return 2
		case value == 9:
			return 3
		case value >= 10 && value <= 11:
			return 4
		case value >= 12 && value <= 14:
			return 5
		case value >= 15 && value <= 18:
			return 6
		case value >= 19 && value <= 22:
			return 7
		case value >= 23 && value <= 26:
			return 8
		default: // value >= 27
			return 9
		}
	} else {
		switch {
		case value <= 6:
			return 1
		case value >= 7 && value <= 8:
			return 2
		case value >= 9 && value <= 10:
			return 3
		case value >= 11 && value <= 12:
			return 4
		case value >= 13 && value <= 16:
			return 5
		case value >= 17 && value <= 21:
			return 6
		case value >= 22 && value <= 25:
			return 7
		case value >= 26 && value <= 30:
			return 8
		default: // value >= 31
			return 9
		}
	}
}

func getEDValue(id, value int) int {
	switch id {
	case 1:
		switch value {
		case 1:
			return 25
		case 2:
			return 49
		case 3:
			return 74
		default:
			return 0
		}
	case 2, 6:
		switch value {
		case 1:
			return 24
		case 2:
			return 49
		case 3:
			return 73
		default:
			return 0
		}
	case 4:
		switch value {
		case 1:
			return 27
		case 2:
			return 53
		case 3:
			return 80
		default:
			return 0
		}
	default:
		return 0
	}
}

func getASTValue(id, value int) int {
	switch id {
	case 8:
		switch value {
		case 1:
			return 30
		case 2:
			return 61
		case 3:
			return 91
		default:
			return 0
		}
	case 13:
		switch value {
		case 1:
			return 41
		case 2:
			return 81
		case 3:
			return 122
		default:
			return 0
		}
	case 14:
		switch value {
		case 1:
			return 29
		case 2:
			return 58
		case 3:
			return 87
		default:
			return 0
		}
	default:
		return 0
	}
}

func getFOBValue(id, value int) int {
	switch id {
	case 7:
		switch value {
		case 1:
			return 37
		case 2:
			return 74
		case 3:
			return 111
		default:
			return 0
		}
	case 9:
		switch value {
		case 1:
			return 28
		case 2:
			return 56
		case 3:
			return 85
		default:
			return 0
		}
	case 12:
		switch value {
		case 1:
			return 29
		case 2:
			return 58
		case 3:
			return 87
		default:
			return 0
		}
	default:
		return 0
	}
}

func getOPValue(id, value int) int {
	switch id {
	case 3:
		switch value {
		case 1:
			return 37
		case 2:
			return 74
		case 3:
			return 110
		default:
			return 0
		}
	case 5:
		switch value {
		case 1:
			return 32
		case 2:
			return 65
		case 3:
			return 98
		default:
			return 0
		}
	case 15:
		switch value {
		case 1:
			return 31
		case 2:
			return 61
		case 3:
			return 92
		default:
			return 0
		}
	default:
		return 0
	}
}

func getSZValue(id, value int) int {
	switch id {
	case 10:
		switch value {
		case 1:
			return 57
		case 2:
			return 114
		case 3:
			return 171
		default:
			return 0
		}
	case 11:
		switch value {
		case 1:
			return 43
		case 2:
			return 86
		case 3:
			return 129
		default:
			return 0
		}
	default:
		return 0
	}
}

func getEDStan(value, age int, sex string) int {
	if sex == "male" || age > 16 {
		switch {
		case value <= 34:
			return 1
		case value >= 35 && value <= 48:
			return 2
		case value >= 49 && value <= 62:
			return 3
		case value >= 63 && value <= 76:
			return 4
		case value >= 77 && value <= 100:
			return 5
		case value >= 101 && value <= 137:
			return 6
		case value >= 138 && value <= 173:
			return 7
		case value >= 174 && value <= 209:
			return 8
		default: // >= 210
			return 9
		}
	} else {
		switch {
		case value <= 42:
			return 1
		case value >= 43 && value <= 58:
			return 2
		case value >= 59 && value <= 75:
			return 3
		case value >= 76 && value <= 92:
			return 4
		case value >= 93 && value <= 117:
			return 5
		case value >= 118 && value <= 150:
			return 6
		case value >= 151 && value <= 183:
			return 7
		case value >= 184 && value <= 217:
			return 8
		default: // >= 218
			return 9
		}
	}
}

func getASTStan(value, age int, sex string) int {
	if sex == "male" || age > 16 {
		switch {
		case value <= 26:
			return 1
		case value >= 27 && value <= 36:
			return 2
		case value >= 37 && value <= 47:
			return 3
		case value >= 48 && value <= 57:
			return 4
		case value >= 58 && value <= 82:
			return 5
		case value >= 83 && value <= 122:
			return 6
		case value >= 123 && value <= 161:
			return 7
		case value >= 162 && value <= 201:
			return 8
		default: // >= 202
			return 9
		}
	} else {
		switch {
		case value <= 31:
			return 1
		case value >= 32 && value <= 44:
			return 2
		case value >= 45 && value <= 57:
			return 3
		case value >= 58 && value <= 70:
			return 4
		case value >= 71 && value <= 95:
			return 5
		case value >= 96 && value <= 132:
			return 6
		case value >= 133 && value <= 169:
			return 7
		case value >= 170 && value <= 206:
			return 8
		default: // >= 207
			return 9
		}
	}
}

func getFOBStan(value, age int, sex string) int {
	if sex == "male" || age > 16 {
		switch {
		case value <= 13:
			return 1
		case value >= 14 && value <= 19:
			return 2
		case value >= 20 && value <= 24:
			return 3
		case value >= 25 && value <= 29:
			return 4
		case value >= 30 && value <= 54:
			return 5
		case value >= 55 && value <= 99:
			return 6
		case value >= 100 && value <= 144:
			return 7
		case value >= 145 && value <= 188:
			return 8
		default: // >= 189
			return 9
		}
	} else {
		switch {
		case value <= 16:
			return 1
		case value >= 17 && value <= 23:
			return 2
		case value >= 24 && value <= 29:
			return 3
		case value >= 30 && value <= 36:
			return 4
		case value >= 37 && value <= 61:
			return 5
		case value >= 62 && value <= 104:
			return 6
		case value >= 105 && value <= 148:
			return 7
		case value >= 149 && value <= 191:
			return 8
		default: // >= 192
			return 9
		}
	}
}

func getOPStan(value, age int, sex string) int {
	if sex == "male" || age > 16 {
		switch {
		case value <= 44:
			return 1
		case value >= 45 && value <= 62:
			return 2
		case value >= 63 && value <= 80:
			return 3
		case value >= 81 && value <= 97:
			return 4
		case value >= 98 && value <= 122:
			return 5
		case value >= 123 && value <= 155:
			return 6
		case value >= 156 && value <= 187:
			return 7
		case value >= 188 && value <= 219:
			return 8
		default: // >= 220
			return 9
		}
	} else {
		switch {
		case value <= 53:
			return 1
		case value >= 54 && value <= 75:
			return 2
		case value >= 76 && value <= 97:
			return 3
		case value >= 98 && value <= 118:
			return 4
		case value >= 119 && value <= 143:
			return 5
		case value >= 144 && value <= 172:
			return 6
		case value >= 173 && value <= 200:
			return 7
		case value >= 201 && value <= 228:
			return 8
		default: // >= 229
			return 9
		}
	}
}

func getSZStan(value, age int, sex string) int {
	if sex == "male" || age > 16 {
		switch {
		case value <= 50:
			return 1
		case value >= 51 && value <= 70:
			return 2
		case value >= 71 && value <= 90:
			return 3
		case value >= 91 && value <= 110:
			return 4
		case value >= 111 && value <= 135:
			return 5
		case value >= 136 && value <= 165:
			return 6
		case value >= 166 && value <= 195:
			return 7
		case value >= 196 && value <= 225:
			return 8
		default: // >= 226
			return 9
		}
	} else {
		switch {
		case value <= 60:
			return 1
		case value >= 61 && value <= 85:
			return 2
		case value >= 86 && value <= 109:
			return 3
		case value >= 110 && value <= 134:
			return 4
		case value >= 135 && value <= 159:
			return 5
		case value >= 160 && value <= 184:
			return 6
		case value >= 185 && value <= 210:
			return 7
		case value >= 211 && value <= 235:
			return 8
		default: // >= 236
			return 9
		}
	}
}

func ITTHandler(s *types.SurveyResults) []byte {
	result := ITTResult{}
	sumS, sumL := 0, 0
	
	for questionID, answer := range s.Picked {
		id := questionID - 90 - 50 - 26 - 21
		if id <= 15 {
			sumS += answer
			result.EmotionalDiscomfortEDS.Value += getEDValue(id, answer)
			result.AsthenicComponentASTS.Value += getASTValue(id, answer)
			result.PhobicComponentFOBS.Value += getFOBValue(id, answer)
			result.AnxietyAssessmentOPS.Value += getOPValue(id, answer)
			result.SocialProtectionSZS.Value += getSZValue(id, answer)
		} else {
			sumL += answer
			result.EmotionalDiscomfortEDL.Value += getEDValue(id, answer)
			result.AsthenicComponentASTL.Value += getASTValue(id, answer)
			result.PhobicComponentFOBL.Value += getFOBValue(id, answer)
			result.AnxietyAssessmentOPL.Value += getOPValue(id, answer)
			result.SocialProtectionSZL.Value += getSZValue(id, answer)
		}
	}
	
	result.EmotionalDiscomfortEDS.Value = getEDStan(result.EmotionalDiscomfortEDS.Value, s.Age, s.Sex)
	result.AsthenicComponentASTS.Value = getASTStan(result.AsthenicComponentASTS.Value, s.Age, s.Sex)
	result.PhobicComponentFOBS.Value = getFOBStan(result.PhobicComponentFOBS.Value, s.Age, s.Sex)
	result.AnxietyAssessmentOPS.Value = getOPStan(result.AnxietyAssessmentOPS.Value, s.Age, s.Sex)
	result.SocialProtectionSZS.Value = getSZStan(result.SocialProtectionSZS.Value, s.Age, s.Sex)
	
	result.EmotionalDiscomfortEDL.Value = getEDStan(result.EmotionalDiscomfortEDL.Value, s.Age, s.Sex)
	result.AsthenicComponentASTL.Value = getASTStan(result.AsthenicComponentASTL.Value, s.Age, s.Sex)
	result.PhobicComponentFOBL.Value = getFOBStan(result.PhobicComponentFOBL.Value, s.Age, s.Sex)
	result.AnxietyAssessmentOPL.Value = getOPStan(result.AnxietyAssessmentOPL.Value, s.Age, s.Sex)
	result.SocialProtectionSZL.Value = getSZStan(result.SocialProtectionSZL.Value, s.Age, s.Sex)
	
	result.SituationalAnxietySTS.MaxValue = 9
	result.EmotionalDiscomfortEDS.MaxValue = 9
	result.AsthenicComponentASTS.MaxValue = 9
	result.PhobicComponentFOBS.MaxValue = 9
	result.AnxietyAssessmentOPS.MaxValue = 9
	result.SocialProtectionSZS.MaxValue = 9
	
	result.PersonalAnxietySTL.MaxValue = 9
	result.EmotionalDiscomfortEDL.MaxValue = 9
	result.AsthenicComponentASTL.MaxValue = 9
	result.PhobicComponentFOBL.MaxValue = 9
	result.AnxietyAssessmentOPL.MaxValue = 9
	result.SocialProtectionSZL.MaxValue = 9
	
	result.SituationalAnxietySTS.Value = getGeneralScore(sumS, s.Age, s.Sex)
	result.PersonalAnxietySTL.Value = getGeneralScore(sumL, s.Age, s.Sex)
	
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	return resultJSON
}
