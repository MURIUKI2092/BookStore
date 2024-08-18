package main

import (
	"BookStore/database"
	"BookStore/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// initialize database connection
	database.Connect()

	// create a new router
	router := mux.NewRouter()

	// register user routes
	routes.RegisterUserRoutes(router)

	// Start the server
	log.Println("Starting server on :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
