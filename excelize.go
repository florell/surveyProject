package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	_ "github.com/go-sql-driver/mysql"
)

type Patient struct {
	Name    string
	Surname string
	Age     string
	Sex     string
}

func getSurveyName(db *sql.DB, surveyID int) (string, error) {
	var surveyName string
	err := db.QueryRow("SELECT Title FROM surveys WHERE id = ?", surveyID).Scan(&surveyName)
	if err != nil {
		return "", err
	}
	return surveyName, nil
}
func getPatient(db *sql.DB, patientID int) (Patient, error) {
	p := Patient{}
	err := db.QueryRow("SELECT Name, Surname, Age, Sex FROM patients WHERE id = ?", patientID).Scan(&p.Name, &p.Surname, &p.Age, &p.Sex)
	if err != nil {
		return Patient{}, err
	}
	return p, nil
}

func makeTable(db *sql.DB) error {
	rows, err := db.Query("SELECT ID, PatientID, SurveyID, Result, CurDate FROM survey_results")
	if err != nil {
		return err
	}
	defer rows.Close()

	xlsx := excelize.NewFile()
	sheetName := "Результаты опроса"
	xlsx.SetSheetName("Sheet1", sheetName)
	patientLastRow := make(map[int]int)
	surveyColumns := make(map[string]int)
	startRow := 2
	h := 4
	styles := []int{}

	style1, _ := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#87CEFA"],"pattern":1}}`)
	style2, _ := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#FFB6C1"],"pattern":1}}`)
	style3, _ := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#98FB98"],"pattern":1}}`)
	style4, _ := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#FFD700"],"pattern":1}}`)
	style5, _ := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#FFA07A"],"pattern":1}}`)
	style6, _ := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#FFE4C4"],"pattern":1}}`)
	style7, _ := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#98FB98"],"pattern":1}}`)
	style8, _ := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#DDA0DD"],"pattern":1}}`)
	style9, _ := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#FFA07A"],"pattern":1}}`)
	style10, _ := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#FFA07A"],"pattern":1}}`)

	styles = append(styles, style1)
	styles = append(styles, style2)
	styles = append(styles, style3)
	styles = append(styles, style4)
	styles = append(styles, style5)
	styles = append(styles, style6)
	styles = append(styles, style7)
	styles = append(styles, style8)
	styles = append(styles, style9)
	styles = append(styles, style10)

	xlsx.SetCellValue(sheetName, "A1", "Пациенты")
	xlsx.MergeCell(sheetName, "A1", "D1")
	xlsx.SetCellValue(sheetName, "A2", "Имя")
	xlsx.SetCellValue(sheetName, "B2", "Фамилия")
	xlsx.SetCellValue(sheetName, "C2", "Возраст")
	xlsx.SetCellValue(sheetName, "D2", "Пол")
	xlsx.SetCellStyle("Sheet1", "A1", "D2", styles[0])

	for rows.Next() {
		s := struct {
			ID        int
			patientID int
			surveyID  int
			resultStr string
			curDate   string
		}{}

		err := rows.Scan(&s.ID, &s.patientID, &s.surveyID, &s.resultStr, &s.curDate)
		if err != nil {
			return err
		}

		surveyName, err := getSurveyName(db, s.surveyID)
		if err != nil {
			return err
		}

		patient, err := getPatient(db, s.patientID)
		if err != nil {
			return err
		}

		var result map[string]interface{}
		err = json.Unmarshal([]byte(s.resultStr), &result)
		if err != nil {
			return err
		}

		lastRow, exists := patientLastRow[s.patientID]
		if !exists {
			lastRow = len(patientLastRow) + 1 + startRow
			patientLastRow[s.patientID] = lastRow
			ruSex := ""
			if patient.Sex == "male" {
				ruSex = "мужской"
			} else {
				ruSex = "женский"
			}

			xlsx.SetCellValue(sheetName, "A"+strconv.Itoa(lastRow), patient.Name)
			xlsx.SetCellValue(sheetName, "B"+strconv.Itoa(lastRow), patient.Surname)
			xlsx.SetCellValue(sheetName, "C"+strconv.Itoa(lastRow), patient.Age)
			xlsx.SetCellValue(sheetName, "D"+strconv.Itoa(lastRow), ruSex)
		}

		columnIndex, exists := surveyColumns[surveyName]
		if !exists {

			columnIndex = h
			surveyColumns[surveyName] = columnIndex

			var t int
			startCell := excelize.ToAlphaString(columnIndex) + "1"
			if s.surveyID != 6 {
				t = len(result) - 1
			} else {
				t = len(result)*2 - 3
			}
			endCell := excelize.ToAlphaString(columnIndex+t) + strconv.Itoa(1)
			xlsx.SetCellStyle("Sheet1", startCell, excelize.ToAlphaString(columnIndex+t)+strconv.Itoa(2), styles[s.surveyID])

			xlsx.SetCellValue(sheetName, startCell, surveyName)
			xlsx.MergeCell(sheetName, startCell, endCell)
			h += len(result)
		}
		sres := []string{}
		for o := range result {
			sres = append(sres, o)
		}
		sort.Strings(sres)
		colInd := columnIndex
		for _, value := range sres {
			switch v := result[value].(type) {
			case map[string]interface{}:
				sk := []string{}
				for g := range v {
					sk = append(sk, g)
				}
				sort.Strings(sk)
				fmt.Println(sk)
				for _, k := range sk {
					if s, ok := v[k].(float64); ok && k == "value" {
						xlsx.SetCellValue(sheetName, excelize.ToAlphaString(colInd)+strconv.Itoa(2), value)
						xlsx.SetCellValue(sheetName, excelize.ToAlphaString(colInd)+strconv.Itoa(lastRow), int(s))
						colInd++
					} else if s, ok := v[k].(float64); ok && k == "tscore" {
						xlsx.SetCellValue(sheetName, excelize.ToAlphaString(colInd)+strconv.Itoa(2), value)
						xlsx.SetCellValue(sheetName, excelize.ToAlphaString(colInd)+strconv.Itoa(lastRow), fmt.Sprintf("%d", int(s))+"T")
						colInd++
					}
				}
			}
		}
		if s.surveyID == 6 {
			for i := columnIndex; i < colInd-1; i += 2 {
				xlsx.MergeCell(sheetName, excelize.ToAlphaString(i)+strconv.Itoa(2), excelize.ToAlphaString(i+1)+strconv.Itoa(2))
			}
		}
	}

	if err := xlsx.SaveAs("Survey_Results.xlsx"); err != nil {
		return err
	}

	return err

}
