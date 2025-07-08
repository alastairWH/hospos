package locations

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"hospos-backend/internal/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	LinkCode string             `json:"linkCode" bson:"linkCode"`
}

func LocationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		coll, err := db.GetCollection("locations")
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
		var locations []Location
		if err := cur.All(ctx, &locations); err != nil {
			log.Printf("decode error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		if locations == nil {
			locations = []Location{}
		}
		if err := json.NewEncoder(w).Encode(locations); err != nil {
			log.Printf("encode error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	case http.MethodPost:
		var l Location
		if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid input"}`))
			return
		}
		if l.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"name required"}`))
			return
		}
		// Generate a unique 12-digit numeric linking code
		l.ID = primitive.NewObjectID()
		l.LinkCode = generateNumericLinkCode(12)
		coll, err := db.GetCollection("locations")
		if err != nil {
			log.Printf("db error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		_, err = coll.InsertOne(ctx, l)
		if err != nil {
			log.Printf("insert error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(l); err != nil {
			log.Printf("encode error: %v", err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// generateNumericLinkCode returns a random numeric string of the given length
func generateNumericLinkCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	code := ""
	for i := 0; i < length; i++ {
		code += strconv.Itoa(rand.Intn(10))
	}
	return code
}
