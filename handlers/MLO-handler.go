package handlers

import (
	"encoding/json"
	"fmt"
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
	LAP                        MLOField `json:"Личностный адаптационный потенциал (ЛАП)"`
	PR                         MLOField `json:"Поведенческая регуляция (ПР)"`
	KP                         MLOField `json:"Коммуникативный потенциал (КП)"`
	MN                         MLOField `json:"Моральная нормативность (МН)"`
	DAN                        MLOField `json:"Дезадаптационные нарушения"`
	AS                         MLOField `json:"Астенические реакции и состояния (АС)"`
	PS                         MLOField `json:"Психотические реакции и состояния (ПС)"`
	D                          MLOField `json:"Шкала искренности (достоверность)"`
	Description                string   `json:"description"`
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
		if s.Picked[value+90+50+26+21+30] == 1 {
			result++
		}
	}
	for _, value := range no {
		if s.Picked[value+90+50+26+21+30] == -1 {
			result++
		}
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
	v := reflect.ValueOf(*result)

	fieldsDigits := make([]MLOField, 0)
	fieldsWords := make([]MLOField, 0)

	for i := 1; i < v.NumField()-1; i++ {
		field := v.Field(i).Interface().(MLOField)
		if field.name == field.num {
			fieldsWords = append(fieldsWords, field)
		} else if field.name != "" {
			fieldsDigits = append(fieldsDigits, field)
		}
	}

	sort.SliceStable(fieldsDigits, func(i, j int) bool {
		//if fieldsDigits[i].TScore == fieldsDigits[j].TScore {
		//	return fieldsDigits[i].num < fieldsDigits[j].num
		//}
		return fieldsDigits[i].TScore > fieldsDigits[j].TScore
	})

	sort.SliceStable(fieldsWords, func(i, j int) bool {
		//if fieldsWords[i].TScore == fieldsWords[j].TScore {
		//	return fieldsWords[i].num < fieldsWords[j].num
		//}
		return fieldsWords[i].TScore > fieldsWords[j].TScore
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
		if i == len(fieldsDigits)-1 || getSymbol(fieldsDigits[i].TScore) != getSymbol(fieldsDigits[i+1].TScore) {
			code += getSymbol(fieldsDigits[i].TScore)
		}
	}
	code += "   "
	for i := 0; i < len(fieldsWords); i++ {
		code += fieldsWords[i].num
		if i == len(fieldsWords)-1 || getSymbol(fieldsWords[i].TScore) != getSymbol(fieldsWords[i+1].TScore) {
			code += getSymbol(fieldsWords[i].TScore)
		}
	}

	return code
}

func PR(count int) int {
	result := 0

	switch {
	case count >= 46:
		result = 1
	case count >= 38 && count <= 45:
		result = 2
	case count >= 30 && count <= 37:
		result = 3
	case count >= 22 && count <= 29:
		result = 4
	case count >= 16 && count <= 21:
		result = 5
	case count >= 13 && count <= 15:
		result = 6
	case count >= 9 && count <= 12:
		result = 7
	case count >= 6 && count <= 8:
		result = 8
	case count >= 4 && count <= 5:
		result = 9
	case count <= 3:
		result = 10
	}

	return result
}

func KP(count int) int {
	result := 0

	switch {
	case count >= 27 && count <= 31:
		result = 1
	case count >= 22 && count <= 26:
		result = 2
	case count >= 17 && count <= 21:
		result = 3
	case count >= 13 && count <= 16:
		result = 4
	case count >= 10 && count <= 12:
		result = 5
	case count >= 7 && count <= 9:
		result = 6
	case count >= 5 && count <= 6:
		result = 7
	case count >= 3 && count <= 4:
		result = 8
	case count >= 1 && count <= 2:
		result = 9
	case count == 0:
		result = 10
	}

	return result
}

func MN(count int) int {
	result := 0

	switch {
	case count >= 18:
		result = 1
	case count >= 15 && count <= 17:
		result = 2
	case count >= 12 && count <= 14:
		result = 3
	case count >= 10 && count <= 11:
		result = 4
	case count >= 7 && count <= 9:
		result = 5
	case count >= 5 && count <= 6:
		result = 6
	case count >= 3 && count <= 4:
		result = 7
	case count == 2:
		result = 8
	case count == 1:
		result = 9
	case count == 0:
		result = 10
	}

	return result
}

func LAP(age, count int) int {
	result := 0

	if age >= 21 {
		switch {
		case count <= 5:
			result = 10
		case count >= 6 && count <= 10:
			result = 9
		case count >= 11 && count <= 15:
			result = 8
		case count >= 16 && count <= 21:
			result = 7
		case count >= 22 && count <= 27:
			result = 6
		case count >= 28 && count <= 32:
			result = 5
		case count >= 33 && count <= 39:
			result = 4
		case count >= 40 && count <= 50:
			result = 3
		case count >= 51 && count <= 60:
			result = 2
		case count >= 61:
			result = 1
		}
	} else {
		switch {
		case count <= 11:
			result = 10
		case count >= 12 && count <= 14:
			result = 9
		case count >= 15 && count <= 17:
			result = 8
		case count >= 18 && count <= 22:
			result = 7
		case count >= 23 && count <= 26:
			result = 6
		case count >= 27 && count <= 32:
			result = 5
		case count >= 33 && count <= 39:
			result = 4
		case count >= 40 && count <= 46:
			result = 3
		case count >= 47 && count <= 57:
			result = 2
		case count >= 58:
			result = 1
		}
	}

	return result
}

func AS(count int) int {
	result := 0

	switch {
	case count >= 36:
		result = 1
	case count >= 30 && count <= 35:
		result = 2
	case count >= 22 && count <= 29:
		result = 3
	case count >= 16 && count <= 21:
		result = 4
	case count >= 10 && count <= 15:
		result = 5
	case count >= 7 && count <= 9:
		result = 6
	case count >= 5 && count <= 6:
		result = 7
	case count >= 3 && count <= 4:
		result = 8
	case count == 2:
		result = 9
	case count >= 0 && count <= 1:
		result = 10
	}

	return result
}

func PS(count int) int {
	result := 0

	switch {
	case count >= 27:
		result = 1
	case count >= 22 && count <= 26:
		result = 2
	case count >= 16 && count <= 21:
		result = 3
	case count >= 13 && count <= 15:
		result = 4
	case count >= 8 && count <= 12:
		result = 5
	case count >= 6 && count <= 7:
		result = 6
	case count >= 4 && count <= 5:
		result = 7
	case count >= 2 && count <= 3:
		result = 8
	case count == 1:
		result = 9
	case count >= 0 && count <= 1:
		result = 10
	}

	return result
}

func DAN(count int) int {
	result := 0

	switch {
	case count >= 51:
		result = 1
	case count >= 43 && count <= 50:
		result = 2
	case count >= 36 && count <= 42:
		result = 3
	case count >= 31 && count <= 35:
		result = 4
	case count >= 21 && count <= 30:
		result = 5
	case count >= 16 && count <= 20:
		result = 6
	case count >= 11 && count <= 15:
		result = 7
	case count >= 6 && count <= 10:
		result = 8
	case count >= 3 && count <= 5:
		result = 9
	case count >= 0 && count <= 2:
		result = 10
	}

	return result
}

func getMLODescription(result MLOResult) string {
	description := ""

	switch {
	case result.LAP.Value >= 1 && result.LAP.Value <= 3:
		description += "Группа 4: Группа сниженной адаптации. Эта группа обладает признаками явных акцентуаций характера и некоторыми признаками психопатий, а психическое состояние можно охарактеризовать как пограничное. Процесс адаптации протекает тяжело. Возможны нервно-психические срывы, длительные нарушения функционального состояния. Лица этой группы обладают низкой нервно-психической устойчивостью, конфликтны, могут допускать делинквентные поступки."
	case result.LAP.Value >= 4 && result.LAP.Value <= 6:
		description += "Группа 3: Группа удовлетворительной адаптации. Большинство лиц этой группы обладают признаками различных акцентуаций, которые в привычных условиях частично компенсированы и могут проявляться при смене деятельности. Поэтому успех адаптации во многом зависит от внешних условий среды. Эти лица, как правило, обладают невысокой эмоциональной устойчивостью. Процесс социализации осложнён, возможны асоциальные срывы, проявление агрессивности и конфликтности. Функциональное состояние в начальные этапы адаптации может быть нарушено. Лица этой группы требуют постоянного контроля."
	case result.LAP.Value >= 7 && result.LAP.Value <= 8:
		description += "Группа 2: Группа хороших адаптационных способностей. Лица той группы легко адаптируются к новым условиям деятельности, быстро «входят» в новый коллектив, достаточно легко и адекватно ориентируются в ситуации, быстро вырабатывают стратегию своего поведения и социализации. Как правило, не конфликтны, обладают высокой эмоциональной устойчивостью. Функциональное состояние лиц этой группы в период адаптации остаётся в пределах нормы, работоспособность сохраняется."
	case result.LAP.Value >= 9 && result.LAP.Value <= 10:
		description += "Группа 1: Группа хороших адаптационных способностей. Лица той группы легко адаптируются к новым условиям деятельности, быстро «входят» в новый коллектив, достаточно легко и адекватно ориентируются в ситуации, быстро вырабатывают стратегию своего поведения и социализации. Как правило, не конфликтны, обладают высокой эмоциональной устойчивостью. Функциональное состояние лиц этой группы в период адаптации остаётся в пределах нормы, работоспособность сохраняется."
	}

	if result.AS.Value >= 1 && result.AS.Value <= 3 {
		description += "Шкала «Астенических реакций и состояний» (АС): Высокий уровень ситуационной тревожности, расстройства сна, ипохондрическая фиксация, повышенная утомляемость, истощаемость, слабость, резкое снижение способности к продолжительному физическому или умственному напряжению, низкая толерантность к неблагоприятным факторам профессиональной деятельности, особенно при чрезвычайных нагрузках, аффективная лабильность с преобладанием пониженного настроения, слезливость, гнетущая безысходность, тоска, хандра, восприятие настоящего окружения и своего будущего только в мрачном свете, наличие суицидальных мыслей, отсутствие мотивации к профессиональной деятельности и др."
	}

	if result.PS.Value >= 1 && result.PS.Value <= 3 {
		description += "Шкала «Психотические реакции и состояния» (ПС): Выраженное нервно-психическое напряжение, импульсивные реакции, приступы неконтролируемого гнева, ухудшение межличностных контактов, нарушение морально-нравственной ориентации, отсутствие стремления соблюдать общепринятые нормы поведения, групповых и корпоративных требований, делинквентное поведение, чрезмерная агрессивность, озлобленность, подозрительность, иногда: аутизм, деперсонализация, наличие галлюцинаций и др."
	}

	if result.DAN.Value >= 1 && result.DAN.Value <= 3 {
		description += "Интегральная шкала ДАН: Выраженные (достаточно выраженные) признаки дезадаптационных нарушений. Требуется консультация психиатра. Показана комплексная психологическая и фармакологическая коррекция."
	}

	return description
}

func MLOHandler(s *types.SurveyResults) []byte {
	fmt.Println(s.Picked)
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

	result.D.Value = getValue([]int{}, []int{1, 10, 11, 19, 31, 51, 69, 78, 92, 101, 116, 128, 138, 148}, s)
	result.D.MaxValue = 14

	result.PR.Value = PR(getValue([]int{4, 6, 7, 8, 11, 12, 15, 16, 17, 18, 20, 21, 28, 29, 30, 37, 39, 40, 41, 47, 57, 60, 63, 65, 67, 68, 70, 71, 73, 80, 82, 83, 84, 86, 89, 94, 95, 96, 98, 102, 103, 108, 109, 110, 111, 112, 113, 115, 117, 118, 119, 120, 122, 123, 124, 127, 129, 131, 135, 136, 137, 139, 143, 146, 149, 153, 154, 155, 156, 157, 158, 161, 162},
		[]int{2, 3, 5, 23, 25, 32, 38, 44, 45, 49, 52, 53, 54, 55, 58, 62, 66, 75, 87, 105, 127, 132, 134, 140}, s))
	result.PR.MaxValue = 10

	result.KP.Value = KP(getValue([]int{9, 24, 27, 33, 43, 46, 61, 64, 81, 88, 90, 99, 104, 106, 114, 121, 126, 133, 142, 151, 152},
		[]int{26, 34, 35, 48, 74, 85, 107, 130, 144, 147, 159}, s))
	result.KP.MaxValue = 10

	result.MN.Value = MN(getValue([]int{14, 22, 36, 42, 50, 56, 59, 72, 77, 79, 91, 93, 125, 141, 145, 150, 164, 165},
		[]int{13, 76, 97, 100, 160, 163}, s))
	result.MN.MaxValue = 10

	result.LAP.Value = LAP(s.Age, getValue([]int{4, 6, 7, 8, 9, 11, 12, 14, 15, 16, 17, 18, 20, 21, 22, 24, 27, 28, 29, 30, 33, 36, 37, 39, 40, 41, 42, 43, 46, 47, 50, 56, 57, 59, 60, 61, 63, 64, 65, 67, 68, 70, 71, 72, 73, 77, 79, 80, 81, 82, 83, 84, 86, 88, 89, 90, 91, 93, 94, 95, 96, 98, 99, 102, 103, 104, 106, 108, 109, 110, 111, 112, 113, 114, 115, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 129, 131, 133, 135, 136, 137, 139, 141, 142, 143, 145, 146, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 161, 162, 164, 165},
		[]int{2, 3, 5, 13, 23, 25, 26, 32, 34, 35, 38, 44, 45, 48, 49, 52, 53, 54, 55, 58, 62, 66, 74, 75, 76, 85, 87, 97, 100, 105, 107, 130, 132, 134, 140, 144, 147, 159, 160, 163}, s))
	result.LAP.MaxValue = 10

	result.AS.Value = AS(getValue([]int{6, 7, 12, 13, 14, 18, 27, 31, 32, 33, 34, 37, 41, 43, 46, 48, 49, 51, 52, 53, 55, 57, 58, 59, 60, 61, 63, 64, 71, 72, 73, 74},
		[]int{1, 2, 9, 11, 21, 25, 26, 30, 38, 42, 67}, s))
	result.AS.MaxValue = 10

	result.PS.Value = PS(getValue([]int{3, 4, 5, 8, 10, 15, 17, 19, 20, 22, 23, 24, 28, 29, 35, 36, 39, 40, 44, 45, 47, 50, 54, 56, 65, 66, 68, 69, 70, 76, 77},
		[]int{16, 62, 75}, s))
	result.PS.MaxValue = 10

	result.DAN.Value = DAN(getValue([]int{3, 4, 5, 6, 7, 8, 10, 12, 13, 14, 15, 17, 18, 19, 20, 22, 23, 24, 27, 28, 29, 31, 32, 33, 34, 35, 36, 37, 39, 40, 41, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 63, 64, 65, 66, 68, 69, 70, 71, 72, 73, 74, 76, 77},
		[]int{1, 2, 9, 11, 16, 21, 25, 26, 30, 38, 42, 62, 67, 75}, s))
	result.DAN.MaxValue = 10

	result.Description = getMLODescription(result)

	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	return resultJSON
}
