package main

import (
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

// Config задает параметры конфигурации приложения
type Config struct {
	Port        string `envconfig:"PORT" default:"8080"`
	StaticsPath string `envconfig:"STATICS_PATH" default:"./static"`
}

func main() {
	config := new(Config)
	err := envconfig.Process("", config)
	if err != nil {
		log.Fatalf("Can't process config: %v", err)
	}

	r := http.NewServeMux()
	r.Handle("/__heartbeat__", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	fs := http.FileServer(http.Dir(config.StaticsPath))
	r.Handle("/", fs)

	log.Printf("start server on port: %s", config.Port)
	err = http.ListenAndServe(":"+config.Port, r)
	if err != nil {
		log.Fatalf("Error while serving: %v", err)
	}

}
