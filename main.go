package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ivanzharinov/go_final_project/internal/db"
	"github.com/ivanzharinov/go_final_project/internal/transport"
	"log"
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

	fmt.Printf("Сервер запущен и слушает порт %s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
