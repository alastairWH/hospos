package bookings

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

type BookingProduct struct {
	ProductID primitive.ObjectID `json:"productId" bson:"productId"`
	Name      string             `json:"name" bson:"name"`
	Qty       int                `json:"qty" bson:"qty"`
	Price     float64            `json:"price" bson:"price"`
}

type Booking struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CustomerID  primitive.ObjectID `json:"customerId" bson:"customerId"`
	TableNumber string             `json:"tableNumber" bson:"tableNumber"`
	Products    []BookingProduct   `json:"products" bson:"products"`
	BillTotal   float64            `json:"billTotal" bson:"billTotal"`
	Status      string             `json:"status" bson:"status"` // open, closed, cancelled
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	BookingTime time.Time          `json:"bookingTime" bson:"bookingTime"`
	ClosedAt    *time.Time         `json:"closedAt,omitempty" bson:"closedAt,omitempty"`
	Receipt     interface{}        `json:"receipt,omitempty" bson:"receipt,omitempty"`
	Notes       string             `json:"notes" bson:"notes"`
}

// Handler for /api/bookings and /api/bookings/{id}
func BookingsHandler(w http.ResponseWriter, r *http.Request) {
	// GET /api/bookings or /api/bookings/{id}
	if r.Method == http.MethodGet {
		parts := splitPath(r.URL.Path)
		coll, err := db.GetCollection("bookings")
		if err != nil {
			log.Printf("db error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		if len(parts) == 3 && parts[2] != "" { // /api/bookings/{id}
			objID, err := primitive.ObjectIDFromHex(parts[2])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"invalid id"}`))
				return
			}
			var booking Booking
			err = coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&booking)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error":"not found"}`))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(booking)
			return
		}
		// List all bookings
		cur, err := coll.Find(ctx, bson.M{})
		if err != nil {
			log.Printf("find error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		defer cur.Close(ctx)
		var bookings []Booking
		if err := cur.All(ctx, &bookings); err != nil {
			log.Printf("decode error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bookings)
		return
	}

	// POST /api/bookings
	if r.Method == http.MethodPost {
		var b Booking
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid input"}`))
			return
		}
		b.ID = primitive.NewObjectID()
		b.CreatedAt = time.Now()
		// Parse bookingTime from string if present (for compatibility)
		if b.BookingTime.IsZero() && r.FormValue("bookingTime") != "" {
			t, err := time.Parse(time.RFC3339, r.FormValue("bookingTime"))
			if err == nil {
				b.BookingTime = t
			}
		}
		if b.BookingTime.IsZero() {
			b.BookingTime = b.CreatedAt
		}
		if b.Status == "" {
			b.Status = "open"
		}
		coll, err := db.GetCollection("bookings")
		if err != nil {
			log.Printf("db error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		_, err = coll.InsertOne(ctx, b)
		if err != nil {
			log.Printf("insert error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(b)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)

	// PATCH /api/bookings/{id} for updating status or bookingTime
	if r.Method == http.MethodPatch {
		parts := splitPath(r.URL.Path)
		if len(parts) == 3 && parts[2] != "" {
			objID, err := primitive.ObjectIDFromHex(parts[2])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"invalid id"}`))
				return
			}
			var update struct {
				Status      string `json:"status"`
				BookingTime string `json:"bookingTime"`
			}
			if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"invalid input"}`))
				return
			}
			updateFields := bson.M{}
			if update.Status != "" {
				updateFields["status"] = update.Status
			}
			if update.BookingTime != "" {
				if t, err := time.Parse("2006-01-02T15:04", update.BookingTime); err == nil {
					updateFields["bookingTime"] = t
				}
			}
			if len(updateFields) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"no fields to update"}`))
				return
			}
			coll, err := db.GetCollection("bookings")
			if err != nil {
				log.Printf("db error: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"db error"}`))
				return
			}
			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			defer cancel()
			res, err := coll.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updateFields})
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
			return
		}
	}
}
