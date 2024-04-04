package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	types "psychward/src"

	_ "github.com/go-sql-driver/mysql"
)

func pushTest(db *sql.DB) {
	// Sample JSON data representing a survey
	jsonData := `
	{
		"ID": 2,
		"Title": "Survey AQ",
		"Questions": [
			{
				"Title": "Question first",
				"Answers": [
					{"Text": "Option 1989", "Value": 80},
					{"Text": "Option 2239", "Value": 13},
					{"Text": "Option 38098098", "Value": 391}
				]
			},
			{
				"Title": "Question second",
				"Answers": [
					{"Text": "Ojf", "Value": 14},
					{"Text": "Ophu dfg", "Value": 72},
					{"Text": "Optiijo Cfhdhg", "Value": 49}
				]
			},
			{
				"Title": "Question third",
				"Answers": [
					{"Text": "udas", "Value": 12},
					{"Text": "opif dfg", "Value": 90},
					{"Text": "wneiweruigor", "Value": 88}
				]
			}
		]
	}`

	// Parse JSON data into Survey struct
	var survey types.Survey
	err := json.Unmarshal([]byte(jsonData), &survey)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Insert survey data into the database
	insertSurveyQuery := "INSERT INTO surveys (id, title) VALUES (?, ?)"
	_, err = db.Exec(insertSurveyQuery, survey.SurveyID, survey.Title)
	if err != nil {
		log.Fatalf("Error inserting survey: %v", err)
	}

	for _, question := range survey.Questions {
		// Insert question data into the database
		insertQuestionQuery := "INSERT INTO questions (surveyid, title) VALUES (?, ?)"
		res, err := db.Exec(insertQuestionQuery, survey.SurveyID, question.Title)
		if err != nil {
			log.Fatalf("Error inserting question: %v", err)
		}
		questionID, err := res.LastInsertId()
		if err != nil {
			log.Fatalf("Error getting last insert ID for question: %v", err)
		}

		// Insert answer data into the database
		for _, answer := range question.Answers {
			insertAnswerQuery := "INSERT INTO answers (questionid, text, value) VALUES (?, ?, ?)"
			_, err := db.Exec(insertAnswerQuery, questionID, answer.Text, answer.Value)
			if err != nil {
				log.Fatalf("Error inserting answer: %v", err)
			}
		}
	}

	fmt.Println("Survey data inserted into the database successfully.")
}
