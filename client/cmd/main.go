package main

import (
	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/client/internal/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("SHORT_CLI_PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", handlers.HomePage)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println("Client shortener stopped...")
	}
}
