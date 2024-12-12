package http

import (
	"AuthenticationService/internal/handler/http/middleware"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log/slog"
	"net/http"

	_ "AuthenticationService/docs"
)

func NewRouter(log *slog.Logger) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(middleware.SetRequestId)

	logger := middleware.NewLogger(log)
	router.Use(logger.LogRequest)

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	return router
}
