package transport

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/ivanzharinov/go_final_project/internal/db"
	"github.com/ivanzharinov/go_final_project/internal/task"
	"github.com/ivanzharinov/go_final_project/internal/utils"
	"net/http"
	"strconv"
	"time"
)

func RegisterAPIRoutes(r *chi.Mux) {
	r.Get("/api/nextdate", HandleNextDate)
	r.Post("/api/task", HandleAddTask)
	r.Get("/api/tasks", task.Tasks)
	r.Get("/api/task", getTaskHandler)
	r.Put("/api/task", updateTaskHandler)
	r.Post("/api/task/done", handleTaskDone)
	r.Delete("/api/task", handleTaskDelete)
}

func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Недопустимый формат даты", http.StatusBadRequest)
		return
	}

	nextDate, err := utils.NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
}

func HandleAddTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		task.AddTask(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "отсутствует id"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "недопустимый параметр id"})
		return
	}

	task, err := db.GetTaskByID(id)
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"error": "задача не найдена"})
		return
	}

	response := map[string]string{
		"id":      strconv.FormatInt(task.ID, 10),
		"date":    task.Date,
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

var req struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверный формат JSON"})
		return
	}

	if req.ID == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Недопустимый формат идентификатора"})
		return
	}

	if req.Date == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Дата не может быть пустой"})
		return
	}

	_, err = time.Parse("20060102", req.Date)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверный формат даты"})
		return
	}

	if req.Title == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Заголовок не может быть пустым"})
		return
	}

	updatedTask := db.Task{
		ID:      id,
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	err = db.UpdateTask(updatedTask)
	if err != nil {
		if err == db.ErrTaskNotFound {
			respondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		} else {
			respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка обновления задачи"})
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{})
}

func handleTaskDone(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Отсутствует идентификатор задачи"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Недопустимый формат идентификатора"})
		return
	}

	task, err := db.GetTaskByID(id)
	if err != nil {
		if err == db.ErrTaskNotFound {
			respondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
			return
		}
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении задачи"})
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Ошибка удаления задачи"})
			return
		}
	} else {
		now := time.Now()
		nextDate, err := utils.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Ошибка вычисления следующей даты"})
			return
		}

		task.Date = nextDate
		err = db.UpdateTask(task)
		if err != nil {
			respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка обновления задачи"})
			return
		}
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{})
}

func handleTaskDelete(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Отсутствует идентификатор задачи"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Недопустимый формат идентификатора"})
		return
	}

	err = db.DeleteTask(id)
	if err != nil {
		if err == db.ErrTaskNotFound {
			respondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
			return
		}
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка удаления задачи"})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{})
}
func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}
