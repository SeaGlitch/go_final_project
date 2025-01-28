package handlers

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/SeaGlitch/go_final_project/tasks"

	"net/http"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Сортировка методов обработчика
func TaskH(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			addTaskH(w, r, db)
		case http.MethodPut:
			putTaskH(w, r, db)
		case http.MethodGet:
			getTaskH(w, r, db)
		case http.MethodDelete:
			delTaskH(w, r, db)
		default:
			http.Error(w, `{"error": "метод не распознан"}`, http.StatusMethodNotAllowed)
		}
	}
}

// Добавление новой задачи
func addTaskH(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, `{"error": "ошибка десериализации JSON"}`, http.StatusBadRequest)
		return
	}

	now := time.Now()
	nowSt := time.Now().Format("20060102")
	nowTD, err := time.Parse("20060102", nowSt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка правил повторений
	if task.Repeat != "" {
		switch task.Repeat[:1] {

		case "y":
			//ничего не делаем

		case "d":

			daysStr := task.Repeat[2:] // Извлекаем число дней
			days, err := strconv.Atoi(daysStr)
			if err != nil || days > 400 {
				response := map[string]string{"error": "недопустимое значение дней в правиле повторения"}
				json.NewEncoder(w).Encode(response)
				return
			}

		case "w":
			parts := strings.Fields(task.Repeat)
			if len(parts) < 2 {
				response := map[string]string{"error": "недопустимое значение дней в правиле повторения"}
				json.NewEncoder(w).Encode(response)
				return
			}

		case "m":
			parts := strings.Fields(task.Repeat)
			if len(parts) < 2 {
				response := map[string]string{"error": "недопустимое значение дней в правиле повторения"}
				json.NewEncoder(w).Encode(response)
				return
			}

		default:
			response := map[string]string{"error": "правило повторения указано в неправильном формате"}
			json.NewEncoder(w).Encode(response)
			return

		}

	}

	if task.Title == "" {
		http.Error(w, `{"error": "не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	if task.Date == "" {
		task.Date = now.Format("20060102")
	} else {
		parsedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			http.Error(w, `{"error": "дата представлена в формате, отличном от 20060102"}`, http.StatusBadRequest)
			return
		}
		if parsedDate.Before(nowTD) && task.Repeat != "" {
			nextDate, err := tasks.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		}
	}

	if task.Date < nowSt {
		task.Date = nowSt
	}

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Сообщаем номер созданной задачи
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

// Внесение изменений задачи
func putTaskH(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, `{"error": "ошибка десериализации JSON"}`, http.StatusBadRequest)
		return
	}

	now := time.Now()
	nowSt := time.Now().Format("20060102")
	nowTD, err := time.Parse("20060102", nowSt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка повторений
	if task.Repeat != "" {
		switch task.Repeat[:1] {

		case "y":
			//ничего не делаем

		case "d":
			daysStr := task.Repeat[2:] // Извлекаем число дней
			days, err := strconv.Atoi(daysStr)
			if err != nil || days > 400 {
				response := map[string]string{"error": "недопустимое значение дней в правиле повторения"}
				json.NewEncoder(w).Encode(response)
				return
			}

		case "w":
			parts := strings.Fields(task.Repeat)
			if len(parts) < 2 {
				response := map[string]string{"error": "недопустимое значение дней в правиле повторения"}
				json.NewEncoder(w).Encode(response)
				return
			}

		case "m":
			parts := strings.Fields(task.Repeat)
			if len(parts) < 2 {
				response := map[string]string{"error": "недопустимое значение дней в правиле повторения"}
				json.NewEncoder(w).Encode(response)
				return
			}

		default:
			response := map[string]string{"error": "правило повторения указано в неправильном формате"}
			json.NewEncoder(w).Encode(response)
			return

		}
	}

	intID, err := strconv.Atoi(task.ID)
	if err != nil || intID < 0 || intID > 32767 {
		http.Error(w, `{"error": "неправильный индекс"}`, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error": "не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	if task.Date == "" {
		task.Date = now.Format("20060102")
	} else {
		parsedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			http.Error(w, `{"error": "дата представлена в формате, отличном от 20060102"}`, http.StatusBadRequest)
			return
		}
		if parsedDate.Before(nowTD) && task.Repeat != "" {
			nextDate, err := tasks.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		}
	}

	if task.Date < nowSt {
		task.Date = nowSt
	}

	_, err = db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?", task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Сообщаем о выполнении
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{})
}

// Вызов задачи по id
func getTaskH(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, `{"error": "отсутствует ID"}`, http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, `{"error": "неправильный id"}`, http.StatusBadRequest)
		return
	}

	var task Task
	err = db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", idInt).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error": "задача не найдена"}`, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(task)
}

// Удаление задачи из списка
func delTaskH(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, `{"error": "отсутствует id"}`, http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, `{"error": "неправильный id"}`, http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", idInt)
	if err != nil {
		http.Error(w, `{"error": "ошибка удаления id"}`, http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Сообщаем о выполнении
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{})
}
