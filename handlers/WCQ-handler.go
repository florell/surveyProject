package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	types "psychward/src"
	"strings"
)

var (
	standardPoints = map[string]map[int][]int{
		"Конфронтация": {
			0:  []int{22, 22, 15, 14, 16, 22, 17, 19},
			1:  []int{26, 25, 19, 18, 20, 26, 21, 23},
			2:  []int{29, 28, 23, 22, 24, 29, 24, 26},
			3:  []int{32, 32, 27, 26, 28, 32, 28, 29},
			4:  []int{35, 35, 31, 30, 32, 35, 32, 33},
			5:  []int{38, 38, 36, 34, 36, 39, 35, 36},
			6:  []int{41, 42, 40, 38, 39, 42, 39, 40},
			7:  []int{45, 45, 44, 42, 43, 45, 43, 43},
			8:  []int{48, 48, 48, 46, 47, 48, 46, 46},
			9:  []int{51, 51, 52, 50, 51, 52, 50, 50},
			10: []int{54, 55, 56, 53, 55, 55, 54, 53},
			11: []int{57, 58, 60, 57, 59, 58, 57, 57},
			12: []int{60, 61, 64, 61, 63, 61, 61, 60},
			13: []int{63, 64, 68, 65, 67, 64, 65, 64},
			14: []int{67, 68, 72, 69, 71, 68, 68, 67},
			15: []int{70, 71, 76, 73, 75, 71, 72, 70},
			16: []int{73, 74, 80, 77, 78, 74, 76, 74},
			17: []int{76, 77, 84, 81, 82, 77, 79, 77},
			18: []int{79, 81, 89, 85, 86, 81, 83, 81},
		},
		"Дистанцирование": {
			0:  []int{24, 22, 14, 21, 21, 24, 19, 19},
			1:  []int{27, 26, 18, 24, 24, 27, 22, 23},
			2:  []int{30, 29, 21, 27, 28, 30, 26, 26},
			3:  []int{33, 32, 25, 31, 31, 33, 29, 29},
			4:  []int{37, 35, 29, 34, 34, 36, 32, 32},
			5:  []int{40, 38, 33, 38, 38, 39, 36, 36},
			6:  []int{43, 41, 37, 41, 41, 42, 39, 39},
			7:  []int{46, 44, 41, 44, 45, 45, 43, 42},
			8:  []int{49, 48, 45, 48, 48, 48, 46, 45},
			9:  []int{52, 51, 49, 51, 51, 52, 50, 49},
			10: []int{56, 54, 53, 54, 55, 55, 53, 52},
			11: []int{59, 57, 56, 58, 58, 58, 57, 55},
			12: []int{62, 60, 60, 61, 61, 61, 60, 59},
			13: []int{65, 63, 64, 64, 65, 64, 64, 62},
			14: []int{68, 67, 68, 68, 68, 67, 67, 65},
			15: []int{71, 70, 72, 71, 72, 70, 71, 68},
			16: []int{75, 73, 76, 74, 75, 73, 74, 72},
			17: []int{78, 76, 80, 78, 78, 76, 78, 75},
			18: []int{81, 79, 84, 81, 82, 79, 81, 78},
		},
		"Самоконтроль": {
			0:  []int{7, 8, 3, 9, 4, 6, 1, 8},
			1:  []int{10, 11, 7, 12, 8, 9, 4, 11},
			2:  []int{13, 14, 10, 15, 11, 12, 8, 14},
			3:  []int{16, 17, 14, 18, 15, 15, 12, 17},
			4:  []int{19, 20, 17, 21, 18, 19, 15, 20},
			5:  []int{22, 23, 21, 24, 21, 22, 19, 24},
			6:  []int{25, 26, 24, 27, 25, 25, 23, 27},
			7:  []int{28, 29, 28, 30, 28, 28, 27, 30},
			8:  []int{31, 32, 31, 33, 32, 31, 30, 33},
			9:  []int{34, 35, 35, 36, 35, 34, 34, 36},
			10: []int{37, 38, 39, 39, 38, 37, 38, 39},
			11: []int{41, 41, 42, 42, 42, 40, 41, 42},
			12: []int{44, 44, 46, 45, 45, 43, 45, 45},
			13: []int{47, 47, 49, 48, 49, 46, 49, 49},
			14: []int{50, 50, 53, 51, 52, 49, 53, 52},
			15: []int{53, 52, 56, 54, 55, 52, 56, 55},
			16: []int{56, 55, 60, 57, 59, 55, 60, 58},
			17: []int{59, 58, 63, 60, 62, 58, 64, 61},
			18: []int{62, 61, 67, 63, 66, 62, 67, 64},
			19: []int{65, 64, 70, 66, 69, 65, 71, 67},
			20: []int{68, 67, 74, 69, 72, 68, 75, 70},
			21: []int{71, 70, 77, 72, 76, 71, 78, 73},
		},
		"Поиск соц. поддержки": {
			0:  []int{18, 10, 12, 19, 15, 13, 12, 20},
			1:  []int{21, 14, 15, 22, 19, 16, 15, 22},
			2:  []int{24, 17, 19, 25, 22, 19, 19, 25},
			3:  []int{27, 21, 22, 28, 25, 22, 21, 28},
			4:  []int{30, 24, 25, 31, 28, 25, 25, 31},
			5:  []int{33, 28, 29, 34, 31, 29, 28, 33},
			6:  []int{36, 31, 32, 37, 34, 32, 31, 36},
			7:  []int{39, 35, 36, 40, 37, 35, 34, 39},
			8:  []int{42, 38, 39, 43, 40, 38, 37, 41},
			9:  []int{45, 42, 43, 46, 43, 41, 40, 44},
			10: []int{48, 45, 46, 48, 46, 44, 44, 47},
			11: []int{51, 49, 49, 51, 49, 47, 47, 49},
			12: []int{54, 52, 53, 54, 52, 51, 50, 52},
			13: []int{57, 56, 56, 57, 55, 54, 53, 55},
			14: []int{60, 59, 60, 60, 58, 57, 56, 57},
			15: []int{63, 63, 63, 63, 62, 60, 59, 60},
			16: []int{66, 66, 67, 66, 65, 63, 63, 63},
			17: []int{69, 70, 70, 69, 68, 66, 66, 66},
			18: []int{72, 76, 73, 72, 71, 69, 69, 68},
		},
		"Принятие ответственности": {
			0:  []int{20, 20, 18, 22, 17, 18, 16, 19},
			1:  []int{24, 24, 22, 25, 21, 22, 20, 23},
			2:  []int{28, 28, 27, 29, 25, 26, 25, 27},
			3:  []int{31, 31, 31, 33, 30, 30, 30, 31},
			4:  []int{35, 35, 35, 36, 34, 34, 34, 35},
			5:  []int{39, 39, 39, 40, 39, 38, 39, 39},
			6:  []int{43, 43, 43, 44, 43, 42, 43, 43},
			7:  []int{47, 46, 48, 47, 47, 46, 48, 47},
			8:  []int{51, 50, 52, 51, 52, 50, 52, 51},
			9:  []int{55, 54, 56, 55, 56, 54, 57, 55},
			10: []int{59, 58, 60, 58, 61, 58, 61, 59},
			11: []int{63, 61, 65, 62, 65, 62, 66, 63},
			12: []int{67, 65, 69, 66, 60, 66, 70, 67},
		},
		"Бегство-избегание": {
			0:  []int{27, 27, 23, 21, 21, 25, 16, 23},
			1:  []int{29, 29, 25, 24, 24, 28, 19, 25},
			2:  []int{31, 31, 28, 26, 26, 30, 22, 28},
			3:  []int{35, 35, 31, 29, 29, 32, 25, 30},
			4:  []int{36, 36, 33, 32, 31, 35, 28, 32},
			5:  []int{38, 38, 36, 35, 34, 37, 31, 35},
			6:  []int{40, 41, 38, 37, 37, 39, 34, 37},
			7:  []int{43, 43, 41, 40, 39, 42, 37, 40},
			8:  []int{45, 45, 44, 43, 42, 44, 40, 42},
			9:  []int{47, 48, 46, 46, 45, 46, 43, 44},
			10: []int{49, 50, 49, 49, 47, 48, 46, 47},
			11: []int{52, 52, 51, 51, 50, 51, 49, 49},
			12: []int{54, 54, 54, 54, 53, 53, 52, 52},
			13: []int{56, 57, 57, 57, 55, 55, 55, 54},
			14: []int{58, 59, 59, 60, 58, 58, 58, 57},
			15: []int{61, 61, 62, 62, 61, 60, 61, 59},
			16: []int{63, 64, 65, 65, 63, 62, 64, 61},
			17: []int{65, 66, 67, 68, 66, 65, 67, 64},
			18: []int{67, 68, 70, 71, 68, 67, 70, 66},
			19: []int{69, 71, 72, 73, 71, 69, 73, 69},
			20: []int{72, 73, 75, 76, 74, 71, 75, 71},
			21: []int{74, 75, 78, 79, 76, 74, 78, 74},
			22: []int{76, 77, 80, 82, 79, 76, 81, 76},
			23: []int{78, 80, 83, 85, 82, 78, 84, 78},
			24: []int{81, 82, 85, 87, 84, 81, 87, 81},
		},
		"Планирование решения": {
			0:  []int{11, 7, 9, 9, 13, 11, 6, 12},
			1:  []int{14, 11, 13, 12, 16, 14, 9, 15},
			2:  []int{17, 14, 16, 15, 19, 18, 13, 18},
			3:  []int{20, 17, 19, 19, 22, 21, 16, 21},
			4:  []int{24, 21, 22, 22, 25, 24, 20, 24},
			5:  []int{27, 24, 25, 25, 28, 27, 24, 28},
			6:  []int{30, 28, 29, 28, 32, 31, 27, 31},
			7:  []int{33, 31, 32, 32, 35, 35, 31, 34},
			8:  []int{37, 34, 35, 35, 38, 38, 35, 37},
			9:  []int{40, 38, 38, 38, 41, 42, 38, 40},
			10: []int{43, 41, 42, 42, 44, 45, 42, 43},
			11: []int{46, 45, 45, 45, 47, 49, 45, 46},
			12: []int{49, 48, 48, 48, 50, 53, 49, 49},
			13: []int{53, 51, 51, 51, 53, 56, 53, 52},
			14: []int{56, 55, 55, 55, 57, 60, 56, 56},
			15: []int{59, 58, 58, 58, 60, 63, 60, 59},
			16: []int{62, 62, 61, 61, 63, 67, 63, 62},
			17: []int{66, 65, 64, 65, 66, 71, 67, 65},
			18: []int{69, 68, 68, 68, 69, 68, 71, 68},
		},
		"Положительная переоценка": {
			0:  []int{18, 14, 13, 20, 18, 14, 13, 17},
			1:  []int{21, 17, 16, 22, 21, 17, 16, 20},
			2:  []int{24, 20, 19, 25, 24, 20, 19, 22},
			3:  []int{26, 23, 22, 27, 26, 23, 22, 25},
			4:  []int{29, 26, 26, 30, 29, 25, 24, 27},
			5:  []int{32, 29, 29, 33, 32, 28, 27, 30},
			6:  []int{34, 31, 32, 35, 34, 31, 30, 32},
			7:  []int{37, 34, 35, 38, 37, 34, 33, 35},
			8:  []int{40, 37, 38, 40, 40, 36, 36, 37},
			9:  []int{42, 40, 41, 43, 42, 39, 39, 40},
			10: []int{45, 43, 45, 45, 45, 42, 42, 42},
			11: []int{48, 46, 48, 48, 48, 45, 45, 45},
			12: []int{51, 49, 51, 50, 50, 47, 47, 47},
			13: []int{53, 52, 54, 53, 53, 50, 50, 50},
			14: []int{56, 54, 57, 55, 56, 53, 53, 52},
			15: []int{59, 57, 61, 58, 58, 56, 56, 55},
			16: []int{61, 60, 64, 60, 61, 58, 58, 57},
			17: []int{64, 63, 67, 63, 64, 61, 61, 60},
			18: []int{67, 66, 70, 65, 66, 64, 64, 62},
			19: []int{69, 69, 73, 68, 69, 67, 67, 65},
			20: []int{72, 72, 76, 70, 72, 69, 69, 67},
			21: []int{75, 75, 80, 73, 74, 72, 72, 70},
		},
	}
	scaleWCQKeys = map[string]map[string]int{
		"Конфронтация": {
			"max_value": 90,
			"2":         1,
			"3":         1,
			"13":        1,
			"21":        1,
			"26":        1,
			"37":        1,
		},
		"Дистанцирование": {
			"max_value": 90,
			"8":         1,
			"9":         1,
			"11":        1,
			"16":        1,
			"32":        1,
			"35":        1,
		},
		"Самоконтроль": {
			"max_value": 90,
			"6":         1,
			"10":        1,
			"27":        1,
			"34":        1,
			"44":        1,
			"49":        1,
			"50":        1,
		},
		"Поиск соц. поддержки": {
			"max_value": 90,
			"4":         1,
			"14":        1,
			"17":        1,
			"24":        1,
			"33":        1,
			"36":        1,
		},
		"Принятие ответственности": {
			"max_value": 90,
			"5":         1,
			"19":        1,
			"22":        1,
			"42":        1,
		},
		"Бегство-избегание": {
			"max_value": 90,
			"7":         1,
			"12":        1,
			"25":        1,
			"31":        1,
			"38":        1,
			"41":        1,
			"46":        1,
			"47":        1,
		},
		"Планирование решения": {
			"max_value": 90,
			"1":         1,
			"20":        1,
			"30":        1,
			"39":        1,
			"40":        1,
			"43":        1,
		},
		"Положительная переоценка": {
			"max_value": 90,
			"15":        1,
			"18":        1,
			"23":        1,
			"28":        1,
			"29":        1,
			"45":        1,
			"48":        1,
		},
	}
)

