package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Redirect struct {
	ShortLink string `json:"short_link"`
	FullLink  string `json:"full_link"`
	IPaddress string `json:"ip_address"`
}

func RedirectPage(w http.ResponseWriter, r *http.Request) {
	//srvHost := os.Getenv("SRV_HOST")
	//if srvHost == "" {
	//	log.Fatal("unknown SRV_HOST = ", srvHost)
	//}
	//
	//srvPort := os.Getenv("SRV_PORT")
	//if srvPort == "" {
	//	log.Fatal("unknown SRV_PORT = ", srvPort)
	//}
	//
	//redirectPath := chi.URLParam(r, "short")

	//redirect := fmt.Sprintf("http://%s:%s/%s", srvHost, srvPort, redirectPath)
	//http.Redirect(w, r, redirect, http.StatusFound)

	if r.Method == http.MethodGet {

		ipaddr := strings.Split(r.RemoteAddr, ":")
		//nolint:staticcheck
		//r = r.WithContext(context.WithValue(r.Context(), "IP_address", ipaddr[0]))

		shortLink := &Redirect{
			ShortLink: chi.URLParam(r, "short"),
			IPaddress: ipaddr[0],
		}

		strJSON, err := json.Marshal(&shortLink)
		if err != nil {
			fmt.Fprintf(w, "func StatPage: error occured json marshaling stat page")
		}

		srvHost := os.Getenv("SRV_HOST")
		if srvHost == "" {
			log.Fatal("unknown SRV_HOST = ", srvHost)
		}

		srvPort := os.Getenv("SRV_PORT")
		if srvPort == "" {
			log.Fatal("unknown SRV_PORT = ", srvPort)
		}

		srv := fmt.Sprintf("http://%s:%s/%s", srvHost, srvPort, shortLink.ShortLink)

		client := &http.Client{Timeout: time.Second * 2}
		req, err := http.NewRequest(http.MethodPost, srv, bytes.NewBuffer(strJSON))
		if err != nil {
			log.Println("func StatPage: error occurred NewRequest: ", err)
		}
		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			log.Println("func StatPage: error occurred client Do: ", err)
		}

		respRedirect := &Redirect{}
		defer res.Body.Close()
		if err = json.NewDecoder(res.Body).Decode(&respRedirect); err != nil {
			http.Error(w, "error unmarshal request", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, respRedirect.FullLink, http.StatusFound)
	}

}
