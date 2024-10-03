package db

import (
	"database/sql"
	"fmt"
	"github.com/ivanzharinov/go_final_project/internal/utils"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      int64
	Date    string
	Title   string
	Comment string
	Repeat  string
}

var db *sql.DB

func InitDB() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Ошибка при получении текущего рабочего каталога: %v", err)
	}
	dbFile := filepath.Join(wd, "scheduler.db")

	dbExists := true
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		dbExists = false
	} else if err != nil {
		log.Fatalf("Ошибка при проверке существования файла базы данных: %v", err)
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatalf("Ошибка при открытии базы данных: %v", err)
	}

	if !dbExists {
		createTables()
	} else {
		fmt.Println("База данных уже существует")
	}
}

func createTables() {
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

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Ошибка при создании таблиц: %v", err)
	}

	fmt.Println("Таблицы созданы успешно")
}

func AddTask(t Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	res, err := db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return 0, fmt.Errorf("Ошибка при добавлении задачи: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Ошибка при получении ID последней вставленной записи: %v", err)
	}

	return id, nil
}

func GetUpcomingTasks() ([]Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при выполнении запроса: %v", err)
	}
	defer rows.Close()

	tasks := []Task{}
	now := time.Now()

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("Ошибка при чтении строки из результата: %v", err)
		}

		taskDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			log.Printf("Ошибка при разборе даты задачи ID %d: %v", task.ID, err)
			continue
		}

		if taskDate.Before(now) || taskDate.Equal(now) {
			nextDateStr, err := utils.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				log.Printf("Ошибка при вычислении следующей даты для задачи ID %d: %v", task.ID, err)
				continue
			}
			task.Date = nextDateStr
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Ошибка при обработке результатов запроса: %v", err)
	}

	// Сортировка задач по дате
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Date < tasks[j].Date
	})

	// Ограничение списка задач до 50
	if len(tasks) > 50 {
		tasks = tasks[:50]
	}

	return tasks, nil
}
