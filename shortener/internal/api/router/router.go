package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/a-ivlev/DZ_Backend_dev_Go_level_2/shortener/internal/app/shortenerBL"
)

type Router struct {
	*http.ServeMux
	short *shortenerBL.ShortenerBL
}

func NewRouter(short *shortenerBL.ShortenerBL) *Router {
	r := &Router{
		ServeMux: http.NewServeMux(),
		short:    short,
	}
	r.Handle("/create",
		// r.AuthMiddleware(
		//r.AuthMiddleware(
		http.HandlerFunc(r.CreateShortener),
		//),
		// ),
	)
	r.Handle("/RedirectAPI", r.AuthMiddleware(http.HandlerFunc(r.RedirectAPI)))
	//r.Handle("/{uuid}", r.AuthMiddleware(http.HandlerFunc(r.RedirectAPI)))
	//r.Handle("/delete", r.AuthMiddleware(http.HandlerFunc(r.DeleteUser)))
	//r.Handle("/search", r.AuthMiddleware(http.HandlerFunc(r.SearchUser)))
	//r.Handle("/whoami", r.AuthMiddleware(http.HandlerFunc(r.Whoami)))
	return r
}

func (rt *Router) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if u, p, ok := r.BasicAuth(); !ok || !(u == "admin" && p == "admin") {
				fmt.Println("ok: ", ok, "user: ", u, "password", p)
				http.Error(w, "unautorized", http.StatusUnauthorized)
				return
			}

			// r = r.WithContext(context.WithValue(r.Context(), 1, 0))
			next.ServeHTTP(w, r)
		},
	)
}

type Shortener struct {
	ShortLink string `json:"short_link"`
	FullLink  string `json:"full_link"`
	// count     int
	CreatedAt time.Time `json:"created_at"`
}

type Link struct {
	Link string `json:"full_link"`
}

func (rt *Router) CreateShortener(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println(r.Method)
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	fulllink := Link{}
	if err := json.NewDecoder(r.Body).Decode(&fulllink); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	shortener := shortenerBL.Shortener{
		FullLink: fulllink.Link,
	}

	newShort, err := rt.short.CreateShort(r.Context(), shortener)
	if err != nil {
		http.Error(w, "error when creating", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(
		Shortener{
			ShortLink: newShort.ShortLink,
			FullLink:  newShort.FullLink,
			CreatedAt: newShort.CreatedAt,
		},
	)
}

func (rt *Router) Redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	//shortlink := Link{}
	//if err := json.NewDecoder(r.Body).Decode(&shortlink); err != nil {
	//	http.Error(w, "bad request", http.StatusBadRequest)
	//	return
	//}
	//link := r.URL.String()//fmt.Sprintf("%s/%s", r.URL, r.URL.)
	//fmt.Println("link", link)
	shortener := shortenerBL.Shortener{
		ShortLink: r.RequestURI,
	}

	shortBD, err := rt.short.GetFullLink(r.Context(), shortener)
	if err != nil {
		http.Error(w, "error when creating", http.StatusInternalServerError)
		return
	}

	client := http.Client{Timeout: time.Second * 2}
	req, _ := http.NewRequest("GET", shortBD.FullLink, nil)
	res, err := client.Do(req)
	if err != nil {
		http.Error(w, "error when creating", http.StatusInternalServerError)
		return
	}

	if res.StatusCode > 399 {
		http.Error(w, "error when creating", http.StatusInternalServerError)
		return
	}

	//w.WriteHeader(res.StatusCode)
	//io.Copy(w, res.Body)

	flusher, ok := w.(http.Flusher)
	if !ok {
		fmt.Println("expected http.ResponseWriter to be an http.Flusher")
	}
	w.Header().Set("X-Content-Type-Options", "nosniff")
	flusher.Flush()
}

func (rt *Router) Whoami(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Ваш IP %s", r.RemoteAddr)
}

func (rt *Router) RedirectAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	shortlink := Link{}
	if err := json.NewDecoder(r.Body).Decode(&shortlink); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	shortener := shortenerBL.Shortener{
		ShortLink: shortlink.Link,
	}

	shortBD, err := rt.short.GetFullLink(r.Context(), shortener)
	if err != nil {
		http.Error(w, "error when creating", http.StatusInternalServerError)
		return
	}

	client := http.Client{Timeout: time.Second * 2}
	req, _ := http.NewRequest("GET", shortBD.FullLink, nil)
	res, err := client.Do(req)
	if err != nil {
		http.Error(w, "error when creating", http.StatusInternalServerError)
		return
	}

	if res.StatusCode > 399 {
		http.Error(w, "error when creating", http.StatusInternalServerError)
		return
	}

	//w.WriteHeader(res.StatusCode)
	//io.Copy(w, res.Body)

	flusher, ok := w.(http.Flusher)
	if !ok {
		fmt.Println("expected http.ResponseWriter to be an http.Flusher")
	}
	w.Header().Set("X-Content-Type-Options", "nosniff")
	flusher.Flush()
}
