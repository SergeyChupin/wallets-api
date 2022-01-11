package httpserver

import (
	"log"
	"net/http"
	"os"

	"github.com/SergeyChupin/wallets-api/internal/app/httpserver/api"
	"github.com/SergeyChupin/wallets-api/internal/app/httpserver/config"
	"github.com/SergeyChupin/wallets-api/internal/database/postgres"
	"github.com/SergeyChupin/wallets-api/internal/repository"
	"github.com/SergeyChupin/wallets-api/internal/server"
	"github.com/SergeyChupin/wallets-api/internal/service"
)

func Run(configPath string) {
	logger := log.New(os.Stdout, "wallets-api ", log.LstdFlags)

	configLoader := config.NewLoader(configPath)
	cfg, err := configLoader.Load()
	if err != nil {
		logger.Fatal(err)
	}

	db, err := postgres.Open(logger, cfg.Postgres)
	if err != nil {
		logger.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	walletRepository := repository.NewWalletRepository(db)
	walletService := service.NewWalletService(walletRepository)

	handler := api.NewHandler(logger, walletService)
	srv := server.NewServer(logger, cfg.Server, handler)

	go func() {
		if err := srv.Start(); err != nil {
			if err != http.ErrServerClosed {
				logger.Fatal(err)
			}
		}
	}()

	if err := srv.GracefulShutdown(); err != nil {
		logger.Fatal(err)
	}
}
