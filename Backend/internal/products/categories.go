package products

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

type Category struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
}

// CategoriesHandler handles /api/categories for GET and POST
func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		coll, err := db.GetCollection("categories")
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
			   var categories []Category
			   if err := cur.All(ctx, &categories); err != nil {
					   log.Printf("decode error: %v", err)
					   w.WriteHeader(http.StatusInternalServerError)
					   w.Write([]byte(`{"error":"db error"}`))
					   return
			   }
			   // Always return an array, even if empty
			   if categories == nil {
					   categories = []Category{}
			   }
			   if err := json.NewEncoder(w).Encode(categories); err != nil {
					   log.Printf("encode error: %v", err)
					   w.WriteHeader(http.StatusInternalServerError)
			   }
	case http.MethodPost:
		var c Category
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid input"}`))
			return
		}
		if c.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"name required"}`))
			return
		}
		coll, err := db.GetCollection("categories")
		if err != nil {
			log.Printf("db error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		res, err := coll.InsertOne(ctx, bson.M{"name": c.Name})
		if err != nil {
			log.Printf("insert error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		c.ID = res.InsertedID.(primitive.ObjectID)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(c)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
