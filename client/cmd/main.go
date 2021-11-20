package main

import (
	"DZ_Backend_dev_Go_level_2/client/internal/handlers"
	"github.com/go-chi/chi/v5"

	"log"
	"net/http"
	"os"
)

func main() {
	cliPort := os.Getenv("CLI_PORT")
	if cliPort == "" {
		log.Fatal("unknown CLI_PORT = ", cliPort)
	}

	r := chi.NewRouter()
	r.Get("/", handlers.HomePage)
	r.Post("/", handlers.HomePage)
	r.Get("/{short}", handlers.RedirectPage)
	r.Get("/stat/{stat}", handlers.StatPage)

	err := http.ListenAndServe(":"+cliPort, r)
	if err != nil {
		log.Println("Client shortener stopped...")
	}
}
