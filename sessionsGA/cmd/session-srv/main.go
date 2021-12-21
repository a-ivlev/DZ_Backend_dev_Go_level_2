package main

import (
	"log"
	"net/http"
	"session-srv/internal/app/sessions"
	"session-srv/internal/db/redisDB"
	"session-srv/internal/handlers"
	"time"
)

func main() {
	const (
		host        = "localhost"
		redisPort   = "6379"
		servicePort = "8080"
	)
	ttl := 1 * time.Hour
	client, err := redisDB.NewRedisClient(host, redisPort, ttl)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	sc := sessions.NewSessionCache(client)
	r := handlers.NewRouter(sc, ttl)

	log.Printf("starting server at :%s", servicePort)
	log.Fatal(http.ListenAndServe(":"+servicePort, r))
}
