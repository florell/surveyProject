package handlers

import (
	"encoding/json"
	"log"
	types "psychward/src"
)

type responseField struct {
	MaxValue int     `json:"max_value"`
	Value    int     `json:"value"`
	Percent  float64 `json:"percent"`
}

type responseResult struct {
	PhysHealth    responseField `json:"Физическое здоровье"`
	PhyslogHealth responseField `json:"Физиологическое здоровье"`
	SocialHealth  responseField `json:"Созиальные взаимоотношения"`
	Environment   responseField `json:"Окружающая среда"`
}

func VOZHandler(s *types.SurveyResults) []byte {
	result := responseResult{}
	
	result.PhysHealth.MaxValue = 35
	result.PhysHealth.Value = (6 - s.Picked[4+90+50]) + (6 - s.Picked[5+90+50]) + s.Picked[11+90+50] + s.Picked[16+90+50] + s.Picked[17+90+50] +
		+s.Picked[18+90+50] + s.Picked[19+90+50]
	result.PhysHealth.Percent = ((float64(result.PhysHealth.Value) - 7) / 28) * 100
	
	result.PhysHealth.MaxValue = 30
	result.PhyslogHealth.Value = s.Picked[6+90+50] + s.Picked[7+90+50] + s.Picked[8+90+50] + s.Picked[12+90+50] +
		+s.Picked[20+90+50] + (6 - s.Picked[27+90+50])
	result.PhyslogHealth.Percent = ((float64(result.PhyslogHealth.Value) - 6) / 24) * 100
	
	result.SocialHealth.MaxValue = 15
	result.SocialHealth.Value = s.Picked[21] + s.Picked[22+90+50] + s.Picked[23+90+50]
	result.SocialHealth.Percent = ((float64(result.SocialHealth.Value) - 3) / 12) * 100
	
	result.Environment.MaxValue = 40
	result.Environment.Value = s.Picked[9+90+50] + s.Picked[10+90+50] + s.Picked[13+90+50] +
		+s.Picked[14+90+50] + s.Picked[15+90+50] + s.Picked[24+90+50] + s.Picked[25+90+50] + s.Picked[26+90+50]
	result.Environment.Percent = ((float64(result.Environment.Value) - 8) / 32) * 100
	
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
	}
	return resultJSON
}
