package routerchi

import (
	"fmt"
	"net/http"
	"shop/internal/api/handlers"
	"shop/internal/app/itemBL"
	"strconv"

	"github.com/go-chi/render"
	"github.com/google/uuid"

	"github.com/go-chi/chi/v5"
)

type RouterChi struct {
	*chi.Mux
	handl *handlers.Handlers
}

func NewRouterChi(hs *handlers.Handlers) *RouterChi {
	chirouter := chi.NewRouter()
	r := &RouterChi{
		handl: hs,
	}

	chirouter.Group(func(hr chi.Router) {
		hr.Post("/items", r.ListItems)           // "POST"
		hr.Put("/item", r.CreateItem)            // "PUT"
		hr.Post("/item/{id}", r.GetItem)         // "POST"
		hr.Put("/item/{id}", r.UpdateItem)       // "PUT"
		hr.Delete("/item/{id}", r.DeleteItem)    // "DELETE"
		hr.Post("/search/{name}", r.SearchItems) // "POST"
	})

	r.Mux = chirouter
	return r
}

type Item handlers.ItemHendler

func (Item) Bind(r *http.Request) error {
	return nil
}

func (Item) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (rt *RouterChi) CreateItem(w http.ResponseWriter, r *http.Request) {
	rItem := Item{}
	if err := render.Bind(r, &rItem); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	newItem, err := rt.handl.CreateItemHandler(r.Context(), handlers.ItemHendler(rItem))
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Render(w, r, Item(newItem))
}

func (rt *RouterChi) GetItem(w http.ResponseWriter, r *http.Request) {
	sid := chi.URLParam(r, "id")

	uid, err := uuid.Parse(sid)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	getItem, err := rt.handl.GetItemHandler(r.Context(), uid)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Render(w, r, Item(*getItem))
}

func (rt *RouterChi) UpdateItem(w http.ResponseWriter, r *http.Request) {
	sid := chi.URLParam(r, "id")

	uid, err := uuid.Parse(sid)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	rItem := Item{}
	if err := render.Bind(r, &rItem); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	rItem.ID = uid

	updItem, err := rt.handl.UpdateItemHandler(r.Context(), handlers.ItemHendler(rItem))
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Render(w, r, Item(*updItem))
}

func (rt *RouterChi) DeleteItem(w http.ResponseWriter, r *http.Request) {
	sid := chi.URLParam(r, "id")

	uid, err := uuid.Parse(sid)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	delItem, err := rt.handl.DeleteItemHandler(r.Context(), uid)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	//render.Data(w, r, []byte("item deleted successfully"))
	render.Render(w, r, Item(*delItem))
}

func (rt *RouterChi) parseItemFilterQuery(r *http.Request) itemBL.ItemFilter {
	filter := itemBL.ItemFilter{}

	if limitRaw := r.FormValue("limit"); limitRaw != "" {
		if limitInput, err := strconv.Atoi(limitRaw); err == nil {
			filter.Limit = limitInput
		}
	}
	if filter.Limit == 0 {
		filter.Limit = 5
	}

	if offsetRaw := r.FormValue("offset"); offsetRaw != "" {
		if offsetInput, err := strconv.Atoi(offsetRaw); err == nil {
			filter.Offset = offsetInput
		}
	}

	if priceRightRaw := r.FormValue("price_right"); priceRightRaw != "" {
		if priceRightInput, err := strconv.ParseInt(priceRightRaw, 10, 64); err == nil {
			filter.PriceRight = &priceRightInput
		}
	}

	if priceLeftRaw := r.FormValue("price_left"); priceLeftRaw != "" {
		if priceLeftInput, err := strconv.ParseInt(priceLeftRaw, 10, 64); err == nil {
			filter.PriceLeft = &priceLeftInput
		}
	}
	return filter
}

type ListItemResponse struct {
	Payload []*itemBL.ItemBL `json:"payload"`
	Limit   int              `json:"limit"`
	Offset  int              `json:"offset"`
}

func (rt *RouterChi) ListItems(w http.ResponseWriter, r *http.Request) {
	filter := rt.parseItemFilterQuery(r)

	fmt.Fprintln(w, "[")
	comma := false
	err := rt.handl.ListItemHandler(r.Context(), filter, func(item handlers.ItemHendler) error {
		if comma {
			fmt.Fprintln(w, ",")
		} else {
			comma = true
		}
		if err := render.Render(w, r, Item(item)); err != nil {
			return err
		}
		w.(http.Flusher).Flush()
		return nil
	})

	if err != nil {
		if comma {
			fmt.Fprint(w, ",")
		}
		render.Render(w, r, ErrRender(err))
	}

	fmt.Fprintln(w, "]")
}

func (rt *RouterChi) SearchItems(w http.ResponseWriter, r *http.Request) {
	s := chi.URLParam(r, "name")

	fmt.Fprintln(w, "[")
	comma := false
	err := rt.handl.SearchItemsHandler(r.Context(), s, func(item handlers.ItemHendler) error {
		if comma {
			fmt.Fprintln(w, ",")
		} else {
			comma = true
		}
		if err := render.Render(w, r, Item(item)); err != nil {
			return err
		}
		w.(http.Flusher).Flush()
		return nil
	})

	if err != nil {
		if comma {
			fmt.Fprint(w, ",")
		}
		render.Render(w, r, ErrRender(err))
	}

	fmt.Fprintln(w, "]")
}
