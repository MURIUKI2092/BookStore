package routes

import (
	"BookStore/handlers"
	"BookStore/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(router *mux.Router) {
	router.HandleFunc("/users/add", handlers.CreateSingleUser).Methods("POST")
	router.HandleFunc("/auth/login", handlers.LoginUser).Methods("POST")
	// protected routes
	router.Handle("/auth/profile", middleware.Authenticate(http.HandlerFunc(handlers.GrabUserProfile))).Methods("GET")
	router.Handle("/users/single", middleware.Authenticate(http.HandlerFunc(handlers.GetSingleUser))).Methods("GET")
	router.Handle("/users/single/update", middleware.Authenticate(http.HandlerFunc(handlers.UpdateSingleUser))).Methods("PUT")
	router.Handle("/users/password/update", middleware.Authenticate(http.HandlerFunc(handlers.UpdateUserPassword))).Methods("PUT")
	// -------------------store-------------------------------------->
	router.Handle("/store/new", middleware.Authenticate(http.HandlerFunc(handlers.CreateSingleBookStore))).Methods("POST")
	router.Handle("/store", middleware.Authenticate(http.HandlerFunc(handlers.GetSingleBookStore))).Methods("GET")
	router.Handle("/store", middleware.Authenticate(http.HandlerFunc(handlers.UpdateSingleBookStore))).Methods("PUT")

	// -------------------book-------------------------------------->
	router.Handle("/book/new", middleware.Authenticate(http.HandlerFunc(handlers.CreatesSingleBook))).Methods("POST")
	router.Handle("/book/all", middleware.Authenticate(http.HandlerFunc(handlers.GetAllBooks))).Methods("GET")
	router.Handle("/book/single", middleware.Authenticate(http.HandlerFunc(handlers.GetSingleBook))).Methods("GET")
	router.Handle("/book/update", middleware.Authenticate(http.HandlerFunc(handlers.UpdateSingleBook))).Methods("PUT")
}