func ageAndSexResolver(age, score int, sex, field string) int {
	var sexInt int
	switch strings.ToLower(sex) {
	case "male":
		sexInt = 1
	case "female":
		sexInt = 0
	default:
		log.Fatalln("'sex' can be only 'male' or 'female'")
	}
	
	switch {
	case age <= 20:
		return standardPoints[field][score][sexInt*4]
	case age >= 21 && age <= 30:
		return standardPoints[field][score][1+sexInt*4]
	case age >= 31 && age <= 45:
		return standardPoints[field][score][2+sexInt*4]
	default: // age >= 46
		return standardPoints[field][score][3+sexInt*4]
	}
}

func WCQHandler(s *types.SurveyResults) []byte {
	result := map[string]map[string]int{
		"Конфронтация": {
			"max_value": 90,
			"value":     0,
		},
		"Дистанцирование": {
			"max_value": 90,
			"value":     0,
		},
		"Самоконтроль": {
			"max_value": 90,
			"value":     0,
		},
		"Поиск соц. поддержки": {
			"max_value": 90,
			"value":     0,
		},
		"Принятие ответственности": {
			"max_value": 90,
			"value":     0,
		},
		"Бегство-избегание": {
			"max_value": 90,
			"value":     0,
		},
		"Планирование решения": {
			"max_value": 90,
			"value":     0,
		},
		"Положительная переоценка": {
			"max_value": 90,
			"value":     0,
		},
	}
	
	for field, keysMap := range scaleWCQKeys {
		result[field]["max_value"] = keysMap["max_value"]
	}
	
	// answer: 0, 1, 2, 3
	for questionID, answer := range s.Picked {
		for field, keysMap := range scaleWCQKeys {
			if value, ok := keysMap[fmt.Sprintf("%d", questionID-90)]; ok && value == 1 && field != "max_value" {
				result[field]["value"] += answer
			}
		}
	}
	
	for field, resMap := range result {
		result[field]["value"] = ageAndSexResolver(s.Age, resMap["value"], s.Sex, field)
	}
	
	fmt.Println("WCQ result: ", result)
	
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	return resultJSON
}
