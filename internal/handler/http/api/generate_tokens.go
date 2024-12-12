package api

import (
	httpApi "AuthenticationService/internal/handler/http"
	"AuthenticationService/internal/handler/http/middleware"
	"fmt"
	"log/slog"
	"net/http"
)

// GenerateTokens - generates tokens
// @Tags auth
// @Description generate tokens
// @Accept  json
// @Produce  json
// @Param	id query string true "id for generate tokens" Format(uuid)
// @Success 200 {object} model.Tokens
// @Failure 400,404,409 {object} model.ResponseError
// @Failure 500 {object} model.InternalError
// @Router /auth/token [get]
func (u *AuthController) GenerateTokens(w http.ResponseWriter, r *http.Request) {
	const op = "handler.http.api.GenerateTokens"

	log := u.log.With(
		slog.String("op", op),
		slog.String("request_id", r.Context().Value(middleware.CtxKeyRequestID).(string)),
	)

	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		log.Error("failed to decode request")
		ErrorResponseWithCode(w, r, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}

	access, refresh, err := u.service.GenerateTokens(id, r.RemoteAddr)
	if err != nil {
		log.Error("failed to generate tokens", slog.String("error", err.Error()))
		ErrorResponse(w, r, err)
		return
	}

	res := httpApi.ToTokenApiModel(access, refresh)

	response(w, r, http.StatusOK, res)
}
