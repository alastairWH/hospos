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

type Product struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Price    float64            `json:"price" bson:"price"`
	Category string             `json:"category,omitempty" bson:"category,omitempty"`
}

// No in-memory products; use MongoDB

// ProductsHandler handles /api/products for GET and POST
func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		coll, err := db.GetCollection("products")
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
		var products []Product
		if err := cur.All(ctx, &products); err != nil {
			log.Printf("decode error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		if err := json.NewEncoder(w).Encode(products); err != nil {
			log.Printf("encode error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	case http.MethodPost:
		var p Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid input"}`))
			return
		}
		coll, err := db.GetCollection("products")
		if err != nil {
			log.Printf("db error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		res, err := coll.InsertOne(ctx, p)
		if err != nil {
			log.Printf("insert error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		p.ID = res.InsertedID.(primitive.ObjectID)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(p); err != nil {
			log.Printf("encode error: %v", err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// ProductByIDHandler handles /api/products/{id} for GET, PUT, DELETE
func ProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/products/"):]
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid id"}`))
		return
	}
	coll, err := db.GetCollection("products")
	if err != nil {
		log.Printf("db error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"db error"}`))
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	switch r.Method {
	case http.MethodGet:
		var p Product
		err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"not found"}`))
			return
		}
		if err := json.NewEncoder(w).Encode(p); err != nil {
			log.Printf("encode error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	case http.MethodPut:
		var p Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid input"}`))
			return
		}
		p.ID = id
		_, err := coll.ReplaceOne(ctx, bson.M{"_id": id}, p)
		if err != nil {
			log.Printf("update error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		if err := json.NewEncoder(w).Encode(p); err != nil {
			log.Printf("encode error: %v", err)
		}
	case http.MethodDelete:
		_, err := coll.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil {
			log.Printf("delete error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
