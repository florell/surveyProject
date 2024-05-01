package handlers

import (
	"encoding/json"
	"log"
	types "psychward/src"
	"reflect"
	"sort"
)

type MLOField struct {
	name       string
	num        string
	MaxValue   int `json:"max_value"`
	Value      int `json:"value"`
	Correction int `json:"correction"`
	TScore     int `json:"tscore"`
}

type MLOResult struct {
	Code                       string   `json:"code"`
	ConfidenceScale            MLOField `json:"Шкала достоверности (L)"`
	ReliabilityScale           MLOField `json:"Шкала надежности (F)"`
	CorrectionScale            MLOField `json:"Шкала коррекции (К)"`
	HypochondriaScale          MLOField `json:"Шкала ипохондрии (Hs)"`
	DepressionScale            MLOField `json:"Шкала депрессии (D)"`
	HysteriaScale              MLOField `json:"Шкала истерии (Ну)"`
	PsychopathyScale           MLOField `json:"Шкала психопатии (Pd)"`
	MasculinityFemininityScale MLOField `json:"Шкала мужественности-женственности (Mf)"`
	ParanoiaScale              MLOField `json:"Шкала паранойяльности (Ра)"`
	PsychastheniaScale         MLOField `json:"Шкала психастении (Рt)"`
	SchizophreniaScale         MLOField `json:"Шкала шизоидности (Sc)"`
	HypomaniaScale             MLOField `json:"Шкала гипомании (Ma)"`
	SocialIntroversionScale    MLOField `json:"Шкала социальной интроверсии (Si)"`
}

var TScore = map[string][]int{
	"L":  {35, 39, 43, 47, 52, 56, 60, 64, 67, 71, 76, 81, 87, 94, 101, 107},
	"F":  {32, 35, 38, 41, 44, 47, 50, 53, 56, 59, 62, 65, 69, 73, 76, 79, 82, 85, 88, 91, 94, 97, 100, 103, 106, 109},
	"K":  {44, 50, 55, 60, 65, 70, 75, 80, 85, 91},
	"Hs": {37, 40, 45, 47, 49, 51, 54, 57, 60, 65, 70, 75, 80, 85, 90, 95, 100, 105, 110},
	"D":  {35, 39, 42, 45, 48, 51, 56, 58, 63, 65, 70, 74, 77, 80, 83, 86, 89, 92, 95, 98, 105, 110},
	"Hy": {30, 33, 36, 39, 42, 45, 47, 49, 52, 54, 56, 60, 62, 65, 68, 70, 73, 75, 78, 82, 85, 89, 95, 100, 105, 111, 118},
	"Pd": {30, 35, 38, 42, 45, 48, 52, 55, 60, 62, 65, 67, 69, 72, 75, 78, 81, 84, 87, 90, 93, 96, 99, 102, 105, 108, 112, 116, 119},
	"Mf": {35, 40, 45, 50, 55, 59, 63, 71, 78, 88},
	"Pa": {27, 31, 36, 40, 43, 46, 49, 52, 55, 58, 63, 67, 71, 74, 78, 81, 84, 87, 91, 95, 100},
	"Pt": {35, 37, 40, 43, 46, 49, 52, 55, 57, 59, 61, 63, 65, 67, 70, 71, 72, 73, 74, 76, 77, 78, 81, 84, 86, 88, 91, 93, 97, 101, 105},
	"Sc": {27, 31, 35, 39, 42, 45, 48, 51, 53, 55, 57, 59, 61, 63, 65, 67, 69, 71, 73, 75, 77, 79, 81, 83, 85, 87, 89, 91, 93, 95, 97, 99, 101, 103, 105, 107, 109, 111, 113, 115, 117, 119, 120},
	"Ma": {33, 36, 39, 42, 47, 50, 53, 55, 57, 62, 68, 71, 74, 77, 80, 83, 86, 89, 92, 95, 98, 105, 109, 115, 119},
	"Si": {35, 39, 43, 47, 51, 55, 59, 64, 68, 73, 77, 82, 86},
}

var scoresTable = map[int]map[string]int{
	9: {"+1 K": 9, "+0.5 K": 5, "+0.4 K": 4, "+0.2 K": 2},
	8: {"+1 K": 8, "+0.5 K": 4, "+0.4 K": 3, "+0.2 K": 2},
	7: {"+1 K": 7, "+0.5 K": 4, "+0.4 K": 3, "+0.2 K": 1},
	6: {"+1 K": 6, "+0.5 K": 3, "+0.4 K": 2, "+0.2 K": 1},
	5: {"+1 K": 5, "+0.5 K": 3, "+0.4 K": 2, "+0.2 K": 1},
	4: {"+1 K": 4, "+0.5 K": 2, "+0.4 K": 2, "+0.2 K": 1},
	3: {"+1 K": 3, "+0.5 K": 2, "+0.4 K": 1, "+0.2 K": 1},
	2: {"+1 K": 2, "+0.5 K": 1, "+0.4 K": 1, "+0.2 K": 0},
	1: {"+1 K": 1, "+0.5 K": 1, "+0.4 K": 1, "+0.2 K": 0},
	0: {"+1 K": 0, "+0.5 K": 0, "+0.4 K": 0, "+0.2 K": 0},
}

