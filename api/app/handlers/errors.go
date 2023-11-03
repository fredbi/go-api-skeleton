package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Errors to wrap the context of a failure
var (
	ErrGetItem    = errors.New("get item failed")
	ErrCreateItem = errors.New("create item failed")
	ErrUpdateItem = errors.New("update item failed")
	ErrDeleteItem = errors.New("delete item failed")
	ErrListItems  = errors.New("list items failed")
)

type Error struct {
	Code int
	Err  error
	_    struct{}
}

func (e Error) MarshalJSON() ([]byte, error) {
	stage := struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}{
		Message: e.Err.Error(),
		Code:    e.Code,
	}

	return json.Marshal(stage)
}

func (e Error) Error() string {
	return fmt.Sprintf("%s:%d", e.Err.Error(), e.Code)
}

// From an origin sentinel error
func (e Error) From(err error) *Error {
	return &Error{
		Code: e.Code,
		Err:  errors.Join(err, e),
	}
}

// WriteJSON sends an error as a json response
func (e Error) WriteJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", fmt.Sprintf("%s;charset=utf-8", jsonMime))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_ = json.NewEncoder(w).Encode(e)

	w.WriteHeader(e.Code)
}

func RequiredErr(param string) *Error {
	return &Error{
		Code: http.StatusBadRequest,
		Err:  fmt.Errorf("{%s} is required", param),
	}
}

func OtherErr() *Error {
	return &Error{
		Code: http.StatusInternalServerError,
		Err:  errors.New("internal Server Error"),
	}
}

func checkErr(w http.ResponseWriter, base, err error) {
	httpErr := Error{
		Err: err,
	}.From(base)

	if errors.Is(err, sql.ErrNoRows) {
		httpErr.Code = http.StatusNotFound
	} else {
		httpErr.Code = http.StatusInternalServerError
	}

	httpErr.WriteJSON(w)
}
