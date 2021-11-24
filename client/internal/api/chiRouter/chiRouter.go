package chiRouter

import (
	"DZ_Backend_dev_Go_level_2/client/internal/api/handlers"
	"DZ_Backend_dev_Go_level_2/client/internal/api/recoverMW"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strings"
)

type ChiRouter struct {
	*chi.Mux
}

func NewChiRouter() *ChiRouter {
	chiNew := chi.NewRouter()

	chiR := &ChiRouter{}

	chiNew.Group(func(r chi.Router) {
		r.Use(recoverMW.RecoverMiddleware)

		r.Get("/", chiR.GetHomePage)
		r.Post("/", chiR.PostHomePage)
		r.Get("/{short}", chiR.RedirectPage)
		r.Get("/stat/{stat}", chiR.StatPage)
	})

	chiNew.Get("/__heartbeat__", func(w http.ResponseWriter, r *http.Request) {})

	chiR.Mux = chiNew

	return chiR
}

func (chr *ChiRouter) GetHomePage(w http.ResponseWriter, r *http.Request) {

	b, err := handlers.GetHomePage()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	_, err = b.WriteTo(w)
	if err != nil {
		log.Println("error rendering home page: ", err)
	}
}

func (chr *ChiRouter) PostHomePage(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}

	fullLink := &handlers.FullLink{
		FullLink: r.PostFormValue("fullLink"),
	}

	strJSON, err := json.Marshal(&fullLink)
	if err != nil {
		fmt.Fprintf(w, "A error occured json.NewEncoder(&b).Encode(p).")
	}

	b, err := handlers.PostHomePage(strJSON)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	_, err = b.WriteTo(w)
	if err != nil {
		log.Println("error rendering home page: ", err)
	}
}

func (chr *ChiRouter) RedirectPage(w http.ResponseWriter, r *http.Request) {
	ipaddr := strings.Split(r.RemoteAddr, ":")
	//nolint:staticcheck
	//r = r.WithContext(context.WithValue(r.Context(), "IP_address", ipaddr[0]))

	shortLink := &handlers.Redirect{
		ShortLink: chi.URLParam(r, "short"),
		IPaddress: ipaddr[0],
	}

	respRedirect, err := handlers.RedirectPage(shortLink)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	http.Redirect(w, r, respRedirect.FullLink, http.StatusFound)

}

func (chr *ChiRouter) StatPage(w http.ResponseWriter, r *http.Request) {
	statLink := &handlers.StatLink{
		StatLink: chi.URLParam(r, "stat"),
	}

	strJSON, err := json.Marshal(&statLink)
	if err != nil {
		fmt.Fprintf(w, "func StatPage: error occured json marshaling stat page")
	}

	b, err := handlers.StatPage(strJSON)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	_, err = b.WriteTo(w)
	if err != nil {
		log.Println("error rendering home page: ", err)
	}
}
