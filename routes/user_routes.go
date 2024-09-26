package routes

import (
	"github.com/gorilla/mux"
	"github.com/suyash-testing/go-backend/handlers"
)

// UserRoutes initializes all user-related routes
func UserRoutes(router *mux.Router) {
    // Define routes for CRUD operations
    router.HandleFunc("/user", handlers.CreateUser).Methods("POST")
    router.HandleFunc("/users", handlers.GetUsers).Methods("GET")
    router.HandleFunc("/user", handlers.UpdateUser).Methods("PUT")
    router.HandleFunc("/user", handlers.DeleteUser).Methods("DELETE")
}