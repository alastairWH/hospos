package linking

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"hospos-backend/internal/db"

	"go.mongodb.org/mongo-driver/bson"
)

type HeartbeatRequest struct {
	TillID     string                 `json:"tillId"`
	DeviceInfo map[string]interface{} `json:"deviceInfo"`
}

type HeartbeatResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// HeartbeatHandler handles /api/heartbeat for terminal online status
func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req HeartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success":false,"error":"invalid input"}`))
		return
	}
	if req.TillID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success":false,"error":"tillId required"}`))
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
	update := bson.M{"$set": bson.M{"lastSeen": time.Now(), "deviceInfo": req.DeviceInfo}}
	_, err = coll.UpdateOne(ctx, bson.M{"_id": req.TillID}, update)
	if err != nil {
		log.Printf("update error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success":false,"error":"db error"}`))
		return
	}
	resp := HeartbeatResponse{Success: true}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
