package handlers

import (
	"encoding/json"
	"log"
	types "psychward/src"
)

type VECResult struct {
	AP          int    `json:"Абсолютный показатель (сумма результатов по субтестам)"`
	KP          int    `json:"Корригированный на возраст показатель"`
	IPP         int    `json:"Эквивалентный интеллекту показатель памяти"`
	Description string `json:"description"`
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
	switch KP {
	case 50:
		IPP = 48
	case 51:
		IPP = 49
	case 52:
		IPP = 49
	case 53:
		IPP = 50
	case 54:
		IPP = 51
	case 55:
		IPP = 52
	case 56:
		IPP = 52
	case 57:
		IPP = 53
	case 58:
		IPP = 54
	case 59:
		IPP = 55
	case 60:
		IPP = 55
	case 61:
		IPP = 56
	case 62:
		IPP = 57
	case 63:
		IPP = 57
	case 64:
		IPP = 58
	case 65:
		IPP = 59
	case 66:
		IPP = 59
	case 67:
		IPP = 60
	case 68:
		IPP = 61
	case 69:
		IPP = 62
	case 70:
		IPP = 62
	case 71:
		IPP = 63
	case 72:
		IPP = 64
	case 73:
		IPP = 64
	case 74:
		IPP = 66
	case 75:
		IPP = 67
	case 76:
		IPP = 69
	case 77:
		IPP = 70
	case 78:
		IPP = 72
	case 79:
		IPP = 73
	case 80:
		IPP = 74
	case 81:
		IPP = 76
	case 82:
		IPP = 77
	case 83:
		IPP = 79
	case 84:
		IPP = 80
	case 85:
		IPP = 81
	case 86:
		IPP = 84
	case 87:
		IPP = 84
	case 88:
		IPP = 86
	case 89:
		IPP = 87
	case 90:
		IPP = 89
	case 91:
		IPP = 90
	case 92:
		IPP = 92
	case 93:
		IPP = 93
	case 94:
		IPP = 94
	case 95:
		IPP = 96
	case 96:
		IPP = 97
	case 97:
		IPP = 99
	case 98:
		IPP = 100
	case 99:
		IPP = 101
	case 100:
		IPP = 103
	case 101:
		IPP = 105
	case 102:
		IPP = 106
	case 103:
		IPP = 108
	case 104:
		IPP = 110
	case 105:
		IPP = 112
	case 106:
		IPP = 114
	case 107:
		IPP = 116
	case 108:
		IPP = 118
	case 109:
		IPP = 120
	case 110:
		IPP = 122
	case 111:
		IPP = 124
	case 112:
		IPP = 128
	case 113:
		IPP = 129
	case 114:
		IPP = 132
	case 115:
		IPP = 135
	case 116:
		IPP = 137
	case 117:
		IPP = 140
	case 118:
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
	
	result := VECResult{AP, KP, IPP, desc}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	return resultJSON
}
