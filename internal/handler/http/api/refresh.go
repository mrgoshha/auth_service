package api

import (
	httpApi "AuthenticationService/internal/handler/http"
	"AuthenticationService/internal/handler/http/middleware"
	apimodel "AuthenticationService/internal/handler/http/model"
	"encoding/json"
	"log/slog"
	"net/http"
)

// Refresh - refresh token
// @Tags auth
// @Description refresh tokens
// @Accept  json
// @Produce  json
// @Param	tokens body model.Tokens true "tokens for refresh"
// @Success 200 {object} model.Tokens
// @Failure 400,404,409 {object} model.ResponseError
// @Failure 500 {object} model.InternalError
// @Router /auth/refresh [post]
func (u *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	const op = "handler.http.api.Refresh"

	log := u.log.With(
		slog.String("op", op),
		slog.String("request_id", r.Context().Value(middleware.CtxKeyRequestID).(string)),
	)

	req := &apimodel.Tokens{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Error("failed to decode request")
		ErrorResponseWithCode(w, r, http.StatusBadRequest, err)
		return
	}

	access, refresh, err := u.service.Refresh(req.AccessToken, req.RefreshToken, r.RemoteAddr)
	if err != nil {
		log.Error("failed to refresh tokens", slog.String("error", err.Error()))
		ErrorResponse(w, r, err)
		return
	}

	res := httpApi.ToTokenApiModel(access, refresh)

	response(w, r, http.StatusOK, res)
}
