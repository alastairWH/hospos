package sales

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

type SaleProduct struct {
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	Name      string             `json:"name" bson:"name"`
	Quantity  int                `json:"qty" bson:"qty"`
	Price     float64            `json:"price" bson:"price"`
}

type SalePayment struct {
	Amount float64 `json:"amount" bson:"amount"`
	Method string  `json:"method" bson:"method"`
}

type Sale struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Products []SaleProduct      `json:"products" bson:"products"`
	Total    float64            `json:"total" bson:"total"`
	VAT      float64            `json:"vat" bson:"vat"`
	Discount float64            `json:"discount" bson:"discount"`
	Paid     float64            `json:"paid" bson:"paid"`
	Payments []SalePayment      `json:"payments" bson:"payments"`
}

// No in-memory sales; use MongoDB

func SalesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		coll, err := db.GetCollection("sales")
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
		var sales []Sale
		if err := cur.All(ctx, &sales); err != nil {
			log.Printf("decode error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		if err := json.NewEncoder(w).Encode(sales); err != nil {
			log.Printf("encode error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	case http.MethodPost:
		var s Sale
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid input"}`))
			return
		}
		// Calculate VAT (20% UK standard) if not provided
		if s.VAT == 0 {
			s.VAT = s.Total * 0.2
		}
		coll, err := db.GetCollection("sales")
		if err != nil {
			log.Printf("db error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		res, err := coll.InsertOne(ctx, s)
		if err != nil {
			log.Printf("insert error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db error"}`))
			return
		}
		s.ID = res.InsertedID.(primitive.ObjectID)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(s); err != nil {
			log.Printf("encode error: %v", err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
