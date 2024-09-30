package api

import (
	"encoding/json"
	"fmt"
	"github.com/ivanzharinov/go_final_project/internal/utils"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ivanzharinov/go_final_project/internal/db"
)

type AddTaskRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type AddTaskResponse struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AddTaskResponse{Error: "не удалось прочитать тело запроса"})
		return
	}
	defer r.Body.Close()

	var req AddTaskRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AddTaskResponse{Error: "неверный формат JSON"})
		return
	}

	// проверка обязательного поля title
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AddTaskResponse{Error: "не указан заголовок задачи"})
		return
	}

	// обработка поля date
	var taskDate time.Time
	now := time.Now()

	if strings.TrimSpace(req.Date) == "" {
		req.Date = now.Format("20060102")
	}

	taskDate, err = time.Parse("20060102", req.Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AddTaskResponse{Error: "дата указана в неверном формате"})
		return
	}

	if taskDate.Before(now) {
		if strings.TrimSpace(req.Repeat) == "" {
			taskDate = now
		} else {
			nextDateStr, err := utils.NextDate(now, req.Date, req.Repeat)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(AddTaskResponse{Error: "неверное правило повторения"})
				return
			}
			taskDate, _ = time.Parse("20060102", nextDateStr)
		}
	}

	// создание объекта задачи
	newTask := db.Task{
		Date:    taskDate.Format("20060102"),
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	// добавление задачи в базу данных
	id, err := db.AddTask(newTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AddTaskResponse{Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(AddTaskResponse{ID: fmt.Sprintf("%d", id)})
}

type TaskResponseItem struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TasksResponse struct {
	Tasks []TaskResponseItem `json:"tasks"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func Tasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tasks, err := db.GetUpcomingTasks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	response := TasksResponse{Tasks: []TaskResponseItem{}}

	for _, t := range tasks {
		taskItem := TaskResponseItem{
			ID:      fmt.Sprintf("%d", t.ID),
			Date:    t.Date,
			Title:   t.Title,
			Comment: t.Comment,
			Repeat:  t.Repeat,
		}
		response.Tasks = append(response.Tasks, taskItem)
	}

	json.NewEncoder(w).Encode(response)
}
