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
		"ID": 1,
		"Title": "Шкала семейных отношений",
		"Description": "Тест поможет вам осознать и оценить качество ваших взаимоотношений в семье, предоставляя вам инсайты для улучшения семейной гармонии и понимания между членами семьи.",
		"Questions": [
			{
				"Title": "Члены нашей семьи оказывают реальную помощь и поддержку друг другу.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Члены нашей семьи часто скрывают свои чувства.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье мы часто ссоримся.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье мы не часто делаем что-либо самостоятельно (отдельно от других членов).",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Мы считаем важным – быть лучшим в любом деле, которое делаешь.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Мы часто говорим о политических и социальных проблемах.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Мы проводим большую часть выходных дней и вечеров дома.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Члены нашей семьи довольно часто смотрят передачи на морально-нравственные темы.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Вся деятельность нашей семьи довольно тщательно планируется.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье редко кто-то командует.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Мы часто дома «убиваем» время.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В своем доме мы говорим все, что хотим.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Члены нашей семьи редко открыто сердятся.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье очень поощряется независимость.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Жизненный успех (продвижение в жизни) очень важен в нашей семье.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Мы редко ходим на лекции, спектакли и концерты.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Друзья часто приходят к нам в гости.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Мы считаем, что семья не несет ответственности за своих членов.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Мы, как правило, очень опрятны и организованы.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Число правил, которым мы следуем в нашей семье, невелико.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Мы вкладываем много энергии в домашние дела.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Трудно «разрядиться» дома, не расстроив кого-нибудь.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Члены нашей семьи иногда настолько разозлятся, что могут швырять вещи.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье мы обдумываем свои дела в одиночку.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Для нас не очень важно, сколько зарабатывает человек.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье является очень важным узнавать о новых вещах, событиях, фактах.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Никто в нашей семье не занимается активно спортом.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Мы часто говорим на морально-нравственные темы.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашем доме часто трудно бывает найти вещь, которая требуется в данный момент.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "У нас есть один член семьи, который принимает большинство решений.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье существует чувство единства.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Мы рассказываем друг другу о своих личных проблемах.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Члены нашей семьи редко выходят из себя.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье мы приходим и уходим, когда захотим.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В любом деле мы верим в соревнование и девиз «Пусть победит сильнейший».",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Мы не очень интересуемся культурной жизнью.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Мы часто ходим в кино, театр, туристические походы.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Высокая нравственность не является уделом нашей семьи.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Быть пунктуальным в нашей семье очень важно.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашем доме все делается по заведенному порядку.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Мы редко вызываемся добровольно, когда что-то нужно сделать дома.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Если нам хочется что-то сделать экспромтом, мы часто тут же собираемся и делаем это.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Члены нашей семьи часто критикуют друг друга.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье очень мало тайн.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Мы всегда стремимся делать дело так, чтобы в следующий раз получилось намного лучше.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "У нас редко бывают интеллектуальные дискуссии.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Все в нашей семье имеют одно или несколько хобби.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "У членов семьи строгие понятия о том, что правильно и что неправильно.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье все часто меняют мнение о домашних делах.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье придается большое значение соблюдению правил.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Мы стараемся делать все во имя сплоченности нашей семьи.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Если у нас в семье начнешь жаловаться, кто-то обычно расстроится.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Члены нашей семьи иногда могут ударить друг друга.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Члены нашей семьи обычно полагаются сами на себя, если возникает какая-то проблема.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Членов нашей семьи мало волнует продвижение по работе, школьные отметки.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Кто-то в нашей семье играет на музыкальном инструменте.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Члены нашей семьи принимают мало участия в развлекательных мероприятиях.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Мы убеждены, что существуют некоторые вещи, которые надо принимать на веру.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Члены нашей семьи содержат свои комнаты в порядке.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В семейных решениях все имеют равное право голоса.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье очень слабо развит дух коллективизма.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье открыто обсуждаются денежные дела и оплата счетов.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Если в нашей семье возникают разногласия, мы изо всех сил стараемся «сгладить углы» и сохранить мир.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Члены нашей семьи усиленно поощряют друг друга отстаивать свои права.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье мы не очень стремимся к успеху.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Члены семьи читают много книг.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Члены семьи иногда посещают курсы или берут уроки по своим интересам и увлечениям (помимо школы).",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье у каждого свои понятия о том, что правильно, а что неправильно.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Обязанности каждого в нашей семье четко определены.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье мы можем делать все, что хотим.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Мы редко по-настоящему ладим друг с другом.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье мы обычно следим, что говорим друг другу.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Члены семьи часто пытаются быть в чем-то выше или превзойти один другого.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашем доме трудно побыть одному, чтобы это кого-нибудь не обидело.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "«Делу – время, потехе – час» – таково правило нашей семьи.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье мы чаще смотрим телевизор, чем читаем книги.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Члены семьи часто выходят «в свет».",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Наша семья придерживается строгих моральных правил.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье с деньгами обращаются не очень бережно.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье царит правило: «Всяк сверчок знай свой шесток!».",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье всем уделяется достаточно много времени и внимания.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье часто возникают спонтанные дискуссии.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье мы считаем, что повышением голоса ничего не добьешься.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье не поощряется, чтобы каждый высказывался сам за себя.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Членов нашей семьи часто сравнивают с другими людьми в отношении того, как они успевают на работе или в школе.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье по-настоящему любят музыку, живопись, литературу.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Главная форма развлечения у нас – смотреть телевизор.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "Члены нашей семьи верят в торжество справедливости.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье посуда моется сразу после еды.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
					}
				]
			},
			{
				"Title": "В нашей семье немногое проходит безнаказанно.",
				"Answers": [
					{
						"Text": "Да",
						"Value": 1
					},
					{
						"Text": "Нет",
						"Value": 0
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
