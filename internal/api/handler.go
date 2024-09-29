package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

func RegisterAPIRoutes(r *chi.Mux) {
	r.Get("/api/nextdate", HandleNextDate)
	r.Post("/api/task", HandleAddTask)
}

func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "недопустимый формат даты", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dateStr, repeatStr)
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
		AddTask(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
