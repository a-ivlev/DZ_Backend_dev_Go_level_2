package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	VersionInfo
	srv         http.Server
	port        string
	staticsPath string
}

type VersionInfo struct {
	Version string
	Commit  string
	Build   string
}

func New(info VersionInfo, port string, staticsPath string) *Server {
	s := &Server{
		VersionInfo: info,
		port:        port,
		staticsPath: staticsPath,
	}

	s.srv = http.Server{
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
	}

	return s
}

func (s Server) Serve(ctx context.Context) error {
	r := http.NewServeMux()
	s.initHandlers(r)

	go func() {
		log.Printf("start server on port: %s", s.port)
		err := http.ListenAndServe(":"+s.port, r)
		if err != nil {
			log.Fatalf("Error while serving: %v", err)
		}
	}()

	<-ctx.Done()

	return s.srv.Shutdown(ctx)
}

func (s Server) initHandlers(r *http.ServeMux) {
	r.Handle("/__heartbeat__", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	r.Handle("/", http.FileServer(http.Dir(s.staticsPath)))
	r.Handle("/__version__", http.HandlerFunc(s.versionHandler))
}

func (s Server) versionHandler(w http.ResponseWriter, r *http.Request) {
	strJSON, err := json.Marshal(map[string]string{
		"version": s.VersionInfo.Version,
		"commit":  s.VersionInfo.Commit,
		"build":   s.VersionInfo.Build,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "func StatPage: error occured json marshaling stat page: %s", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(strJSON)
}
