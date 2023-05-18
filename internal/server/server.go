package server

import (
	"net/http"

	"github.com/NevostruevK/GopherMart.git/internal/db"
	"github.com/NevostruevK/GopherMart.git/internal/server/handlers"
	"github.com/NevostruevK/GopherMart.git/internal/server/middleware"
	"github.com/NevostruevK/GopherMart.git/internal/util/token"
	"github.com/go-chi/chi/v5"
)

func NewServer(db *db.DB, address string) (*http.Server, error) {
	tk, err := token.NewToken()
	if err != nil {
		return nil, err
	}
	r := chi.NewRouter()

	handler := middleware.GzipCompressMiddleware(r)
	handler = middleware.GzipDecompressMiddleware(handler)
	handler = middleware.AuthMiddleware(handler, tk)
	handler = middleware.CheckHeadersMiddleware(handler)
	handler = middleware.LoggerMiddleware(handler)

	r.Post("/api/user/register", handlers.Authentication(db, tk, handlers.Register))
	r.Post("/api/user/login", handlers.Authentication(db, tk, handlers.Login))
	r.Post("/api/user/orders", handlers.PostOrder(db))
	r.Post("/api/user/balance/withdraw", handlers.PostWithdraw(db))
	r.Get("/api/user/orders", handlers.GetOrders(db))
	r.Get("/api/user/balance", handlers.GetBalance(db))
	r.Get("/api/user/withdrawals", handlers.GetWithdrawals(db))

	return &http.Server{
		Addr:    address,
		Handler: handler,
	}, nil
}
