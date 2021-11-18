package chiRouter

import (
	"CourseProjectBackendDevGoLevel-1/shortener/internal/api/handler"
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"strings"

)

type ChiRouter struct {
	*chi.Mux
	hs *handler.Handlers
}

func NewChiRouter(handlers *handler.Handlers) *ChiRouter {
	chiNew := chi.NewRouter()

	chiR := &ChiRouter{
		hs: handlers,
	}

	chiNew.Group(func(r chi.Router) {
		r.Post("/create", chiR.CreateShortener)
		r.Get("/{short}", chiR.Redirect)
		r.Post("/stat", chiR.Statistic)
	})

	chiR.Mux = chiNew

	return chiR
}

type Shortener handler.Shortener
func (Shortener) Bind(r *http.Request) error {
	return nil
}
func (Shortener) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (chr *ChiRouter) CreateShortener(w http.ResponseWriter, r *http.Request) {
	rShortener := Shortener{}
	if err := render.Bind(r, &rShortener); err != nil {
		//nolint:errcheck
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	newShort, err := chr.hs.CreateShortener(r.Context(), handler.Shortener(rShortener))
	if err != nil {
		log.Println(err)
		return
	}

	err = render.Render(w, r, Shortener(newShort))
	if err != nil {
		log.Println(err)
	}
}

func (chr *ChiRouter) Redirect(w http.ResponseWriter, r *http.Request) {
	rShortener := Shortener{}

	rShortener.ShortLink = chi.URLParam(r, "short")

	ipaddr := strings.Split(r.RemoteAddr, ":")
	//nolint:staticcheck
	ctx := context.WithValue(r.Context(), "IP_address", ipaddr[0])

	getFullink, err := chr.hs.Redirect(ctx, handler.Shortener(rShortener))
	if err != nil {
		//nolint:errcheck
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	http.Redirect(w, r, getFullink.FullLink, http.StatusFound)
}

type Statistic handler.Statistic
func (Statistic) Bind(r *http.Request) error {
	return nil
}
func (Statistic) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (chr *ChiRouter) Statistic(w http.ResponseWriter, r *http.Request) {
	rShortener := Shortener{}
	if err := render.Bind(r, &rShortener); err != nil {
		//nolint:errcheck
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	statistic, err := chr.hs.GetStatisticList(r.Context(), rShortener.StatLink)
	if err != nil {
		//nolint:errcheck
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	//nolint:errcheck
	render.Render(w, r, Statistic(statistic))
}
