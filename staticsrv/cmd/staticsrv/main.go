package main

import (
	"context"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
	"os/signal"
	"staticsrv/internal/api/server"
	"staticsrv/internal/app/version"
	"syscall"
)

// Config задает параметры конфигурации приложения
type Config struct {
	Port        string `envconfig:"PORT" default:"8080"`
	StaticsPath string `envconfig:"STATICS_PATH" default:"./static"`
	//StaticsPath string `envconfig:"STATICS_PATH" default:"../test_static"`
}

func main() {
	config := new(Config)
	err := envconfig.Process("", config)
	if err != nil {
		log.Fatalf("Can't process config: %v", err)
	}

	//r := http.NewServeMux()
	//r.Handle("/__heartbeat__", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	//fs := http.FileServer(http.Dir(config.StaticsPath))
	//r.Handle("/", fs)

	//go func() {
	//	log.Printf("start server on port: %s", config.Port)
	//	err = http.ListenAndServe(":"+config.Port, r)
	//	if err != nil {
	//		log.Fatalf("Error while serving: %v", err)
	//	}
	//}()

	info := server.VersionInfo{
		Version: version.Version,
		Commit:  version.Commit,
		Build:   version.Build,
	}

	srv := server.New(info, config.Port, config.StaticsPath)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		err := srv.Serve(ctx)
		if err != nil {
			log.Printf("start server on port: %s", config.Port)
			log.Println(fmt.Errorf("serve: %w", err))
			return
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case killSignal := <-interrupt:
		switch killSignal {
		case os.Interrupt:
			log.Print("Got SIGINT...")
		case syscall.SIGTERM:
			log.Print("Got SIGTERM...")
		}
	}

	cancel()
}
