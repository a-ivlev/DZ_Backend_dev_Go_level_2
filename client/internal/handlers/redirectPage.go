package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
)

func RedirectPage(w http.ResponseWriter, r *http.Request)  {
	srvHost := os.Getenv("SRV_HOST")
	if srvHost == "" {
		srvHost = "localhost"
	}

	srvPort := os.Getenv("SRV_PORT")
	if srvPort == "" {
		srvPort = "8035"
	}

	redirectPath := chi.URLParam(r, "short")

	redirect := fmt.Sprintf("http://%s:%s/%s", srvHost, srvPort, redirectPath)

	http.Redirect(w, r, redirect, http.StatusFound)
}
