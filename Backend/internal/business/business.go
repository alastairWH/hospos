package business

import (
	"context"
	"encoding/json"
	"hospos-backend/internal/db"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BusinessInfo struct {
	CompanyName      string   `json:"companyName" bson:"companyName"`
	CompanyAddress   string   `json:"companyAddress" bson:"companyAddress"`
	FinanceEmail     string   `json:"financeEmail" bson:"financeEmail"`
	VATID            string   `json:"vatId" bson:"vatId"`
	CompanyRegNumber string   `json:"companyRegNumber" bson:"companyRegNumber"`
	Phone            string   `json:"phone" bson:"phone"`
	Website          string   `json:"website" bson:"website"`
	LogoURL          string   `json:"logoUrl" bson:"logoUrl"`
	SalesIDPrefix    string   `json:"salesIdPrefix" bson:"salesIdPrefix"`
	Currency         string   `json:"currency" bson:"currency"`
	DefaultTaxRate   float64  `json:"defaultTaxRate" bson:"defaultTaxRate"`
	BankDetails      string   `json:"bankDetails" bson:"bankDetails"`
	LegalFooter      string   `json:"legalFooter" bson:"legalFooter"`
	OpeningHours     string   `json:"openingHours" bson:"openingHours"`
	SocialLinks      []string `json:"socialLinks" bson:"socialLinks"`
	CustomReceiptMsg string   `json:"customReceiptMsg" bson:"customReceiptMsg"`
	InvoiceFormat    string   `json:"invoiceFormat" bson:"invoiceFormat"`
	Country          string   `json:"country" bson:"country"`
	LastSalesNumber  int      `json:"lastSalesNumber" bson:"lastSalesNumber"`
}

const businessCollection = "business"

func BusinessInfoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetBusinessInfoHandler(w, r)
	case http.MethodPost, http.MethodPut:
		SetBusinessInfoHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /api/business
func GetBusinessInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	coll, _ := db.GetCollection(businessCollection)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var info BusinessInfo
	err := coll.FindOne(ctx, bson.M{}).Decode(&info)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"not found"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

// POST /api/business
func SetBusinessInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var info BusinessInfo
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid input"}`))
		return
	}
	coll, _ := db.GetCollection(businessCollection)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	_, err := coll.ReplaceOne(ctx, bson.M{}, info, options.Replace().SetUpsert(true))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"db error"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}
