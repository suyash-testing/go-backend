package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/suyash-testing/go-backend/db"
)

var userCollection *mongo.Collection

func InitUserCollection() {
    client := db.ConnectMongo()
    userCollection = db.GetCollection(client, "users")
}

// User structure
type UserTesting struct {
    ID       string `json:"id,omitempty" bson:"_id,omitempty"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

// CreateUser inserts a new user into the database.
func CreateUser(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := userCollection.InsertOne(ctx, user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "User created with ID: %v", result.InsertedID)
}

// GetUsers retrieves all users from the database.
func GetUsers(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := userCollection.Find(ctx, bson.M{})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var users []User
    for cursor.Next(ctx) {
        var user User
        cursor.Decode(&user)
        users = append(users, user)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

// UpdateUser updates an existing user's email or password.
func UpdateUser(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    filter := bson.M{"email": user.Email}
    update := bson.M{
        "$set": bson.M{
            "password": user.Password,
        },
    }

    result, err := userCollection.UpdateOne(ctx, filter, update)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Matched %v document(s) and updated %v document(s).\n", result.MatchedCount, result.ModifiedCount)
}

// DeleteUser removes a user from the database by email.
func DeleteUser(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := userCollection.DeleteOne(ctx, bson.M{"email": user.Email})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Deleted %v document(s)\n", result.DeletedCount)
}