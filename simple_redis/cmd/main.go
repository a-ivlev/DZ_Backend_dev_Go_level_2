package main

import (
	"log"
	"simple-redis/internal/db/redisDB"
	"simple-redis/internal/db/redisDB/client"
)

func main() {
	const (
		host = "localhost"
		port = "6379"
	)

	client, err := client.NewRedisClient(host, port)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	//if err := redisDB.BasicWork(client); err != nil {
	//	log.Fatal(err)
	//}

	if err := redisDB.WithStructWork(client); err != nil {
		log.Fatal(err)
	}
}
