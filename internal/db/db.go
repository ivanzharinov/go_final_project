package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDB() {
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
