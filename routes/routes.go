package routes

import (
	"restapi-go/controllers"

	"github.com/go-chi/chi"
)

func RegisterRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/users", controllers.GetUsers)
	router.Post("/users", controllers.CreateUser)
	router.Get("/users/{id}", controllers.GetUser)
	router.Put("/users/{id}", controllers.UpdateUser)
	router.Delete("/users/{id}", controllers.DeleteUser)

	return router
}
