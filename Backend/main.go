package main

import (
	"log"
	"net/http"
	"os"

	"hospos-backend/internal/bookings"
	"hospos-backend/internal/customers"
	"hospos-backend/internal/dbinit"
	"hospos-backend/internal/discounts"
	"hospos-backend/internal/inventory"
	"hospos-backend/internal/locations"
	"hospos-backend/internal/payments"
	"hospos-backend/internal/products"
	"hospos-backend/internal/receipts"
	"hospos-backend/internal/reminders"
	"hospos-backend/internal/reports"
	"hospos-backend/internal/sales"
	"hospos-backend/internal/sync"
	"hospos-backend/internal/users"
)

func main() {
	mux := http.NewServeMux()

	// Product management
	mux.HandleFunc("/api/products", products.ProductsHandler)
	mux.HandleFunc("/api/products/", products.ProductByIDHandler)
	// Sales
	mux.HandleFunc("/api/sales", sales.SalesHandler)
	// Table bookings
	mux.HandleFunc("/api/bookings", bookings.BookingsHandler)
	// Inventory
	mux.HandleFunc("/api/inventory", inventory.InventoryHandler)
	// Users
	mux.HandleFunc("/api/users", users.UsersHandler)
	// Reports
	mux.HandleFunc("/api/reports", reports.ReportsHandler)
	// Customers
	mux.HandleFunc("/api/customers", customers.CustomersHandler)
	// Payments
	mux.HandleFunc("/api/payments", payments.PaymentsHandler)
	// Receipts
	mux.HandleFunc("/api/receipts", receipts.ReceiptsHandler)
	// Discounts
	mux.HandleFunc("/api/discounts", discounts.DiscountsHandler)
	// Locations
	mux.HandleFunc("/api/locations", locations.LocationsHandler)
	// Offline sync
	mux.HandleFunc("/api/sync", sync.SyncHandler)
	// Reservation reminders
	mux.HandleFunc("/api/reminders", reminders.RemindersHandler)
	// DB initialization
	mux.HandleFunc("/api/dbinit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		err := dbinit.InitDB()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db init failed"}`))
			return
		}
		w.Write([]byte(`{"status":"db initialized"}`))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on :%s...", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
