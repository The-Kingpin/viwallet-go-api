package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gitlab.com/code-harbor/viwallet/internal/handlers"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(CSRFProtect)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)

	mux.Get("/user/login", handlers.Repo.Login)
	mux.Post("/user/login", handlers.Repo.PostLogin)

	mux.Get("/user/register", handlers.Repo.Register)
	mux.Post("/user/register", handlers.Repo.PostRegister)

	mux.Get("/transactions", handlers.Repo.Trasnsactions)
	mux.Get("/cards", handlers.Repo.Cards)
	mux.Get("/wallets", handlers.Repo.Wallets)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
