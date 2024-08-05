package handlers

import (
	"gophermart/internal/middleware"

	"github.com/go-chi/chi/v5"
)

type GmHandler struct {
	gmService    gmService
	signatureKey string
}

func New(gmService gmService, signatureKey string) (*GmHandler, error) {
	gmHandler := &GmHandler{
		gmService:    gmService,
		signatureKey: signatureKey,
	}

	return gmHandler, nil
}

func (h *GmHandler) InitRouter() chi.Router {

	r := chi.NewRouter()
	r.Use(middleware.WithLogging, middleware.WithGzip)

	r.Route("/api/user", func(r chi.Router) {

		// Вложенный маршрут с промежуточным обработчиком WithMakeAuth для /register и /login
		r.Route("/", func(r chi.Router) {
			r.Use(middleware.WithMakeAuth(h.signatureKey))

			r.Post("/register", h.register())
			r.Post("/login", h.login())
		})

		// Вложенный маршрут для /orders с промежуточным обработчиком CheckAuth
		r.Route("/orders", func(r chi.Router) {
			r.Use(middleware.WithCheckAuth(h.signatureKey))

			r.Post("/", h.addOrder())
			r.Get("/", h.getOrders())
		})

		// Вложенный маршрут для /balance с промежуточным обработчиком CheckAuth
		r.Route("/balance", func(r chi.Router) {
			r.Use(middleware.WithCheckAuth(h.signatureKey))

			r.Get("/", h.getBalance())
			r.Post("/withdraw", h.withdraw())
		})

		// Вложенный маршрут для /withdrawals с промежуточным обработчиком CheckAuth
		r.Route("/withdrawals", func(r chi.Router) {
			r.Use(middleware.WithCheckAuth(h.signatureKey))

			r.Get("/", h.getWithdrawals())
		})

	})

	return r
}
