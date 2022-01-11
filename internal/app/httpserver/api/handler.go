package api

import (
	"log"
	"net/http"

	"github.com/SergeyChupin/wallets-api/internal/app/httpserver/api/v1"
	"github.com/SergeyChupin/wallets-api/internal/service"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

type handler struct {
	logger *log.Logger
	router *mux.Router
}

func NewHandler(logger *log.Logger, walletService service.WalletService) *handler {
	handler := &handler{
		logger: logger,
	}
	handler.initRoutes(walletService)
	return handler
}

func (handler *handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	handler.router.ServeHTTP(rw, req)
}

func (handler *handler) initRoutes(walletService service.WalletService) {
	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	v1.NewWalletsApi(handler.logger, apiRouter, walletService)

	redocOpts := middleware.RedocOpts{SpecURL: "/api.yaml"}
	redocHandler := middleware.Redoc(redocOpts, nil)

	router.Handle("/docs", redocHandler)
	router.Handle("/api.yaml", http.FileServer(http.Dir("./docs")))

	handler.router = router
}
