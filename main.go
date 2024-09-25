package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

func main() {
	CheckDB()
	port := os.Getenv("TODO_PORT")
	webDir := "./web"

	r := chi.NewRouter()

	fileServer := http.FileServer(http.Dir(webDir))
	r.Handle("/*", fileServer)

	fmt.Printf("Сервер запущен и слушает порт %s", port)
	http.ListenAndServe(":"+port, r)
}

func CheckDB() {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	if install {
		create()
	} else {
		fmt.Println("Таблица уже существует")
	}
}

func create() {
	db, err := sql.Open("sqlite", "./scheduler.db")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date CHAR(8) NOT NULL,
        title TEXT NOT NULL,
        comment TEXT,
        repeat VARCHAR(128),
        UNIQUE(date, title)
    );

    CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
    `

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Таблица создана успешно")
}
