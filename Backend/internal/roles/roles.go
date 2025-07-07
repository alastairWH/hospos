package roles

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

type Role struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Role string             `json:"role" bson:"role"`
}

func RolesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		coll, err := db.GetCollection("roles")
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
		var roles []Role
		if err := cur.All(ctx, &roles); err != nil {
			log.Printf("decode error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		if err := json.NewEncoder(w).Encode(roles); err != nil {
			log.Printf("encode error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	case http.MethodPost:
		var role Role
		if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid input"}`))
			return
		}
		coll, err := db.GetCollection("roles")
		if err != nil {
			log.Printf("db error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		res, err := coll.InsertOne(ctx, role)
		if err != nil {
			log.Printf("insert error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		role.ID = res.InsertedID.(primitive.ObjectID)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(role); err != nil {
			log.Printf("encode error: %v", err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
