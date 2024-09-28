package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ivanzharinov/go_final_project/internal/api"
	"github.com/ivanzharinov/go_final_project/internal/db"
	"net/http"
)

func main() {
	db.InitDB()
	//port := os.Getenv("TODO_PORT")
	port := "7540"
	webDir := "./web"

	r := chi.NewRouter()

	r.Get("/api/nextdate", api.HandleNextDate)

	fileServer := http.FileServer(http.Dir(webDir))
	r.Handle("/*", fileServer)

	fmt.Printf("Сервер запущен и слушает порт %s", port)
	http.ListenAndServe(":"+port, r)
}
