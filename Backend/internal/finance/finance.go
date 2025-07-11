package finance

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"hospos-backend/internal/db"
)

type FinanceSummary struct {
	TotalSales    float64 `json:"totalSales"`
	TotalVAT      float64 `json:"totalVAT"`
	TotalPayments float64 `json:"totalPayments"`
	TotalReceipts int     `json:"totalReceipts"`
}

// GET /api/finance/summary
func FinanceSummaryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Aggregate sales
	salesColl, err := db.GetCollection("sales")
	if err != nil {
		log.Printf("db error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"db error"}`))
		return
	}
	var salesTotal, vatTotal float64
	cur, err := salesColl.Find(ctx, map[string]interface{}{})
	if err == nil {
		var sales []struct {
			Total float64 `json:"total"`
			VAT   float64 `json:"vat"`
		}
		if err := cur.All(ctx, &sales); err == nil {
			for _, s := range sales {
				salesTotal += s.Total
				vatTotal += s.VAT
			}
		}
	}

	// Aggregate payments
	paymentsColl, err := db.GetCollection("payments")
	var paymentsTotal float64
	if err == nil {
		cur, err := paymentsColl.Find(ctx, map[string]interface{}{})
		if err == nil {
			var payments []struct {
				Amount float64 `json:"amount"`
			}
			if err := cur.All(ctx, &payments); err == nil {
				for _, p := range payments {
					paymentsTotal += p.Amount
				}
			}
		}
	}

	// Count receipts
	receiptsColl, err := db.GetCollection("receipts")
	receiptsCount := 0
	if err == nil {
		count, err := receiptsColl.CountDocuments(ctx, map[string]interface{}{})
		if err == nil {
			receiptsCount = int(count)
		}
	}

	summary := FinanceSummary{
		TotalSales:    salesTotal,
		TotalVAT:      vatTotal,
		TotalPayments: paymentsTotal,
		TotalReceipts: receiptsCount,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
