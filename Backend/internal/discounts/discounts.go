package discounts

import (
	"encoding/json"
	"hospos-backend/internal/db"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Discount struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name      string             `json:"name"`
	Percent   float64            `json:"percent"`
	Type      string             `json:"type"`
	Code      string             `json:"code,omitempty"`
	ExpiresAt *time.Time         `json:"expiresAt,omitempty"`
	Active    bool               `json:"active"`
}

var collectionName = "discounts"

func DiscountsHandler(w http.ResponseWriter, r *http.Request) {
	coll, err := db.GetCollection(collectionName)
	if err != nil {
		log.Printf("db error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"db error"}`))
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/discounts")
	idStr := strings.TrimPrefix(path, "/")
	switch r.Method {
	case http.MethodGet:
		cursor, err := coll.Find(r.Context(), bson.M{})
		if err != nil {
			log.Printf("mongo find error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		var results []Discount
		if err := cursor.All(r.Context(), &results); err != nil {
			log.Printf("mongo decode error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		// Normalize type field to lowercase for all results
		for i := range results {
			results[i].Type = strings.ToLower(results[i].Type)
		}
		if err := json.NewEncoder(w).Encode(results); err != nil {
			log.Printf("error encoding discounts: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"internal error"}`))
		}
	case http.MethodPost:
		var d Discount
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid input"}`))
			return
		}
		// Normalize and validate type
		d.Type = strings.ToLower(d.Type)
		if d.Type != "static" && d.Type != "code" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"type must be 'static' or 'code'"}`))
			return
		}
		if d.Type == "code" && d.Code == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"code required for code-based discount"}`))
			return
		}
		if d.ExpiresAt != nil {
			d.Active = d.ExpiresAt.After(time.Now())
		} else {
			d.Active = true
		}
		_, err := coll.InsertOne(r.Context(), d)
		if err != nil {
			log.Printf("mongo insert error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(d); err != nil {
			log.Printf("error encoding discount: %v", err)
		}
	case http.MethodPut:
		// Edit discount
		if idStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"missing id"}`))
			return
		}
		oid, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid id"}`))
			return
		}
		var d Discount
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid input"}`))
			return
		}
		update := bson.M{
			"$set": bson.M{
				"name":      d.Name,
				"percent":   d.Percent,
				"type":      d.Type,
				"code":      d.Code,
				"expiresat": d.ExpiresAt,
				"active":    d.Active,
			},
		}
		_, err = coll.UpdateByID(r.Context(), oid, update)
		if err != nil {
			log.Printf("mongo update error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		// Remove discount
		if idStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"missing id"}`))
			return
		}
		oid, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid id"}`))
			return
		}
		_, err = coll.DeleteOne(r.Context(), bson.M{"_id": oid})
		if err != nil {
			log.Printf("mongo delete error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case http.MethodPatch:
		// PATCH /api/discounts/{id}/renew
		if !strings.HasSuffix(path, "/renew") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		idStr = strings.TrimSuffix(idStr, "/renew")
		oid, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid id"}`))
			return
		}
		// Set new expiry (default: 1 month from now)
		newExpiry := time.Now().AddDate(0, 1, 0)
		update := bson.M{"$set": bson.M{"expiresat": newExpiry, "active": true}}
		_, err = coll.UpdateByID(r.Context(), oid, update)
		if err != nil {
			log.Printf("mongo renew error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
