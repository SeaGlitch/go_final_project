package main

import (
	"go_final_project/database"
	"go_final_project/handlers"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	//Путь к базе данных
	dbPath := os.Getenv("TODO_DBFILE") // Путь из переменной окружения
	if dbPath == "" {
		dbPath = "database/scheduler.db" // Путь по умолчанию
	}

	// Подключение базы данных
	db := database.ConnectDB(dbPath)
	defer db.Close()

	// Таблица базы данных
	database.TableDB(db, dbPath)

	//Обработчики
	mux := http.NewServeMux()
	mux.HandleFunc("/api/nextdate", handlers.NextDateH(db))
	mux.HandleFunc("/api/task", handlers.TaskH(db))
	mux.HandleFunc("/api/tasks", handlers.TasksH(db))
	mux.HandleFunc("/api/task/done", handlers.DoneTaskH(db))

	mux.Handle("/*", http.FileServer(http.Dir("./web")))

	// Получение порта
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	// Запуск сервера
	log.Printf("Сервер запущен на http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
