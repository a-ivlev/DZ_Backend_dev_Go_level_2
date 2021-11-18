package server

import (
	"context"
	"net/http"
	"time"

	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/shortener/internal/app/shortenerBL"
	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/shortener/internal/app/starter"
)

var _ starter.APIServer = &Server{}

type Server struct {
	srv         http.Server
	shortenerBL *shortenerBL.ShortenerBL
}

func NewServer(addr string, h http.Handler) *Server {
	server := &Server{}

	server.srv = http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
	}
	return server
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	s.srv.Shutdown(ctx)
	cancel()
}

func (s *Server) Start(shortBL *shortenerBL.ShortenerBL) {
	s.shortenerBL = shortBL
	go s.srv.ListenAndServe()
}
