package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"psychward/handlers"
	types "psychward/src"
	"strconv"
	"strings"
	"syscall"
	"time"

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
		stmt, err := db.Prepare("INSERT INTO patients (name, surname, sex, age) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE id = LAST_INSERT_ID(id)")
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

	tmpl := template.Must(template.New("survey.html").Funcs(template.FuncMap{"sum": func(a, b int) int {
		return a + b
	}}).ParseFiles("templates/survey.html"))
	log.Println(survey.Questions[0].QuestionID)

	if err != nil {
		log.Println("Error parsing template:", err)
		return
	}

	// Execute the template with survey data and question count
	err = tmpl.Execute(w, struct {
		Survey          types.Survey
		QuestionCount   int
		FirstQuestionID int
	}{
		Survey:          survey,
		QuestionCount:   len(survey.Questions),
		FirstQuestionID: survey.Questions[0].QuestionID,
	})
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
		log.Println(key, values)
		if strings.HasPrefix(key, "question") {
			// Extract question number and answer ID
			questionID := strings.Split(key, "question_")[1]
			qIDint, ok := strconv.Atoi(questionID)
			if ok != nil {
				log.Println(ok)
			}
			intValue, ok := strconv.Atoi(values[0])
			if ok != nil {
				log.Println(ok)
			}
			selectedAnswers[qIDint] = intValue
		}
	}

	surveyResults := types.SurveyResults{
		SurveyID:  surveyID,
		PatientID: int(patientID),
		Age:       patientAgeInt,
		Sex:       patientGender,
		Picked:    selectedAnswers,
	}

	fmt.Println(surveyResults)

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
		http.Error(w, fmt.Sprintf("Survey ID is not supported: %d", surveyID), http.StatusBadRequest)
		return
	}
	err = session.Save(r, w)
	if err != nil {
		log.Fatalln("Error saving results in session:", err)
	}

	// Prepare SQL statement
	stmt, err := db.Prepare("INSERT INTO survey_results (PatientID, SurveyID, CurDate, Result) VALUES (?, ?, CURDATE(), ?) ON DUPLICATE KEY UPDATE CurDate = CURDATE(), Result = ?")
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
	_, err = stmt.Exec(patientID, surveyID, string(analysis), string(analysis))
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

// func patientLogHandler(w http.ResponseWriter, r *http.Request) {
// 	session, err := store.Get(r, "session-name")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		log.Println("Error getting session:", err)
// 		return
// 	}

// 	type SurveyResult struct {
// 		SurveyID int
// 		Result   string
// 		Date     string
// 	}

// 	patientID, ok := session.Values["patientID"].(int64)
// 	if !ok {
// 		http.Error(w, "Patient ID not found in session", http.StatusInternalServerError)
// 		return
// 	}

// 	rows, err := db.Query("SELECT SurveyID, Result, CurDate FROM survey_results WHERE PatientID = ?", patientID)
// 	if err != nil {
// 		log.Println("Error querying database:", err)
// 		return
// 	}
// 	defer rows.Close()

// 	var results []SurveyResult
// 	for rows.Next() {
// 		var result SurveyResult
// 		if err := rows.Scan(&result.SurveyID, &result.Result, &result.Date); err != nil {
// 			log.Println("Error scanning row:", err)
// 			return
// 		}
// 		results = append(results, result)
// 	}

// 	if err := rows.Err(); err != nil {
// 		log.Println("Error iterating over rows:", err)
// 		return
// 	}

// 	for _, result := range results {
// 		log.Printf("ID: %d, Result: %s, Date: %s\n", result.SurveyID, result.Result, result.Date)
// 	}

// }

func generateTable(w http.ResponseWriter, r *http.Request) {
	err := makeTable(db)
	if err != nil {
		http.Error(w, "Ошибка при создании таблицы", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func downloadTable(w http.ResponseWriter, r *http.Request) {
	// Открываем файл с результатами таблицы
	file, err := os.Open("Survey_Results.xlsx")
	if err != nil {
		http.Error(w, "Ошибка при открытии файла таблицы", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Отправляем файл пользователю
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=survey_results.xlsx")
	stat, _ := file.Stat()
	http.ServeContent(w, r, "Survey_Results.xlsx", stat.ModTime(), file)
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
	excFlag := flag.Bool("excel", false, "Use this flag to create an excel table")
	flag.Parse()
	if *pushFlag {
		pushTest(db)
		return
	}
	if *excFlag {
		makeTable(db)
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/choose", chooseHandler)
	r.HandleFunc("/survey/{id}", surveyHandler)
	r.HandleFunc("/submit_survey", submitSurveyHandler)
	r.HandleFunc("/get_table", generateTable)
	r.HandleFunc("/download_table", downloadTable)
	r.HandleFunc("/result", resultHandler)
	http.Handle("/", r)

	srv := &http.Server{
		Addr: ":8080",
	}
	// Запуск сервера в горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Printf("Server started serving on port: %s", srv.Addr)

	// Ожидание сигнала остановки сервера
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server shutting down...")

	// Создание контекста с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Завершение работы сервера
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server gracefully stopped")
}
