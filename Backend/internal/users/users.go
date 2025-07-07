package users

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"hospos-backend/internal/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ...existing code...

// AuthHandler handles POST /api/auth for PIN-based login
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Name string `json:"name"`
		Pin  string `json:"pin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid input"}`))
		return
	}
	coll, err := db.GetCollection("users")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"db error"}`))
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var user User
	err = coll.FindOne(ctx, bson.M{"name": req.Name, "pin": req.Pin}).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"invalid credentials"}`))
		return
	}
	// For now, return user info as a simple token (replace with JWT in production)
	resp := struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Role  string `json:"role"`
		Token string `json:"token"`
	}{
		ID:    user.ID,
		Name:  user.Name,
		Role:  user.Role,
		Token: user.ID + ":" + user.Role, // placeholder token
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type User struct {
	ID   string `json:"id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Role string `json:"role" bson:"role"`
	Pin  string `json:"pin" bson:"pin"` // 3-6 digit pin for login
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		coll, err := db.GetCollection("users")
		if err != nil {
			log.Printf("db error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		cur, err := coll.Find(ctx, bson.M{})
		if err != nil {
			log.Printf("find error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		defer cur.Close(ctx)
		var users []User
		if err := cur.All(ctx, &users); err != nil {
			log.Printf("decode error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		if err := json.NewEncoder(w).Encode(users); err != nil {
			log.Printf("encode error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	case http.MethodPost:
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid input"}`))
			return
		}
		// Validate pin: must be 3-6 digits
		if len(u.Pin) < 3 || len(u.Pin) > 6 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"pin must be 3-6 digits"}`))
			return
		}
		for _, c := range u.Pin {
			if c < '0' || c > '9' {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"pin must be digits only"}`))
				return
			}
		}
		coll, err := db.GetCollection("users")
		if err != nil {
			log.Printf("db error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		res, err := coll.InsertOne(ctx, u)
		if err != nil {
			log.Printf("insert error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		u.ID = res.InsertedID.(primitive.ObjectID).Hex()
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(u); err != nil {
			log.Printf("encode error: %v", err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
