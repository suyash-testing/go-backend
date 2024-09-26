package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/suyash-testing/go-backend/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var authCollection *mongo.Collection

// Initialize the auth collection
func InitAuthCollection() {
    client := db.ConnectMongo()
    authCollection = db.GetCollection(client, "users")
}

// User structure for authentication
type User struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// HashPassword hashes a plain password using bcrypt.
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// CheckPasswordHash checks if the provided password matches the stored hash.
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// SignupHandler handles user signup by creating a new user with a hashed password.
func SignupHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	
	// Hash the user's password before storing it in the database
    hashedPassword, err := HashPassword(user.Password)
    if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }
    user.Password = hashedPassword
	
    // Insert the user into the database
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	_, err = authCollection.InsertOne(ctx, user)
    if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
        return
    }
	
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("User created successfully"))
}

// LoginHandler handles user login by checking if the provided credentials are valid.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Find the user in the database by email
    var foundUser User
    err = authCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
    if err != nil {
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    // Check if the provided password matches the hashed password
    if !CheckPasswordHash(user.Password, foundUser.Password) {
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    w.Write([]byte("Login successful"))
}