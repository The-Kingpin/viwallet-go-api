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

	mux.Get("/", handlers.Home)

	mux.Get("/user/login", handlers.Login)
	mux.Post("/user/login", handlers.PostLogin)

	mux.Get("/user/register", handlers.Register)
	mux.Post("/user/register", handlers.PostRegister)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
