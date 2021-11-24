package main

import (
	"DZ_Backend_dev_Go_level_2/client/internal/api/chiRouter"
	"log"
	"net/http"
	"os"
)

func main() {
	cliPort := os.Getenv("CLI_PORT")
	if cliPort == "" {
		log.Fatal("unknown CLI_PORT = ", cliPort)
	}

	//r := chi.NewRouter()
	//r.Get("/", handlers.HomePage)
	//r.Post("/", handlers.HomePage)
	//r.Get("/{short}", handlers.RedirectPage)
	//r.Get("/stat/{stat}", handlers.StatPage)
	//
	//r.Get("/__heartbeat__", func(w http.ResponseWriter, r *http.Request) {})

	r := chiRouter.NewChiRouter()

	err := http.ListenAndServe(":"+cliPort, r)
	if err != nil {
		log.Println("Client shortener stopped...")
	}
}
