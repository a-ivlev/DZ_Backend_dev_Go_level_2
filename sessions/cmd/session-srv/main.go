package main

import (
	"log"
	"net/http"
	"session-srv/internal/db/redisDB/client"
	"time"
)

func main() {
	const (
		host        = "localhost"
		redisPort   = "6379"
		servicePort = "8080"
	)
	ttl := 1 * time.Hour
	client, err := client.NewRedisClient(host, redisPort, ttl)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	http.HandleFunc("/", client.RootHandler)
	http.HandleFunc("/login", client.LoginHandler)
	http.HandleFunc("/logout", client.LogoutHandler)

	log.Printf("starting server at :%s", servicePort)
	log.Fatal(http.ListenAndServe(":"+servicePort, nil))
}
