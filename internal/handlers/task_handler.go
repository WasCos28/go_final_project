package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"my_education/go/go_final_project/internal/logic"
	"net/http"
	"time"
)

type Task struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что запрос выполнен методом POST
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"только метод POST поддерживается"}`, http.StatusMethodNotAllowed)
			return
		}

		var task Task

		// Десериализуем JSON-запрос
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&task)
		if err != nil {
			http.Error(w, `{"error":"ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}

		// Проверяем обязательное поле title
		if task.Title == "" {
			http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
			return
		}

		// Проверяем и парсим поле date
		var taskDate time.Time
		if task.Date == "" {
			taskDate = time.Now().Truncate(24 * time.Hour)
		} else {
			taskDate, err = time.Parse("20060102", task.Date)
			if err != nil {
				http.Error(w, `{"error":"Дата указана в неверном формате"}`, http.StatusBadRequest)
				return
			}
		}

		// Проверяем, если дата меньше сегодняшнего дня
		now := time.Now().Truncate(24 * time.Hour)
		if taskDate.Before(now) {
			if task.Repeat == "" {
				taskDate = now // Устанавливаем на текущую дату
			} else {
				taskDateStr, err := logic.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusBadRequest)
					return
				}
				taskDate, _ = time.Parse("20060102", taskDateStr)
			}
		}

		// Вставляем задачу в базу данных
		query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
		res, err := db.Exec(query, taskDate.Format("20060102"), task.Title, task.Comment, task.Repeat)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка при добавлении задачи: %v"}`, err), http.StatusInternalServerError)
			return
		}

		// Получаем ID вставленной записи
		id, err := res.LastInsertId()
		if err != nil {
			http.Error(w, `{"error":"Ошибка при получении идентификатора задачи"}`, http.StatusInternalServerError)
			return
		}

		// Отправляем успешный ответ
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		response := fmt.Sprintf(`{"id":"%d"}`, id)
		fmt.Fprintln(w, response)
	}
}
