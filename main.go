package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func chooseHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch surveys from the database
	rows, err := db.Query("SELECT id, title FROM surveys")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	
	var surveys []Survey
	for rows.Next() {
		var survey Survey
		err := rows.Scan(&survey.ID, &survey.Title)
		if err != nil {
			log.Fatal(err)
		}
		surveys = append(surveys, survey)
	}
	
	tmpl := template.Must(template.ParseFiles("templates/choose.html"))
	tmpl.Execute(w, surveys)
}

func surveyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	// Fetch survey data from the database using survey ID
	var survey Survey
	
	// Fetch survey data from the database
	row := db.QueryRow("SELECT id, title FROM surveys WHERE id = ?", id)
	err := row.Scan(&survey.ID, &survey.Title)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		log.Fatal(err)
	}
	
	// Fetch questions and answers for the survey from the database
	rows, err := db.Query("SELECT id, title FROM questions WHERE surveyid = ?", survey.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var question Question
		err := rows.Scan(&question.ID, &question.Title)
		if err != nil {
			log.Fatal(err)
		}
		
		// Fetch answers for each question
		answerRows, err := db.Query("SELECT id, text, value FROM answers WHERE questionid = ?", question.ID)
		if err != nil {
			log.Fatal(err)
		}
		defer answerRows.Close()
		
		for answerRows.Next() {
			var answer Answer
			err := answerRows.Scan(&answer.ID, &answer.Text, &answer.Value)
			if err != nil {
				log.Fatal(err)
			}
			question.Answers = append(question.Answers, answer)
		}
		
		survey.Questions = append(survey.Questions, question)
	}
	
	// Temporary solution
	// data, err := json.Marshal(survey)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// w.Header().Set("Content-Type", "application/json")
	// w.Write(data)
	
	tmpl, err := template.ParseFiles("templates/survey.html")
	if err != nil {
		log.Fatal(err)
	}
	
	// Execute the template with survey data
	err = tmpl.Execute(w, survey)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var err error
	db, err = sql.Open("mysql", "psy_admin:pw2319#@tcp(127.0.0.1:3306)/psy_data")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// pushTest(db)
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/choose", chooseHandler)
	r.HandleFunc("/survey/{id}", surveyHandler)
	
	http.Handle("/", r)
	
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