func getValue(yes []int, no []int, s *types.SurveyResults) int {
	result := 0
	
	for _, value := range yes {
		result += s.Picked[value+90+50+26+21+30]
	}
	for _, value := range no {
		result -= s.Picked[value+90+50+26+21+30]
	}
	
	return result
}

func getTScore(field *MLOField) int {
	return TScore[field.name][field.Value]
}

func getCorrection(score int, correction string) int {
	return scoresTable[score][correction]
}

func getCode(result *MLOResult) string {
	v := reflect.ValueOf(&result).Elem()
	
	fieldsDigits := make([]MLOField, 0)
	fieldsWords := make([]MLOField, 0)
	
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i).Addr().Interface().(MLOField)
		if field.name == field.num {
			fieldsWords = append(fieldsWords, field)
		} else {
			fieldsDigits = append(fieldsDigits, field)
		}
	}
	
	sort.SliceStable(fieldsDigits, func(i, j int) bool {
		//if fieldsDigits[i].TScore == fieldsDigits[j].TScore {
		//	return fieldsDigits[i].num < fieldsDigits[j].num
		//}
		return fieldsDigits[i].TScore < fieldsDigits[j].TScore
	})
	
	sort.SliceStable(fieldsWords, func(i, j int) bool {
		//if fieldsWords[i].TScore == fieldsWords[j].TScore {
		//	return fieldsWords[i].num < fieldsWords[j].num
		//}
		return fieldsWords[i].TScore < fieldsWords[j].TScore
	})
	
	getSymbol := func(score int) string {
		switch {
		case score > 120:
			return "!!"
		case score > 110:
			return "!"
		case score > 100:
			return "**"
		case score > 90:
			return "*"
		case score > 80:
			return "”"
		case score > 70:
			return "‘"
		case score > 60:
			return "–"
		case score > 50:
			return "/"
		case score > 40:
			return ":"
		case score > 30:
			return "#"
		default:
			return ""
		}
	}
	
	code := ""
	for i := 0; i < len(fieldsDigits); i++ {
		code += fieldsDigits[i].num
		if fieldsDigits[i].TScore > fieldsDigits[i+1].TScore || i == len(fieldsDigits)-1 {
			code += getSymbol(fieldsDigits[i].TScore)
		}
	}
	for i := 0; i < len(fieldsWords); i++ {
		code += fieldsWords[i].num
		if fieldsWords[i].TScore > fieldsWords[i+1].TScore || i == len(fieldsWords)-1 {
			code += getSymbol(fieldsWords[i].TScore)
		}
	}
	
	return code
}

