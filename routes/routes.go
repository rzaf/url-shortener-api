package routes

import (
	"github.com/rzaf/url-shortener-api/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func GetRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)
	router.Use(handlers.RecoverServerPanics)

	router.Group(func(r chi.Router) {
		r.Use(handlers.AuthApiKeyMiddleware)

		r.Get("/urls", handlers.GetUrls)
		r.Post("/urls", handlers.CreateUrl)
		r.Delete("/urls/{short}", handlers.DeleteUrl)
		r.Put("/urls/{short}", handlers.EditUrl)

		r.Get("/users/{id}", handlers.GetUserById)
		r.Get("/users", handlers.GetUsers)
		r.Delete("/users/{id}", handlers.DeleteUser)
		r.Post("/users/{id}/newApiKey", handlers.EditUserApiKey)
		r.Put("/users/{id}", handlers.EditUser)
	})

	router.Get("/urls/{short}", handlers.GetUrl)
	router.Post("/users", handlers.CreateUser)

	baseRouter := chi.NewRouter()
	baseRouter.Mount("/api/", router)
	return baseRouter
}
