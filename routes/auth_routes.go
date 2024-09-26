package routes

import (
	"github.com/suyash-testing/go-backend/handlers"

	"github.com/gorilla/mux"
)

// AuthRoutes initializes the authentication-related routes.
func AuthRoutes(router *mux.Router) {
    // Signup and login routes
    router.HandleFunc("/signup", handlers.SignupHandler).Methods("POST")
    router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
}