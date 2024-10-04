package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ivanzharinov/go_final_project/internal/db"
	"github.com/ivanzharinov/go_final_project/internal/transport"
	"net/http"
)

func main() {
	db.InitDB()
	port := "7540"
	webDir := "./web"

	r := chi.NewRouter()

	transport.RegisterAPIRoutes(r)

	fileServer := http.FileServer(http.Dir(webDir))
	r.Handle("/*", fileServer)

	fmt.Printf("Сервер запущен и слушает порт %s", port)
	http.ListenAndServe(":"+port, r)
}
