package handler
//
//import (
//	"CourseProjectBackendDevGoLevel-1/shortener/internal/app/redirectBL"
//	"context"
//	"encoding/json"
//	"fmt"
//	"io"
//	"log"
//	"net/http"
//	"time"
//)
//
//type Router struct {
//	*http.ServeMux
//	redirectBL *redirectBL.Redirect
//}
//
//func NewRouter(redirectBL *redirectBL.Redirect) *Router {
//	r := &Router{
//		ServeMux:   http.NewServeMux(),
//		redirectBL: redirectBL,
//	}
//	r.HandleFunc("/create", r.CreateShortener)
//	r.HandleFunc("/{id}", r.Redirect)
//	r.HandleFunc("/stat/{id}", r.Statistic)
//	//r.Handle("/{uuid}", r.AuthMiddleware(http.HandlerFunc(r.RedirectAPI)))
//	//r.Handle("/delete", r.AuthMiddleware(http.HandlerFunc(r.DeleteUser)))
//	//r.HandleFunc("/search", r.AuthMiddleware(http.HandlerFunc(r.SearchUser)).ServeHTTP)
//	//r.Handle("/whoami", r.AuthMiddleware(http.HandlerFunc(r.Whoami)))
//	return r
//}
//
//func (rt *Router) AuthMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(
//		func(w http.ResponseWriter, r *http.Request) {
//			if u, p, ok := r.BasicAuth(); !ok || !(u == "admin" && p == "admin") {
//				fmt.Println("ok: ", ok, "user: ", u, "password", p)
//				http.Error(w, "unautorized", http.StatusUnauthorized)
//				return
//			}
//			// r = r.WithContext(context.WithValue(r.Context(), 1, 0))
//			next.ServeHTTP(w, r)
//		},
//	)
//}
//
////type Shortener struct {
////	ShortLink     string    `json:"short_link"`
////	FullLink      string    `json:"full_link"`
////	StatisticLink string    `json:"statistic_link"`
////	TotalCount    int       `json:"total_count"`
////	CreatedAt     time.Time `json:"created_at"`
////}
////
////type Statistic struct {
////	ShortLink  string                  `json:"short_link"`
////	TotalCount int                     `json:"total_count"`
////	CreatedAt  time.Time               `json:"created_at"`
////	FollowList []followingBL.Following `json:"follow_list"`
////}
//
//func (rt *Router) CreateShortener(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodPost {
//		fmt.Println(r.Method)
//		http.Error(w, "bad method", http.StatusMethodNotAllowed)
//		return
//	}
//	defer r.Body.Close()
//
//	shortener := Shortener{}
//	if err := json.NewDecoder(r.Body).Decode(&shortener); err != nil {
//		http.Error(w, "bad request", http.StatusBadRequest)
//		return
//	}
//
//	newShort, err := rt.redirectBL.CreateShortLink(r.Context(), []byte(shortener.FullLink))
//	if err != nil {
//		http.Error(w, "error when creating", http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusCreated)
//	_ = json.NewEncoder(w).Encode(
//		Shortener{
//			ShortLink:     newShort.ShortLink,
//			FullLink:      newShort.FullLink,
//			StatisticLink: newShort.StatisticLink,
//			TotalCount:    newShort.TotalCount,
//			CreatedAt:     newShort.CreatedAt,
//		},
//	)
//}
//
//func (rt *Router) Redirect(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodGet {
//		http.Error(w, "bad method", http.StatusMethodNotAllowed)
//		return
//	}
//	defer r.Body.Close()
//
//	//shortlink := Link{}
//	//if err := json.NewDecoder(r.Body).Decode(&shortlink); err != nil {
//	//	http.Error(w, "bad request", http.StatusBadRequest)
//	//	return
//	//}
//
//	shortlink := r.URL.String()
//
//	ctx := context.WithValue(r.Context(), "ipAddress", r.RemoteAddr)
//	getFullink, err := rt.redirectBL.GetFullLink(ctx, []byte(shortlink))
//	if err != nil {
//		http.Error(w, "error when get URL", http.StatusInternalServerError)
//		return
//	}
//
//	client := http.Client{Timeout: time.Second * 2}
//	req, _ := http.NewRequest("GET", getFullink, nil)
//	res, err := client.Do(req)
//	if err != nil {
//		http.Error(w, "error when creating", http.StatusInternalServerError)
//		return
//	}
//
//	if res.StatusCode > 399 {
//		http.Error(w, "error when creating", http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(res.StatusCode)
//	_, err = io.Copy(w, res.Body)
//	if err != nil {
//		log.Println("func Redirect error response write: ", err)
//	}
//}
//
//func (rt *Router) Statistic(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodGet {
//		fmt.Println(r.Method)
//		http.Error(w, "bad method", http.StatusMethodNotAllowed)
//		return
//	}
//	defer r.Body.Close()
//
//	//statisticlink := Link{}
//	//if err := json.NewDecoder(r.Body).Decode(&statisticlink); err != nil {
//	//	http.Error(w, "bad request", http.StatusBadRequest)
//	//	return
//	//}
//
//	statisticlink := r.URL.String()
//
//	statistic, err := rt.redirectBL.GetStatisticList(r.Context(), []byte(statisticlink))
//	if err != nil {
//		http.Error(w, "error when creating", http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusCreated)
//	_ = json.NewEncoder(w).Encode(statistic)
//}
