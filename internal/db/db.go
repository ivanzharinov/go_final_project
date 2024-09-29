package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      int64
	Date    string
	Title   string
	Comment string
	Repeat  string
}

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

func AddTask(t Task) (int64, error) {

	db, err := sql.Open("sqlite", "./scheduler.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
