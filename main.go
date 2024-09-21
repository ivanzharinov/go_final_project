package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("TODO_PORT")
	webDir := "./web"

	r := chi.NewRouter()

	fileServer := http.FileServer(http.Dir(webDir))
	r.Handle("/*", fileServer)

	fmt.Printf("Сервер запущен и слушает порт %s", port)
	http.ListenAndServe(":"+port, r)
}
