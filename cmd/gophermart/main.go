package main

import (
	"net/http"

	"github.com/KaziPHone/go-musthave-diploma-tpl/internal/config"
	"github.com/KaziPHone/go-musthave-diploma-tpl/internal/handler"
	"github.com/KaziPHone/go-musthave-diploma-tpl/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	router := chi.NewRouter()
	h := handler.NewHandler(cfg.DataBaseDsn, []byte(cfg.SecretKey))

	router.Use(middleware.LoggingMiddleware)

	router.Post("/api/user/register", h.RegisterUserHandler)
	router.Post("/api/user/login", h.LoginHandler)

	router.Group(func(router chi.Router) {
		router.Use(middleware.AuthMiddleware(cfg.SecretKey))

		router.Post("/api/user/orders", h.UploadOrderHandler)
		router.Get("/api/user/orders", h.GetOrdersHandler)

		router.Get("/api/user/balance", h.GetBalanceHandler)

		router.Post("/api/user/balance/withdraw", h.WithdrawHandler)
		router.Get("/api/user/withdrawals", h.GetWithdrawalsHandler)
	})

	log.Printf("Starting server on: %s...", cfg.Host)
	err = http.ListenAndServe(cfg.Host, router)
	if err != nil {
		log.Err(err)
	}
}
