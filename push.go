package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func pushTest(db *sql.DB) {
	// Sample JSON data representing a survey
	jsonData := `
	{
		"ID": 2,
		"Title": "Survey 3",
		"Questions": [
			{
				"Title": "Question 1_3",
				"Answers": [
					{"Text": "Option 1", "Value": "value1"},
					{"Text": "Option 2", "Value": "value2"},
					{"Text": "Option 3", "Value": "value3"}
				]
			},
			{
				"Title": "Question 2_3",
				"Answers": [
					{"Text": "Option A", "Value": "valueA"},
					{"Text": "Option B", "Value": "valueB"},
					{"Text": "Option C", "Value": "valueC"}
				]
			}
		]
	}`

	// Parse JSON data into Survey struct
	var survey Survey
	err := json.Unmarshal([]byte(jsonData), &survey)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Insert survey data into the database
	insertSurveyQuery := "INSERT INTO surveys (id, title) VALUES (?, ?)"
	_, err = db.Exec(insertSurveyQuery, survey.ID, survey.Title)
	if err != nil {
		log.Fatalf("Error inserting survey: %v", err)
	}

	for _, question := range survey.Questions {
		// Insert question data into the database
		insertQuestionQuery := "INSERT INTO questions (surveyid, title) VALUES (?, ?)"
		res, err := db.Exec(insertQuestionQuery, survey.ID, question.Title)
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
