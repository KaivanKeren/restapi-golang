package routes

import (
	"net/http"
	"restapi-go/controllers"
	"restapi-go/middleware"

	"github.com/go-chi/chi"
)

func RegisterRoutes() *chi.Mux {
	router := chi.NewRouter()

	// User Routes
	
	router.With(middleware.AuthMiddleware).Get("/users", controllers.GetUsers)
	router.With(middleware.AuthMiddleware).Post("/users", controllers.CreateUser)
	router.With(middleware.AuthMiddleware).Get("/users/{id}", controllers.GetUser)
	router.With(middleware.AuthMiddleware).Put("/users/{id}", controllers.UpdateUser)
	router.With(middleware.AuthMiddleware).Delete("/users/{id}", controllers.DeleteUser)

	// Auth Routes
	router.Post("/login", controllers.Login)
	router.With(middleware.AuthMiddleware).Post("/logout", controllers.Logout)

	http.ListenAndServe(":8080", middleware.CORSMiddleware(router))

	return router
}
