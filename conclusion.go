package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/fumiama/go-docx"
	_ "github.com/go-sql-driver/mysql"
)

func declineAge(sAge string) string {
	yearForms := []string{"год", "года", "лет"}
	age, _ := strconv.Atoi(sAge)

	var form string
	if age%10 == 1 && age%100 != 11 {
		form = yearForms[0]
	} else if (age%10 >= 2 && age%10 <= 4) && (age%100 < 10 || age%100 >= 20) {
		form = yearForms[1]
	} else {
		form = yearForms[2]
	}

	return fmt.Sprintf("%d %s", age, form)
}

func makeConclusion(db *sql.DB, patientID string) error {
	// Проверка соединения
	err := db.Ping()
	if err != nil {
		return err
	}

	// Выполнение SQL-запроса
	rows, err := db.Query("SELECT SurveyID, Result, CurDate, Description FROM survey_results WHERE PatientID = ? ORDER BY CurDate ASC", patientID)
	if err != nil {
		return err
	}
	defer rows.Close()

	intPId, _ := strconv.Atoi(patientID)
	patient, err := getPatient(db, intPId)

	if err != nil {
		return err
	}

	doc := docx.New().WithDefaultTheme()

	lastDate := ""

	// Добавление заголовка
	paraTitle := doc.AddParagraph()
	paraTitle.AddText("Заключение по результатам экспериментально-психологического обследования").Bold()
	paraFIO := doc.AddParagraph()
	paraFIO.AddText("ФИ: ").Bold()
	paraFIO.AddText(patient.Name + " " + patient.Surname)
	paraAge := doc.AddParagraph()
	paraAge.AddText("Возраст: ").Bold()
	paraAge.AddText(declineAge(patient.Age))

	paraDates := doc.AddParagraph()
	paraDates.AddText("Даты обследования: ").Bold()

	paraR := doc.AddParagraph()
	paraR.AddText("Обследование когнитивной сферы").Bold()

	for rows.Next() {
		s := struct {
			SurveyID    int
			Result      string
			CurDate     string
			Description sql.NullString
		}{}

		err = rows.Scan(&s.SurveyID, &s.Result, &s.CurDate, &s.Description)
		if err != nil {
			log.Fatalln(err)
		}
		if s.CurDate != lastDate {
			paraDates.AddText(s.CurDate + "; ")
			lastDate = s.CurDate
		}
		paraTemp := doc.AddParagraph()
		surveyName, _ := getSurveyName(db, s.SurveyID)
		paraTemp.AddText(surveyName + ": ").Bold()

		var data map[string]interface{}
		err = json.Unmarshal([]byte(s.Result), &data)
		if err != nil {
			return err
		}

		sdata := []string{}
		for o := range data {
			sdata = append(sdata, o)
		}
		sort.Strings(sdata)
		for _, value := range sdata {
			switch v := data[value].(type) {
			case map[string]interface{}:
				sk := []string{}
				for g := range v {
					sk = append(sk, g)
				}
				sort.Strings(sk)
				fr := make([]string, 4, 10)
				fr[0] = value + " "
				for _, k := range sk {
					if s, ok := v[k].(float64); ok && k == "value" {
						ss := strconv.FormatFloat(s, 'f', 0, 64)
						fr[1] = ss + " / "
					} else if s, ok := v[k].(float64); ok && k == "max_value" {
						ss := strconv.FormatFloat(s, 'f', 0, 64)
						fr[2] = ss
					} else if s, ok := v[k].(string); ok && k == "description" {
						if s == "" {
							continue
						}
						fr[3] = " - " + s
					}
				}
				rrr := strings.TrimRight(strings.Join(fr, ""), " ") + "; "
				paraTemp.AddText(rrr)
			}
		}

	}

	// Проверка на ошибки
	if err := rows.Err(); err != nil {
		return err
	}

	// Сохранение документа
	f, err := os.Create("conclusion_" + patientID + ".docx")
	if err != nil {
		panic(err)
	}
	_, err = doc.WriteTo(f)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}

	return nil
}
