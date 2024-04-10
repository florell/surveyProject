package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	handlers "psychward/handlers"
	types "psychward/src"
	"strconv"
	"strings"

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
		sex := r.FormValue("sex")
		age := r.FormValue("age")

		// Prepare SQL statement
		stmt, err := db.Prepare("INSERT IGNORE patients (name, surname, sex, age) VALUES (?, ?, ?, ?)")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		// Execute the SQL statement
		res, err := stmt.Exec(name, surname, sex, age)
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
		session.Values["patientGender"] = sex
		session.Values["patientAge"] = age
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
	rows, err := db.Query("SELECT id, title, description FROM surveys")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var surveys []types.Survey
	for rows.Next() {
		var survey types.Survey
		err := rows.Scan(&survey.SurveyID, &survey.Title, &survey.Description)
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

	tmpl := template.Must(template.ParseFiles("templates/survey.html"))
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
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
	surveyIDstr := r.FormValue("survey_id")
	surveyID, err := strconv.Atoi(surveyIDstr)
	if err != nil {
		http.Error(w, "Invalid survey ID", http.StatusBadRequest)
		return
	}

	selectedAnswers := make(map[int]int)
	for key, values := range r.Form {
		fmt.Println(key, values)
		if strings.HasPrefix(key, "question") {
			// Extract question number and answer ID
			questionID := strings.Split(key, "question")[1]
			fmt.Println(questionID)
			for _, value := range values {
				qIDint, ok := strconv.Atoi(questionID)
				if ok != nil {
					log.Println(ok)
				}
				intValue, ok := strconv.Atoi(value)
				if ok != nil {
					log.Println(ok)
				}
				selectedAnswers[qIDint] = intValue
			}
		}
	}

	fmt.Println("Selected Answers:")
	for questionID, answerID := range selectedAnswers {
		fmt.Printf("Question %d: Answer %d\n", questionID, answerID)
	}

	surveyResults := types.SurveyResults{
		SurveyID:  surveyID,
		PatientID: int(patientID),
		Picked:    selectedAnswers,
	}
	fmt.Println(surveyResults.Picked)
	analysis := handlers.FamilyEnvironmentalScaleHandler(&surveyResults)
	fmt.Println(analysis)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
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
	_, err = stmt.Exec(patientID, surveyID, string(analysis))
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

	pushFlag := flag.Bool("push", false, "Use this flag to enable push")
	flag.Parse()
	if *pushFlag {
		pushTest(db)
		return
	}

	// pushTest(db)
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/choose", chooseHandler)
	r.HandleFunc("/survey/{id}", surveyHandler)
	r.HandleFunc("/submit_survey", submitSurveyHandler)
	r.HandleFunc("/thankyou", thankyouHandler)
	http.Handle("/", r)

	fmt.Println("Server is running on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
