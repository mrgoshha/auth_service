package api

import (
	"AuthenticationService/internal/adapter/dbs"
	"AuthenticationService/internal/handler/http/middleware"
	"AuthenticationService/internal/handler/http/model"
	"encoding/json"
	"errors"
	"net/http"
)

func ErrorResponseWithCode(w http.ResponseWriter, r *http.Request, code int, err error) {
	var modelErr interface{}
	if code == http.StatusInternalServerError {
		modelErr = model.InternalError{
			RequestId: r.Context().Value(middleware.CtxKeyRequestID).(string),
		}
	} else {
		modelErr = model.ResponseError{
			Message:   err.Error(),
			RequestId: r.Context().Value(middleware.CtxKeyRequestID).(string),
		}
	}
	response(w, r, code, modelErr)

}

func ErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	var code int
	switch {
	case errors.Is(err, dbs.ErrorRecordNotFound):
		code = http.StatusNotFound
	case errors.Is(err, dbs.ErrorRecordAlreadyExists):
		code = http.StatusConflict
	default:
		code = http.StatusInternalServerError
	}
	ErrorResponseWithCode(w, r, code, err)
}

func response(w http.ResponseWriter, _ *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
