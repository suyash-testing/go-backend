package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/suyash-testing/go-backend/handlers"
	"github.com/suyash-testing/go-backend/routes"
)

func main() {
   // Initialize MongoDB connection and collection
   handlers.InitAuthCollection()

   r := mux.NewRouter()

   // Initialize authentication-related routes
   routes.AuthRoutes(r)

   log.Println("Server is running on port 8080")
   log.Fatal(http.ListenAndServe(":8080", r))
}