func MLOHandler(s *types.SurveyResults) []byte {
	result := MLOResult{}
	
	result.ConfidenceScale.Value = getValue([]int{}, []int{1, 10, 31, 69, 78, 92, 101, 116, 128, 148}, s)
	result.ConfidenceScale.MaxValue = 10
	result.ConfidenceScale.name = "L"
	result.ConfidenceScale.num = "L"
	result.ConfidenceScale.TScore = getTScore(&result.ConfidenceScale)
	
	result.ReliabilityScale.Value = getValue(
		[]int{4, 8, 11, 18, 20, 22, 37, 41, 47, 60, 72, 82, 84, 86, 91, 96, 98, 103, 115, 153},
		[]int{2, 25, 43, 44, 53}, s)
	result.ReliabilityScale.MaxValue = 25
	result.ReliabilityScale.name = "F"
	result.ReliabilityScale.num = "F"
	result.ReliabilityScale.TScore = getTScore(&result.ReliabilityScale)
	
	result.CorrectionScale.Value = getValue([]int{35}, []int{15, 46, 48, 64, 73, 90, 102, 151}, s)
	result.CorrectionScale.MaxValue = 9
	result.CorrectionScale.name = "K"
	result.CorrectionScale.num = "K"
	result.CorrectionScale.TScore = getTScore(&result.CorrectionScale)
	
	result.HypochondriaScale.Value = getValue([]int{17, 67}, []int{2, 3, 5, 23, 38, 53, 55, 58, 62, 75, 93}, s)
	result.HypochondriaScale.MaxValue = 13
	result.HypochondriaScale.name = "Hs"
	result.HypochondriaScale.num = "1"
	result.HypochondriaScale.Correction = getCorrection(result.HypochondriaScale.Value, "+0.5 K")
	result.HypochondriaScale.Value += result.HypochondriaScale.Correction
	result.HypochondriaScale.TScore = getTScore(&result.HypochondriaScale)
	
	result.DepressionScale.Value = getValue([]int{16, 17, 30, 39, 46},
		[]int{5, 14, 23, 26, 27, 32, 34, 50, 52, 53, 54, 55, 67, 68, 77, 102}, s)
	result.DepressionScale.MaxValue = 21
	result.DepressionScale.name = "D"
	result.DepressionScale.num = "2"
	result.DepressionScale.TScore = getTScore(&result.DepressionScale)
	
	result.HysteriaScale.Value = getValue([]int{11, 17, 20, 21, 28, 65, 67},
		[]int{2, 3, 23, 33, 38, 42, 45, 48, 53, 58, 61, 62, 64, 75, 88, 90, 95, 97, 99}, s)
	result.HysteriaScale.MaxValue = 26
	result.HysteriaScale.name = "Hy"
	result.HysteriaScale.num = "3"
	result.HysteriaScale.TScore = getTScore(&result.HysteriaScale)
	
	result.PsychopathyScale.Value = getValue([]int{6, 8, 11, 12, 14, 41, 42, 56, 72, 81, 82, 91, 114},
		[]int{13, 35, 45, 48, 55, 79, 90, 97, 100, 102}, s)
	result.PsychopathyScale.MaxValue = 23 + 4
	result.PsychopathyScale.name = "Pd"
	result.PsychopathyScale.num = "4"
	result.PsychopathyScale.Correction = getCorrection(result.PsychopathyScale.Value, "+0.4 K")
	result.PsychopathyScale.Value += result.PsychopathyScale.Correction
	result.PsychopathyScale.TScore = getTScore(&result.PsychopathyScale)
	
	result.MasculinityFemininityScale.Value = getValue([]int{63, 66, 73}, []int{9, 43, 50, 74, 86, 87}, s)
	result.MasculinityFemininityScale.MaxValue = 9
	result.MasculinityFemininityScale.name = "Mf"
	result.MasculinityFemininityScale.num = "5"
	result.MasculinityFemininityScale.TScore = getTScore(&result.MasculinityFemininityScale)
	
	result.ParanoiaScale.Value = getValue([]int{4, 7, 8, 10, 18, 39, 43, 46, 48, 98, 104, 125, 150, 152},
		[]int{33, 42, 84, 137, 145, 155}, s)
	result.ParanoiaScale.MaxValue = 20
	result.ParanoiaScale.name = "Pa"
	result.ParanoiaScale.num = "6"
	result.ParanoiaScale.TScore = getTScore(&result.ParanoiaScale)
	
	result.PsychastheniaScale.Value = getValue([]int{7, 10, 11, 16, 28, 30, 37, 41, 67, 73, 80, 88, 103, 104, 110, 117, 120, 122, 123},
		[]int{2, 52}, s)
	result.PsychastheniaScale.MaxValue = 21 + 9
	result.PsychastheniaScale.name = "Pt"
	result.PsychastheniaScale.num = "7"
	result.PsychastheniaScale.Correction = getCorrection(result.PsychastheniaScale.Value, "+1 K")
	result.PsychastheniaScale.Value += result.PsychastheniaScale.Correction
	result.PsychastheniaScale.TScore = getTScore(&result.PsychastheniaScale)
	
	result.SchizophreniaScale.Value = getValue([]int{4, 6, 7, 8, 10, 11, 12, 14, 16, 21, 24, 36, 39, 56, 60, 63, 70, 80, 89, 98, 103, 105, 106, 108, 111, 119, 123, 124},
		[]int{13, 38, 44, 66, 107}, s)
	result.SchizophreniaScale.MaxValue = 33 + 9
	result.SchizophreniaScale.name = "Sc"
	result.SchizophreniaScale.num = "8"
	result.SchizophreniaScale.Correction = getCorrection(result.SchizophreniaScale.Value, "+1 K")
	result.SchizophreniaScale.Value += result.SchizophreniaScale.Correction
	result.SchizophreniaScale.TScore = getTScore(&result.SchizophreniaScale)
	
	result.HypomaniaScale.Value = getValue([]int{6, 7, 27, 36, 42, 49, 56, 59, 76, 77, 80, 89, 90, 93, 95},
		[]int{40, 43, 64, 96}, s)
	result.HypomaniaScale.MaxValue = 19 + 2
	result.HypomaniaScale.name = "Ma"
	result.HypomaniaScale.num = "9"
	result.HypomaniaScale.Correction = getCorrection(result.HypomaniaScale.Value, "+0.2 K")
	result.HypomaniaScale.Value += result.HypomaniaScale.Correction
	result.HypomaniaScale.TScore = getTScore(&result.HypomaniaScale)
	
	result.SocialIntroversionScale.Value = getValue([]int{64, 85, 126, 160, 163}, []int{12, 49, 90, 74, 144, 147, 159}, s)
	result.SocialIntroversionScale.MaxValue = 12
	result.SocialIntroversionScale.name = "Si"
	result.SocialIntroversionScale.num = "0"
	result.SocialIntroversionScale.TScore = getTScore(&result.SocialIntroversionScale)
	
	result.Code = getCode(&result)
	
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	return resultJSON
}
