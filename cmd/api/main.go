package main

import (
	"net/http"
	"restapi-go/config"
	"restapi-go/routes"
)

func main() {
	config.ConnectDB()

	router := routes.RegisterRoutes()
	http.ListenAndServe(":8080", router)
}
