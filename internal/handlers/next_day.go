package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"my_education/go/go_final_project/internal/logic"
)

// NextDateHandler обрабатывает запросы на вычисление следующей даты.
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Парсим дату now из строки
	now, err := time.Parse(logic.FormatDate, nowStr)
	if err != nil {
		http.Error(w, "некорректная дата now", http.StatusBadRequest)
		return
	}

	// Вызываем функцию для получения следующей даты
	nextDate, err := logic.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, "ошибка: "+err.Error(), http.StatusBadRequest)
		log.Printf("Ошибка при обработке next_day: %v", err)
	}

	// Возвращаем результат в формате текста
	fmt.Fprintln(w, nextDate)
}
