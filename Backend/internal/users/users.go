package users

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"hospos-backend/internal/db"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	log.Printf("[AUTH] Attempt login: name='%s', pin='%s'", req.Name, req.Pin)
	coll, err := db.GetCollection("users")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"db error"}`))
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var user User
	err = coll.FindOne(ctx, bson.M{"name": req.Name}).Decode(&user)
	if err != nil {
		log.Printf("[AUTH] User not found: name='%s'", req.Name)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"invalid credentials"}`))
		return
	}
	log.Printf("[AUTH] DB hash for user '%s': %s", user.Name, user.Pin)
	// Compare hashed PIN
	if bcrypt.CompareHashAndPassword([]byte(user.Pin), []byte(req.Pin)) != nil {
		log.Printf("[AUTH] PIN mismatch for user '%s'", user.Name)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"invalid credentials"}`))
		return
	}
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
	Pin  string `json:"pin" bson:"pin"`
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
		// Validate pin: Must only be Numbers
		for _, c := range u.Pin {
			if c < '0' || c > '9' {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"pin must be digits only"}`))
				return
			}
		}
		// Hash the PIN
		hashedPin, err := bcrypt.GenerateFromPassword([]byte(u.Pin), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("bcrypt error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"hash error"}`))
			return
		}
		u.Pin = string(hashedPin)
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
		// Support /api/users/{id}/pin (PUT) and /api/users/{id} (DELETE)
		parts := splitPath(r.URL.Path)
		if len(parts) >= 4 && parts[3] == "pin" && r.Method == http.MethodPut {
			// Handle PIN update
			id := parts[2]
			var req struct {
				Pin string `json:"pin"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"invalid input"}`))
				return
			}
			//Validator: Must be within the 3 - 6 digits
			if len(req.Pin) < 3 || len(req.Pin) > 6 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"pin must be 3-6 digits"}`))
				return
			}
			//Validator: Must be only Number, no Letters
			for _, c := range req.Pin {
				if c < '0' || c > '9' {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"error":"pin must be digits only"}`))
					return
				}
			}
			// Hash the new PIN
			hashedPin, err := bcrypt.GenerateFromPassword([]byte(req.Pin), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("bcrypt error: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"hash error"}`))
				return
			}
			// DB Save Error
			coll, err := db.GetCollection("users")
			if err != nil {
				log.Printf("db error: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"db error"}`))
				return
			}
			// DB Error: Invalid ID read from Table
			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"invalid user id"}`))
				return
			}
			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			defer cancel()
			res, err := coll.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"pin": string(hashedPin)}})
			// DB Error: Pin Update Error
			if err != nil {
				log.Printf("update pin error: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"db error"}`))
				return
			}
			// Error Handle: No user found In DB
			if res.MatchedCount == 0 {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error":"user not found"}`))
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		} else if len(parts) >= 3 && parts[2] != "" && r.Method == http.MethodDelete {
			// Handle user delete
			id := parts[2]
			coll, err := db.GetCollection("users")
			if err != nil {
				log.Printf("db error: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"db error"}`))
				return
			}
			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"invalid user id"}`))
				return
			}
			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			defer cancel()
			res, err := coll.DeleteOne(ctx, bson.M{"_id": objID})
			if err != nil {
				log.Printf("delete error: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"db error"}`))
				return
			}
			if res.DeletedCount == 0 {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error":"user not found"}`))
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// splitPath splits a URL path into its segments, ignoring leading/trailing slashes.
func splitPath(path string) []string {
	var segs []string
	start := 0
	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			if i > start {
				segs = append(segs, path[start:i])
			}
			start = i + 1
		}
	}
	if start < len(path) {
		segs = append(segs, path[start:])
	}
	return segs
}
