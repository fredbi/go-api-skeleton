package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/fredbi/go-api-skeleton/api/pkg/injected"
	"github.com/fredbi/go-api-skeleton/api/pkg/repos"
	"github.com/fredbi/go-trace/log"
	"github.com/fredbi/go-trace/tracer"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

const jsonMime = "application/json"

// Register a CRUD json API for sample items
func Register(rt injected.Iface) {
	h := New(rt)

	rt.App().Route("/sample", func(r chi.Router) {
		r.Get("/item/{id}", http.HandlerFunc(h.GetItem))
		r.Put("/item/{id}", http.HandlerFunc(h.UpdateItem))
		r.Post("/item", http.HandlerFunc(h.CreateItem))
		r.Delete("/item/{id}", http.HandlerFunc(h.DeleteItem))
	})
}

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

	id, err := h.getResourceFromPathParam(w, r, "id", ErrGetItem)
	if err != nil {
		return
	}

	lg.Debug("get item", zap.String("id", id))

	item, err := h.rt.Repos().Sample().Get(ctx, id)
	if err != nil {
		h.checkErr(w, ErrGetItem, err)

		return
	}

	if err = json.NewEncoder(w).Encode(item); err != nil {
		h.checkErr(w, ErrGetItem, err)

		return
	}
}

func (h Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	ctx, span, lg := tracer.StartSpan(r.Context(), h)
	defer span.End()

	id, err := h.getResourceFromPathParam(w, r, "id", ErrUpdateItem)
	if err != nil {
		return
	}

	item, err := h.parseItemFromJSON(w, r, ErrUpdateItem)
	if err != nil {
		return
	}

	item.ID = id
	lg.Debug("update item", zap.Any("item", item))

	if err := h.rt.Repos().Sample().Update(ctx, item); err != nil {
		h.checkErr(w, ErrUpdateItem, err)

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

	lg.Debug("create item", zap.Any("item", item))

	id, err := h.rt.Repos().Sample().Create(ctx, item)
	if err != nil {
		h.checkErr(w, ErrCreateItem, err)

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

	id, err := h.getResourceFromPathParam(w, r, "id", ErrDeleteItem)
	if err != nil {
		return
	}

	lg.Debug("delete item", zap.String("id", id))

	if err := h.rt.Repos().Sample().Delete(ctx, id); err != nil {
		h.checkErr(w, ErrDeleteItem, err)

		return
	}
}

func (h Handler) getResourceFromPathParam(w http.ResponseWriter, r *http.Request, param string, base error) (string, error) {
	id := chi.URLParam(r, param)
	if err := h.requireID(w, id, base); err != nil {
		return "", err
	}

	return id, nil
}

func (h Handler) requireID(w http.ResponseWriter, id string, base error) error {
	if id != "" {
		return nil
	}

	err := fmt.Errorf("{id} is required")
	h.Error(w, errors.Join(base, err), http.StatusBadRequest)

	return err
}

func (h Handler) checkErr(w http.ResponseWriter, base, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		h.Error(w, errors.Join(base, err), http.StatusNotFound)

		return
	}

	h.Error(w, errors.Join(base, err), http.StatusInternalServerError)
}

// Error as a json response
func (h Handler) Error(w http.ResponseWriter, err error, code int) {
	msg := struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}{
		Message: err.Error(),
		Code:    code,
	}

	w.Header().Set("Content-Type", fmt.Sprintf("%s;charset=utf-8", jsonMime))
	_ = json.NewEncoder(w).Encode(msg)

	w.WriteHeader(code)
}

var (
	ErrGetItem    = errors.New("get item failed")
	ErrCreateItem = errors.New("create item failed")
	ErrUpdateItem = errors.New("update item failed")
	ErrDeleteItem = errors.New("delete item failed")
)

func (h Handler) parseItemFromJSON(w http.ResponseWriter, r *http.Request, base error) (repos.Item, error) {
	var item repos.Item

	if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, jsonMime) {
		err := errors.New("Content-Type must be application/json")
		h.checkErr(w, base, err)

		return item, err
	}

	// TODO: security - use ReadLimiter like in net/http
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.checkErr(w, base, err)

		return item, err
	}

	if len(body) == 0 {
		err = errors.New("empty body")
		h.checkErr(w, base, err)

		return item, err
	}

	if err = json.Unmarshal(body, &item); err != nil {
		h.checkErr(w, base, err)

		return item, err
	}

	return item, nil
}
