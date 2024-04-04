package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	types "psychward/src"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var db *sql.DB
var store = sessions.NewCookieStore([]byte("merogrek"))

func indexHandler(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Extract name and surname from form
		name := r.FormValue("name")
		surname := r.FormValue("surname")

		// Prepare SQL statement
		stmt, err := db.Prepare("INSERT INTO patients (name, surname) VALUES (?, ?)")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		// Execute the SQL statement
		res, err := stmt.Exec(name, surname)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		insertedID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Values["patientID"] = insertedID
		session.Save(r, w)

		// Redirect after successful form submission
		http.Redirect(w, r, "/choose", http.StatusSeeOther)
		return
	}

	// For GET request, serve the HTML form
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
	var surveys []types.Survey
	for rows.Next() {
		var survey types.Survey
		err := rows.Scan(&survey.SurveyID, &survey.Title)
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
	var survey types.Survey

	// Fetch survey data from the database
	row := db.QueryRow("SELECT id, title FROM surveys WHERE id = ?", id)
	err := row.Scan(&survey.SurveyID, &survey.Title)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		log.Fatal(err)
	}

	// Fetch questions and answers for the survey from the database
	rows, err := db.Query("SELECT id, title FROM questions WHERE surveyid = ?", survey.SurveyID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var question types.Question
		err := rows.Scan(&question.QuestionID, &question.Title)
		if err != nil {
			log.Fatal(err)
		}

		// Fetch answers for each question
		answerRows, err := db.Query("SELECT id, text, value FROM answers WHERE questionid = ?", question.QuestionID)
		if err != nil {
			log.Fatal(err)
		}
		defer answerRows.Close()

		for answerRows.Next() {
			var answer types.Answer
			err := answerRows.Scan(&answer.AnswerID, &answer.Text, &answer.Value)
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

func submitSurveyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve the session
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the patientID exists in the session
	patientID, ok := session.Values["patientID"].(int64)
	if !ok {
		http.Error(w, "Patient ID not found in session", http.StatusInternalServerError)
		return
	}

	// Parse form data
	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract survey answers from form
	surveyID := r.FormValue("survey_id")
	answers := r.Form["answers"]

	// Convert answers to JSON format
	answersJSON, err := json.Marshal(answers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare SQL statement
	stmt, err := db.Prepare("INSERT INTO survey_results (PatientID, SurveyID, Result) VALUES (?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(patientID, surveyID, string(answersJSON))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect after successful form submission
	http.Redirect(w, r, "/thankyou", http.StatusSeeOther)
}

func thankyouHandler(w http.ResponseWriter, r *http.Request) {
	// Render a thank you message
	fmt.Fprintln(w, "Thank you for submitting the survey!")
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
	r.HandleFunc("/submit_survey", submitSurveyHandler)
	r.HandleFunc("/thankyou", thankyouHandler)
	http.Handle("/", r)

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
