package handlers

import (
	"encoding/json"
	"log"
	types "psychward/src"
)

type VECField struct {
	Value    int `json:"value"`
	MaxValue int `json:"max_value"`
}

type VECResult struct {
	AP          VECField `json:"Абсолютный показатель (сумма результатов по субтестам)"`
	KP          VECField `json:"Корригированный на возраст показатель"`
	IPP         VECField `json:"Эквивалентный интеллекту показатель памяти"`
	Description string   `json:"description"`
}

func getKP(AP int) int {
	switch {
	case AP <= 19:
		AP += 32
	case AP >= 20 && AP <= 24:
		AP += 33
	case AP >= 25 && AP <= 29:
		AP += 34
	case AP >= 30 && AP <= 34:
		AP += 36
	case AP >= 35 && AP <= 39:
		AP += 38
	case AP >= 40 && AP <= 44:
		AP += 40
	case AP >= 45 && AP <= 49:
		AP += 42
	case AP >= 50 && AP <= 54:
		AP += 44
	case AP >= 55 && AP <= 59:
		AP += 46
	case AP >= 60 && AP <= 64:
		AP += 48
	case AP >= 65 && AP <= 69:
		AP += 50
	case AP >= 70 && AP <= 74:
		AP += 52
	case AP >= 75 && AP <= 79:
		AP += 54
	case AP >= 80:
		AP += 56
	}
	
	return AP
}

func getIPP(KP int) int {
	IPP := 0
	switch {
	case KP <= 50:
		IPP = 48
	case KP == 51:
		IPP = 49
	case KP == 52:
		IPP = 49
	case KP == 53:
		IPP = 50
	case KP == 54:
		IPP = 51
	case KP == 55:
		IPP = 52
	case KP == 56:
		IPP = 52
	case KP == 57:
		IPP = 53
	case KP == 58:
		IPP = 54
	case KP == 59:
		IPP = 55
	case KP == 60:
		IPP = 55
	case KP == 61:
		IPP = 56
	case KP == 62:
		IPP = 57
	case KP == 63:
		IPP = 57
	case KP == 64:
		IPP = 58
	case KP == 65:
		IPP = 59
	case KP == 66:
		IPP = 59
	case KP == 67:
		IPP = 60
	case KP == 68:
		IPP = 61
	case KP == 69:
		IPP = 62
	case KP == 70:
		IPP = 62
	case KP == 71:
		IPP = 63
	case KP == 72:
		IPP = 64
	case KP == 73:
		IPP = 64
	case KP == 74:
		IPP = 66
	case KP == 75:
		IPP = 67
	case KP == 76:
		IPP = 69
	case KP == 77:
		IPP = 70
	case KP == 78:
		IPP = 72
	case KP == 79:
		IPP = 73
	case KP == 80:
		IPP = 74
	case KP == 81:
		IPP = 76
	case KP == 82:
		IPP = 77
	case KP == 83:
		IPP = 79
	case KP == 84:
		IPP = 80
	case KP == 85:
		IPP = 81
	case KP == 86:
		IPP = 84
	case KP == 87:
		IPP = 84
	case KP == 88:
		IPP = 86
	case KP == 89:
		IPP = 87
	case KP == 90:
		IPP = 89
	case KP == 91:
		IPP = 90
	case KP == 92:
		IPP = 92
	case KP == 93:
		IPP = 93
	case KP == 94:
		IPP = 94
	case KP == 95:
		IPP = 96
	case KP == 96:
		IPP = 97
	case KP == 97:
		IPP = 99
	case KP == 98:
		IPP = 100
	case KP == 99:
		IPP = 101
	case KP == 100:
		IPP = 103
	case KP == 101:
		IPP = 105
	case KP == 102:
		IPP = 106
	case KP == 103:
		IPP = 108
	case KP == 104:
		IPP = 110
	case KP == 105:
		IPP = 112
	case KP == 106:
		IPP = 114
	case KP == 107:
		IPP = 116
	case KP == 108:
		IPP = 118
	case KP == 109:
		IPP = 120
	case KP == 110:
		IPP = 122
	case KP == 111:
		IPP = 124
	case KP == 112:
		IPP = 128
	case KP == 113:
		IPP = 129
	case KP == 114:
		IPP = 132
	case KP == 115:
		IPP = 135
	case KP == 116:
		IPP = 137
	case KP == 117:
		IPP = 140
	case KP >= 118:
		IPP = 143
	}
	
	if KP > 118 {
		IPP = 143
	}
	
	return IPP
}

func VECHandler(s *types.SurveyResults) []byte {
	AP := 0
	for _, val := range s.Picked {
		AP += val
	}
	
	KP := getKP(AP)
	IPP := getIPP(KP)
	
	desc := ""
	switch {
	case IPP >= 110:
		desc = "нормальное функционирование памяти"
	case IPP >= 93 && IPP <= 106:
		desc = "нарушения памяти легкой степени"
	case IPP >= 73 && IPP <= 87:
		desc = "нарушения памяти умеренной степени"
	case IPP >= 48 && IPP <= 66:
		desc = "нарушения памяти выраженной степени"
	}
	
	switch {
	case IPP >= 130:
		desc += "/ очень высокий"
	case IPP >= 120 && IPP <= 129:
		desc += "/ высокий"
	case IPP >= 110 && IPP <= 119:
		desc += "/ хорошая норма"
	case IPP >= 90 && IPP <= 109:
		desc += "/ средний"
	case IPP >= 80 && IPP <= 89:
		desc += "/ низкая (плохая) норма"
	case IPP >= 70 && IPP <= 79:
		desc += "/ пограничная зона"
	case IPP <= 69:
		desc += "/ умственный дефект"
	}
	
	switch {
	case IPP >= 68 && IPP <= 80:
		desc += "/ пограничная УО"
	case IPP >= 52 && IPP <= 57:
		desc += "/ легкая УО (дебильность)"
	case IPP >= 36 && IPP <= 51:
		desc += "/ умеренная (средняя) УО, невыраженная имбецильность"
	case IPP >= 20 && IPP <= 35:
		desc += "/ глубокая УО, выраженная имбецильность"
	case IPP < 20:
		desc += "/ полная уо (идиотия)"
	}
	
	result := VECResult{VECField{
		Value:    AP,
		MaxValue: 93,
	}, VECField{
		Value:    KP,
		MaxValue: 93 + 56,
	}, VECField{
		Value:    IPP,
		MaxValue: 143,
	}, desc}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	return resultJSON
}
