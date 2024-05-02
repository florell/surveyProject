package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"psychward/handlers"
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
		log.Println("Error getting session:", err)
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
		defer func(stmt *sql.Stmt) {
			err := stmt.Close()
			if err != nil {
				log.Println(err)
			}
		}(stmt)
		
		// Execute the SQL statement
		res, err := stmt.Exec(name, surname, sex, age)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		insertedID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalln(err)
		}
		
		log.Printf("%T %v %T %v\n", sex, sex, age, age)
		
		session.Values["patientID"] = insertedID
		session.Values["patientGender"] = sex
		session.Values["patientAge"] = age
		
		if err := session.Save(r, w); err != nil {
			log.Println("Error saving cookies:", err)
			return
		}
		
		// Redirect after successful form submission
		http.Redirect(w, r, "/choose", http.StatusSeeOther)
		return
	}
	
	// For GET request, serve the HTML form
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		log.Println(err)
		return
	}
}

func chooseHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch surveys from the database
	rows, err := db.Query("SELECT id, title, description FROM surveys")
	if err != nil {
		log.Fatal(err)
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	}(rows)
	
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
	
	if err := tmpl.Execute(w, surveys); err != nil {
		log.Println(err)
		return
	}
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
		if errors.Is(err, sql.ErrNoRows) {
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
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	}(rows)
	
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
		defer func(answerRows *sql.Rows) {
			if err := answerRows.Close(); err != nil {
				log.Println(err)
			}
		}(answerRows)
		
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
	
	tmpl := template.Must(template.ParseFiles("templates/survey.html"))
	if err != nil {
		log.Println("Error parsing template:", err)
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
		log.Println("Error getting session:", err)
		return
	}
	
	// Check if the patientID exists in the session
	patientID, ok := session.Values["patientID"].(int64)
	if !ok {
		http.Error(w, "Patient ID not found in session", http.StatusInternalServerError)
		return
	}
	
	patientGender, ok := session.Values["patientGender"].(string)
	if !ok {
		http.Error(w, "Patient gender not found in this session", http.StatusInternalServerError)
	}
	
	patientAge, ok := session.Values["patientAge"].(string)
	if !ok {
		http.Error(w, "Patient age not found in this session", http.StatusInternalServerError)
	}
	patientAgeInt, err := strconv.Atoi(patientAge)
	if err != nil {
		log.Fatalln(err)
	}
	
	// Parse form data
	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Extract survey answers from form
	surveyID, err := strconv.Atoi(r.FormValue("survey_id"))
	if err != nil {
		http.Error(w, "Invalid survey ID", http.StatusBadRequest)
		return
	}
	
	selectedAnswers := make(map[int]int)
	for key, values := range r.Form {
		if strings.HasPrefix(key, "question") {
			// Extract question number and answer ID
			questionID := strings.Split(key, "question")[1]
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
	
	surveyResults := types.SurveyResults{
		SurveyID:  surveyID,
		PatientID: int(patientID),
		Age:       patientAgeInt,
		Sex:       patientGender,
		Picked:    selectedAnswers,
	}
	
	var analysis []byte
	switch surveyID {
	case 1:
		analysis = handlers.FESHandler(&surveyResults)
	case 2:
		analysis = handlers.WCQHandler(&surveyResults)
	case 3:
		analysis = handlers.VOZHandler(&surveyResults)
	case 4:
		analysis = handlers.BDIHandler(&surveyResults)
	case 5:
		analysis = handlers.ITTHandler(&surveyResults)
	case 6:
		analysis = handlers.MLOHandler(&surveyResults)
	default:
		log.Println("Survey ID is not supported:", surveyID)
		http.Error(w, fmt.Sprintf("Survey ID is not supported: %s", surveyID), http.StatusBadRequest)
		return
	}
	err = session.Save(r, w)
	if err != nil {
		log.Fatalln("Error saving results in session:", err)
	}
	
	// Prepare SQL statement
	stmt, err := db.Prepare("INSERT INTO survey_results (PatientID, SurveyID, Result) VALUES (?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Println(err)
		}
	}(stmt)
	
	// Execute the SQL statement
	_, err = stmt.Exec(patientID, surveyID, string(analysis))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	var resultID string
	if err := db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&resultID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Redirect after successful form submission
	http.Redirect(w, r, "/result?survey_id="+strconv.Itoa(surveyID)+"&result_id="+resultID, http.StatusSeeOther)
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	resultId := r.URL.Query().Get("result_id")
	
	var resD string
	if err := db.QueryRow("SELECT Result FROM survey_results WHERE ID = ?", resultId).Scan(&resD); err != nil {
		http.Error(w, fmt.Sprintf("Error getting results from database: %s", err.Error()), http.StatusInternalServerError)
	}
	
	tmpl := template.Must(template.ParseFiles("templates/results.html"))
	
	fmt.Println("___", resD, "___")
	
	if err := tmpl.Execute(w, resD); err != nil {
		log.Fatal(err)
	}
	
}

func main() {
	var err error
	db, err = sql.Open("mysql", "psy_admin:pw2319#@tcp(127.0.0.1:3306)/psy_data")
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Println(err)
		}
	}(db)
	
	pushFlag := flag.Bool("push", false, "Use this flag to push")
	flag.Parse()
	if *pushFlag {
		pushTest(db)
		return
	}
	
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/choose", chooseHandler)
	r.HandleFunc("/survey/{id}", surveyHandler)
	r.HandleFunc("/submit_survey", submitSurveyHandler)
	r.HandleFunc("/result", resultHandler)
	http.Handle("/", r)
	
	fmt.Println("Server is running on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
