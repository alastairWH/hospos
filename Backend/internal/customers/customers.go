package customers

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

// CustomerResponse is used to marshal ObjectID as string for JSON responses
type CustomerResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Notes     string    `json:"notes"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
}

// Helper to convert Customer to CustomerResponse
func customerToResponse(c Customer) CustomerResponse {
	return CustomerResponse{
		ID:        c.ID.Hex(),
		Name:      c.Name,
		Email:     c.Email,
		Phone:     c.Phone,
		Notes:     c.Notes,
		Tags:      c.Tags,
		CreatedAt: c.CreatedAt,
	}
}

type Customer struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Phone     string             `json:"phone" bson:"phone"`
	Notes     string             `json:"notes" bson:"notes"`
	Tags      []string           `json:"tags" bson:"tags"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

// GET: list/search, POST: add, PUT: update, DELETE: remove
// Handles /api/customers and /api/customers/{id}
func CustomersHandler(w http.ResponseWriter, r *http.Request) {
	// Check if path is /api/customers/{id}
	if r.URL.Path != "/api/customers" && r.Method == http.MethodGet {
		// Extract ID from path
		parts := splitPath(r.URL.Path)
		if len(parts) == 3 && parts[1] == "customers" && parts[2] != "" {
			id := parts[2]
			objID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"invalid id"}`))
				return
			}
			coll, err := db.GetCollection("customers")
			if err != nil {
				log.Printf("db error: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"db error"}`))
				return
			}
			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			defer cancel()
			var customer Customer
			err = coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&customer)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error":"not found"}`))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(customerToResponse(customer))
			return
		}
	}

	switch r.Method {
	// Helper to split path into parts

	case http.MethodGet:
		// Search by name/email/phone if query param present
		query := r.URL.Query().Get("q")
		filter := bson.M{}
		if query != "" {
			regex := primitive.Regex{Pattern: query, Options: "i"}
			filter = bson.M{"$or": []bson.M{
				{"name": regex},
				{"email": regex},
				{"phone": regex},
			}}
		}
		coll, err := db.GetCollection("customers")
		if err != nil {
			log.Printf("db error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		cur, err := coll.Find(ctx, filter)
		if err != nil {
			log.Printf("find error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		var customers []Customer
		if err := cur.All(ctx, &customers); err != nil {
			log.Printf("decode error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		// Convert to []CustomerResponse
		responses := make([]CustomerResponse, len(customers))
		for i, c := range customers {
			responses[i] = customerToResponse(c)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responses)

	case http.MethodPost:
		var c Customer
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid input"}`))
			return
		}
		c.ID = primitive.NewObjectID()
		c.CreatedAt = time.Now()
		if c.Tags == nil {
			c.Tags = []string{}
		}
		coll, err := db.GetCollection("customers")
		if err != nil {
			log.Printf("db error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		_, err = coll.InsertOne(ctx, c)
		if err != nil {
			log.Printf("insert error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(c)

	case http.MethodPut, http.MethodPatch:
		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"missing id"}`))
			return
		}
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid id"}`))
			return
		}
		var update Customer
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid input"}`))
			return
		}
		coll, err := db.GetCollection("customers")
		if err != nil {
			log.Printf("db error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		updateMap := bson.M{}
		if update.Name != "" {
			updateMap["name"] = update.Name
		}
		if update.Email != "" {
			updateMap["email"] = update.Email
		}
		if update.Phone != "" {
			updateMap["phone"] = update.Phone
		}
		if update.Notes != "" {
			updateMap["notes"] = update.Notes
		}
		if update.Tags != nil {
			updateMap["tags"] = update.Tags
		}
		if len(updateMap) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"no fields to update"}`))
			return
		}
		res, err := coll.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateMap})
		if err != nil {
			log.Printf("update error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		if res.MatchedCount == 0 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"not found"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true}`))

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"missing id"}`))
			return
		}
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid id"}`))
			return
		}
		coll, err := db.GetCollection("customers")
		if err != nil {
			log.Printf("db error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
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
			w.Write([]byte(`{"error":"not found"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true}`))

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
