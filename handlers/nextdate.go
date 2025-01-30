package handlers

import (
	"database/sql"
	"encoding/json"

	"net/http"
	"time"

	"github.com/SeaGlitch/go_final_project/tasks"
)

// Вычисление следующей даты задачи
func NextDateH(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := r.FormValue("now")
		date := r.FormValue("date")
		repeat := r.FormValue("repeat")

		nowTime, err := time.Parse("20060102", now)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		nextDate, err := tasks.NextDate(nowTime, date, repeat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(nextDate))

		//fmt.Fprintln(w, nextDate)

	}
}

// Функция отмечает задачу выполненной, запрашивает следующую дату
func DoneTaskH(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error": "метод не распознан"}`, http.StatusMethodNotAllowed)
			return
		}

		id := r.FormValue("id")
		if id == "" {
			http.Error(w, `{"error": "неправильный id"}`, http.StatusBadRequest)
			return
		}

		var task Task
		err := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "задача не найдена"}`, http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		now := time.Now()
		nextDate, err := tasks.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if nextDate == "" {
			_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		} else {
			_, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{})
	}
}
