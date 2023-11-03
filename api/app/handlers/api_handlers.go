package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/fredbi/go-api-skeleton/api/pkg/injected"
	"github.com/fredbi/go-api-skeleton/api/pkg/repos"
	"github.com/fredbi/go-trace/log"
	"github.com/fredbi/go-trace/tracer"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const jsonMime = "application/json"

// Register a CRUD json API for sample items
func Register(rt injected.Iface) {
	h := New(rt)

	rt.App().Route("/sample", func(r chi.Router) {
		r.Get("/items", h.ListItems)
		r.Get("/item/{id}", h.GetItem)
		r.Put("/item/{id}", h.UpdateItem)
		r.Post("/item", h.CreateItem)
		r.Delete("/item/{id}", h.DeleteItem)
	})
}

// Handler knows how to expose the sample API REST endpoints.
type Handler struct {
	rt     injected.Iface
	logger log.Factory
}

func New(rt injected.Iface) *Handler {
	return &Handler{
		rt:     rt,
		logger: rt.Logger(),
	}
}

func (h Handler) Logger() log.Factory {
	return h.logger
}

func (h Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	ctx, span, lg := tracer.StartSpan(r.Context(), h)
	defer span.End()

	id, httpErr := h.requireResourceFromPathParam(r, "id")
	if httpErr != nil {
		httpErr.From(ErrGetItem).WriteJSON(w)

		return
	}

	lg.Debug("get item", zap.String("id", id))

	item, err := h.rt.Repos().Sample().Get(ctx, id)
	if err != nil {
		checkErr(w, ErrGetItem, err)

		return
	}

	if err = json.NewEncoder(w).Encode(item); err != nil {
		checkErr(w, ErrGetItem, err)

		return
	}
}

func (h Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	ctx, span, lg := tracer.StartSpan(r.Context(), h)
	defer span.End()

	id, httpErr := h.requireResourceFromPathParam(r, "id")
	if httpErr != nil {
		httpErr.From(ErrUpdateItem).WriteJSON(w)

		return
	}

	item, err := h.parseItemFromJSON(w, r, ErrUpdateItem)
	if err != nil {
		return
	}

	item.ID = id
	lg.Debug("update item", zap.Any("item", item))
	// TODO(fredbi): check for required fields

	if err := h.rt.Repos().Sample().Update(ctx, item); err != nil {
		checkErr(w, ErrUpdateItem, err)

		return
	}
}

func (h Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	ctx, span, lg := tracer.StartSpan(r.Context(), h)
	defer span.End()

	item, err := h.parseItemFromJSON(w, r, ErrCreateItem)
	if err != nil {
		return
	}
	item.ID = ""
	// TODO(fredbi): check for required fields

	lg.Debug("create item", zap.Any("item", item))

	id, err := h.rt.Repos().Sample().Create(ctx, item)
	if err != nil {
		checkErr(w, ErrCreateItem, err)

		return
	}

	h.putResourceInHeader(w, fmt.Sprintf("/sample/item/%s", id))
}

func (h Handler) putResourceInHeader(w http.ResponseWriter, resource string) {
	w.Header().Set("Content-Type", fmt.Sprintf("%s;charset=utf-8", jsonMime))

	// identifies the new resource as a response header, with its full path
	w.Header().Set("ID", resource)
}

func (h Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	ctx, span, lg := tracer.StartSpan(r.Context(), h)
	defer span.End()

	id, err := h.requireResourceFromPathParam(r, "id")
	if err != nil {
		err.From(ErrDeleteItem).WriteJSON(w)
		return
	}

	lg.Debug("delete item", zap.String("id", id))

	if err := h.rt.Repos().Sample().Delete(ctx, id); err != nil {
		checkErr(w, ErrDeleteItem, err)

		return
	}
}

func (h Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	ctx, span, lg := tracer.StartSpan(r.Context(), h)
	defer span.End()

	lg.Debug("list all items")

	// TODO(fredbi): feature - support filters, support pagination
	iterator, err := h.rt.Repos().Sample().List(ctx)
	if err != nil {
		checkErr(w, ErrListItems, err)

		return
	}
	defer func() {
		_ = iterator.Close()
	}()

	items := make([]repos.Item, 0, 100)
	for iterator.Next() {
		item, errItem := iterator.Item()
		if errItem != nil {
			checkErr(w, ErrListItems, errItem)

			return
		}
		items = append(items, item)
	}

	// TODO(fredbi): performance - encode chunks and push back partial responses
	// while iterating.
	if err = json.NewEncoder(w).Encode(items); err != nil {
		checkErr(w, ErrGetItem, err)

		return
	}
}

func (h Handler) requireResourceFromPathParam(r *http.Request, param string) (string, *Error) {
	id := chi.URLParam(r, param)
	if id != "" {
		return id, nil
	}

	return "", RequiredErr(param)
}

func (h Handler) parseItemFromJSON(w http.ResponseWriter, r *http.Request, base error) (repos.Item, error) {
	var item repos.Item

	if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, jsonMime) {
		err := errors.New("Content-Type must be application/json")
		checkErr(w, base, err)

		return item, err
	}

	// TODO(fredbi): performance - don't use encoding/json in real conditions (too slow)
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		checkErr(w, base, err)

		return item, err
	}

	return item, nil
}
