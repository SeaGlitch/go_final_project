package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Предоставляет список текущих задач
func TasksH(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error": "метод не распознан"}`, http.StatusMethodNotAllowed)
			return
		}

		const limit = 50
		var rows *sql.Rows
		var err error

		search := r.FormValue("search")
		if search == "" {
			rows, err = db.Query("SELECT * FROM scheduler ORDER BY date LIMIT ?", limit)
		} else {

			searchDate, errS := time.Parse("02.01.2006", search)
			if errS == nil {
				rows, err = db.Query("SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT ?", searchDate.Format("20060102"), limit)
			} else {
				search = "%" + search + "%"
				rows, err = db.Query("SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?", search, search, limit)
			}
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var tasks []Task
		for rows.Next() {
			var task Task
			if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tasks = append(tasks, task)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if tasks == nil {
			tasks = []Task{} //Заготовка для пустого списка
		}

		resp := struct {
			Tasks []Task `json:"tasks"` //Формирование структуры для ответа
		}{
			Tasks: tasks,
		}

		//Кодирование JSON-ответа
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		err = json.NewEncoder(w).Encode(resp)

		if err != nil {
			log.Println("Error encoding JSON:", err)
			http.Error(w, `{"error": "Ошибка кодирования JSON"}`, http.StatusInternalServerError)
			return
		}

	}
}
