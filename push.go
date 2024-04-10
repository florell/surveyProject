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
	jsonData := `{
		"ID": 2,
		"Title": "Копинг-стратении Лазарус",
		"Description": "Тест поможет вам осознать, какие стратегии преодоления стресса вы чаще всего используете, и даст понимание, какие методы могут быть наиболее эффективными для вас в управлении жизненными трудностями.",
		"Questions": [
			{
				"Title": "Сосредотачивался на том, что мне нужно делать дальше, на следующем шаге.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Начинал что-то делать, зная, что это всё равно не будет работать: главное – делать хоть что-нибудь.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Пытался склонить вышестоящих к тому, чтобы они изменили свое мнение.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Говорил с другими, чтобы больше узнать о ситуации.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Критиковал и укорял себя.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Пытался не сжигать за собой мосты.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Надеялся на чудо.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Смирялся с судьбой: бывает, что мне не везет.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Вел себя, как будто ничего не произошло.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Старался не показывать своих чувств.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Пытался увидеть плюсы этой ситуации.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Спал больше обычного.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Срывал свою досаду на тех, кто навлек на меня проблемы.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Искал сочувствия и понимания у кого-нибудь.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Во мне возникла потребность выразить себя творчески.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Пытался забыть всё это.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Обращался за помощью к специалистам.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Менялся или рос как личность в положительную сторону.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Извинялся или старался всё загладить.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Составлял план действий.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Старался дать какой-то выход своим чувствам.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Понимал, что я сам вызывал эту проблему.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Набирался опыта в этой ситуации.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Говорил с кем-либо, кто мог конкретно помочь в этой ситуации.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Пытался улучшить свое самочувствие едой, выпивкой, курением или лекарствами.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Рисковал напропалую.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Старался действовать не слишком поспешно, не доверяясь первому порыву.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Находил новую веру во что-то.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Вновь открывал для себя что-то важное.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Что-то менял так, что всё улаживалось.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "В целом избегал общения с людьми.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Не допускал это до себя, стараясь об этом особенно не задумываться.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Спрашивал совета у родственника или друга, которых уважал.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Старался, чтобы другие не узнали, как плохо обстоят дела.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Отказывался воспринимать это слишком серьезно.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Говорил с кем-то о том, что я чувствую.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Стоял на своем и боролся за то, что хотел.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Вымещал это на других людях.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Пользовался прошлым опытом – мне приходилось уже попадать в такие ситуации.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Знал, что надо делать, и удваивал свои усилия, чтобы всё наладить.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Отказывался верить, что это действительно произошло.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Я давал себе обещание, что в следующий раз всё будет по-другому.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Находил пару других способов решения проблемы.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Старался, чтобы мои эмоции не слишком мешали мне в других делах.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Что-то менял в себе.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Хотел, чтобы всё это скорее как-то образовалось или кончилось.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Представлял себе, фантазировал, как всё это могло бы обернуться.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Молился.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Прокручивал в уме, что мне сказать или сделать.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
				]
			},
			{
				"Title": "Думал о том, как бы в данной ситуации действовал человек, которым я восхищаюсь, и старался подражать ему.",
				"Answers": [
					{
						"Text": "Никогда",
						"Value": 0
					},
					{
						"Text": "Редко",
						"Value": 1
					},
					{
						"Text": "Иногда",
						"Value": 2
					},
					{
						"Text": "Часто",
						"Value": 3
					}
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

	fmt.Println("Survey data inserted into the database successfully.")
}
