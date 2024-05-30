package main

import (
	"context"
	"database/sql"
	"errors"
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

func checkSessionData(session *sessions.Session) bool {
	patientID, ok1 := session.Values["patientID"].(int64)

	_, ok2 := session.Values["patientGender"].(string)

	_, ok3 := session.Values["patientAge"].(string)

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM patients WHERE id = ?", patientID).Scan(&count); err != nil {
		log.Fatalln("Error while checking session data:", err)
	}

	return ok1 && ok2 && ok3 && (count != 0)
}

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

		// log.Printf("%T %v %T %v\n", sex, sex, age, age)

		session.Values["patientID"] = insertedID
		session.Values["patientGender"] = sex
		session.Values["patientAge"] = age
		session.Values["patientName"] = name + " " + surname

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

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error getting session:", err)
		return
	}

	if !checkSessionData(session) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	patientName, ok := session.Values["patientName"].(string)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	patientGender, ok := session.Values["patientGender"].(string)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if patientGender == "male" {
		patientGender = "мужской"
	} else {
		patientGender = "женский"
	}

	patientAge, ok := session.Values["patientAge"].(string)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	patientId, ok := session.Values["patientID"].(int64)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	surveysWithCheck := []struct {
		Survey    types.Survey
		Completed bool
	}{}

	for rows.Next() {
		surveyWithCheck := struct {
			Survey    types.Survey
			Completed bool
		}{}
		err := rows.Scan(&surveyWithCheck.Survey.SurveyID, &surveyWithCheck.Survey.Title, &surveyWithCheck.Survey.Description)
		if err != nil {
			log.Fatal(err)
		}
		row := db.QueryRow("SELECT EXISTS (SELECT 1 FROM survey_results WHERE SurveyID = ? AND PatientID = ?) AS RecordExists", surveyWithCheck.Survey.SurveyID, patientId)
		err = row.Scan(&surveyWithCheck.Completed)
		if err != nil {
			log.Fatal(err)
		}
		surveysWithCheck = append(surveysWithCheck, surveyWithCheck)
	}

	tmpl := template.Must(template.ParseFiles("templates/choose.html"))

	data := struct {
		SurveysWithCheck []struct {
			Survey    types.Survey
			Completed bool
		}
		ID     string
		Age    string
		Name   string
		Gender string
	}{surveysWithCheck, strconv.Itoa(int(patientId)), patientAge, patientName, patientGender}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
		return
	}
}

func surveyHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error getting session:", err)
		return
	}

	if !checkSessionData(session) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	vars := mux.Vars(r)
	id := vars["id"]

	// Fetch survey data from the database using survey ID
	var survey types.Survey

	// Fetch survey data from the database
	row := db.QueryRow("SELECT id, title FROM surveys WHERE id = ?", id)
	if err := row.Scan(&survey.SurveyID, &survey.Title); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		log.Fatal(err)
	}

	// Fetch questions and answers for the survey from the database
	rows, err := db.Query("SELECT id, title, maxval FROM questions WHERE surveyid = ?", survey.SurveyID)
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
		err := rows.Scan(&question.QuestionID, &question.Title, &question.MaxValue)
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	patientGender, ok := session.Values["patientGender"].(string)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	patientAge, ok := session.Values["patientAge"].(string)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		// log.Println(key, values)
		if strings.HasPrefix(key, "question") {
			// Extract question number and answer ID
			questionID := strings.Split(key, "question_")[1]
			qIDint, ok := strconv.Atoi(questionID)
			if ok != nil {
				log.Println("r.Form - questionID:", ok)
			}
			intValue, ok := strconv.Atoi(values[0])
			if ok != nil {
				log.Println("r.Form - values[0]:", ok)
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

	// fmt.Println(surveyResults)
	// fmt.Println("^^^^ ", surveyID)

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
	case 7:
		analysis = handlers.FABHandler(&surveyResults)
	case 8:
		analysis = handlers.MMSEHandler(&surveyResults)
	case 9:
		analysis = handlers.CDTHandler(&surveyResults)
	case 10:
		analysis = handlers.VECHandler(&surveyResults)
	case 11:
		analysis = handlers.ACEHandler(&surveyResults)
	case 12:
		analysis = handlers.SFTHandler(&surveyResults)
	case 13:
		analysis = handlers.PFTHandler(&surveyResults)
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
		log.Println("Error inserting result into db (prepare)", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(stmt *sql.Stmt) {
		if err := stmt.Close(); err != nil {
			log.Println(err)
		}
	}(stmt)

	// fmt.Println("Analysis:", string(analysis))
	// Execute the SQL statement
	_, err = stmt.Exec(patientID, surveyID, string(analysis), string(analysis))
	if err != nil {
		log.Println("Error inserting result into db (exec)", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect after successful form submission
	http.Redirect(w, r, "/result?survey_id="+strconv.Itoa(surveyID)+"&patient_id="+fmt.Sprintf("%d", patientID), http.StatusSeeOther)
}

func generateTable(w http.ResponseWriter, r *http.Request) {
	err := makeTable(db)
	if err != nil {
		http.Error(w, "Ошибка при создании таблицы", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func downloadTable(w http.ResponseWriter, r *http.Request) {
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Println(fmt.Sprintf("Error deleting %s file:", "Survey_Results.xlsx"), err)
		}
	}("Survey_Results.xlsx")

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

func generateConclusion(w http.ResponseWriter, r *http.Request) {
	patientId := r.URL.Query().Get("patient_id")
	err := makeConclusion(db, patientId)
	if err != nil {
		http.Error(w, "Ошибка при создании заключения", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func downloadConclusion(w http.ResponseWriter, r *http.Request) {
	patientId := r.URL.Query().Get("patient_id")
	fileName := "conclusion_" + patientId + ".docx"

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Println(fmt.Sprintf("Error deleting %s file:", fileName), err)
		}
	}(fileName)

	file, err := os.Open(fileName)
	if err != nil {
		http.Error(w, "Ошибка при открытии файла заключения", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Отправляем файл пользователю
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	stat, _ := file.Stat()
	http.ServeContent(w, r, fileName, stat.ModTime(), file)
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	patientId := r.URL.Query().Get("patient_id")
	surveyId := r.URL.Query().Get("survey_id")

	var resD string
	if err := db.QueryRow(
		"SELECT Result FROM survey_results WHERE PatientID = ? AND SurveyID = ?", patientId, surveyId,
	).Scan(&resD); err != nil {
		http.Error(w, fmt.Sprintf("Error getting results from database: %s", err.Error()), http.StatusInternalServerError)
	}

	tmpl := template.Must(template.ParseFiles("templates/results.html"))

	if err := tmpl.Execute(w, resD); err != nil {
		log.Fatal(err)
	}

}

func main() {

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	name := os.Getenv("DB_NAME")

	var err error
	for {
		db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, pass, host, name))
		if err != nil {
			log.Println("Error opening database:", err)
		} else {
			err = db.Ping()
			if err == nil {
				break
			}
			log.Println("Waiting for database to be ready:", err)
			db.Close()
		}
		time.Sleep(5 * time.Second) // Повторяем пинг каждые 5 секунд
	}
	
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Println(err)
		}
	}(db)

	pushTest(db)

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/choose", chooseHandler)
	r.HandleFunc("/survey/{id}", surveyHandler)
	r.HandleFunc("/submit_survey", submitSurveyHandler)
	r.HandleFunc("/get_table", generateTable)
	r.HandleFunc("/download_table", downloadTable)
	r.HandleFunc("/get_conclusion", generateConclusion)
	r.HandleFunc("/download_conclusion", downloadConclusion)
	r.HandleFunc("/result", resultHandler)
	http.Handle("/", r)

	js := http.FileServer(http.Dir("src/script"))
	css := http.FileServer(http.Dir("src/styles"))
	img := http.FileServer(http.Dir("src/img"))

	http.Handle("/script/", http.StripPrefix("/script/", js))
	http.Handle("/styles/", http.StripPrefix("/styles/", css))
	http.Handle("/img/", http.StripPrefix("/img/", img))

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
