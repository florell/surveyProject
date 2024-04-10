package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	types "psychward/src"

	_ "github.com/go-sql-driver/mysql"
)

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func pushTest(db *sql.DB) {
	files, err := FilePathWalkDir("json_tests")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(files)

	for _, file := range files {
		jsonFile, err := os.Open(file)
		if err != nil {
			log.Fatalln(err)
		}
		defer jsonFile.Close()
		byteValue, _ := io.ReadAll(jsonFile)

		var survey types.Survey

		if err := json.Unmarshal([]byte(byteValue), &survey); err != nil {
			panic(err)
		}

		insertSurveyQuery := "INSERT INTO surveys (id, title, description) VALUES (?, ?, ?)"
		res, err := db.Exec(insertSurveyQuery, survey.SurveyID, survey.Title, survey.Description)
		if err != nil {
			log.Fatalf("Error inserting survey: %v", err)
		}

		surveyID, _ := res.LastInsertId()
		for _, question := range survey.Questions {
			// Insert question data into the database
			insertQuestionQuery := "INSERT INTO questions (surveyid, title) VALUES (?, ?)"
			res, err := db.Exec(insertQuestionQuery, surveyID, question.Title)
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

	}

	fmt.Println("Survey data inserted into the database successfully.")
}
