package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
)

// Подключение базы данных
func ConnectDB(dbPath string) *sql.DB {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("ошибка подключения к базе данных", err)
	}
	return db
}

// Таблица базы данных
func TableDB(db *sql.DB, dbPath string) {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), dbPath)
	log.Printf("dbFile:%s", dbFile)
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	// Создание таблицы базы данных
	if install {
		table := `CREATE TABLE IF NOT EXISTS scheduler (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					date CHAR(8) NOT NULL DEFAULT "",
					title TEXT NOT NULL DEFAULT "",
					comment TEXT NOT NULL DEFAULT "",
					repeat VARCHAR(128) NOT NULL DEFAULT ""
				);
					CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);
			`

		_, err := db.Exec(table)
		if err != nil {
			log.Fatal("ошибка при создании таблицы базы данных", err)
		}
		log.Println("База данных запущена.")
	} else {
		log.Println("База данных не существует.")
	}
}
