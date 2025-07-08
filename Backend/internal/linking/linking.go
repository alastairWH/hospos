package linking

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

type LinkRequest struct {
	LinkCode   string                 `json:"linkCode"`
	DeviceInfo map[string]interface{} `json:"deviceInfo"`
}

type LinkResponse struct {
	Success     bool                   `json:"success"`
	TillID      string                 `json:"tillId,omitempty"`
	InitialData map[string]interface{} `json:"initialData,omitempty"`
	Error       string                 `json:"error,omitempty"`
}

// LinkHandler handles /api/linking/link for till registration
func LinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req LinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success":false,"error":"invalid input"}`))
		return
	}
	if req.LinkCode == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success":false,"error":"link code required"}`))
		return
	}
	coll, err := db.GetCollection("locations")
	if err != nil {
		log.Printf("db error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success":false,"error":"db error"}`))
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var location bson.M
	err = coll.FindOne(ctx, bson.M{"linkCode": req.LinkCode}).Decode(&location)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"success":false,"error":"invalid link code"}`))
		return
	}
	// Optionally: mark as linked, store device info, update last seen
	update := bson.M{"$set": bson.M{"linked": true, "deviceInfo": req.DeviceInfo, "lastSeen": time.Now()}}
	_, err = coll.UpdateOne(ctx, bson.M{"_id": location["_id"]}, update)
	if err != nil {
		log.Printf("update error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success":false,"error":"db error"}`))
		return
	}
	// Gather initial data (example: products, categories, users, roles)
	initialData := make(map[string]interface{})
	initialData["products"] = fetchAll(ctx, "products")
	initialData["categories"] = fetchAll(ctx, "categories")
	initialData["users"] = fetchAll(ctx, "users")
	initialData["roles"] = fetchAll(ctx, "roles")
	resp := LinkResponse{
		Success:     true,
		TillID:      location["_id"].(primitive.ObjectID).Hex(),
		InitialData: initialData,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// fetchAll returns all documents from a collection as []bson.M
func fetchAll(ctx context.Context, collection string) []bson.M {
	coll, err := db.GetCollection(collection)
	if err != nil {
		return nil
	}
	cur, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil
	}
	defer cur.Close(ctx)
	var results []bson.M
	if err := cur.All(ctx, &results); err != nil {
		return nil
	}
	return results
}
