package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type server struct {
	logger     *log.Logger
	httpServer *http.Server
	config     Config
}

func NewServer(logger *log.Logger, config Config, handler http.Handler) *server {
	server := &server{
		logger: logger,
		config: config,
		httpServer: &http.Server{
			Addr:         ":" + config.Port,
			Handler:      handler,
			ErrorLog:     logger,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
		},
	}
	return server
}

func (server *server) Start() error {
	server.logger.Println("Starting server on port " + server.config.Port)
	return server.httpServer.ListenAndServe()
}

func (server *server) GracefulShutdown() error {
	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)
	<-sigChannel
	server.logger.Println("Shutdown server")
	ctx, cancel := context.WithTimeout(context.Background(), server.config.ShutdownGracePeriod)
	defer cancel()
	return server.httpServer.Shutdown(ctx)
}
