package main

import (
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
	"hospos-backend/internal/roles"
	"hospos-backend/internal/sales"
	"hospos-backend/internal/sync"
	"hospos-backend/internal/users"
	"log"
	"net/http"
	"os"
)

// CORS middleware
func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}

func main() {
	mux := http.NewServeMux()

	// Product management
	mux.HandleFunc("/api/products", withCORS(products.ProductsHandler))
	mux.HandleFunc("/api/products/", withCORS(products.ProductByIDHandler))
	// Sales
	mux.HandleFunc("/api/sales", withCORS(sales.SalesHandler))
	// Table bookings
	mux.HandleFunc("/api/bookings", withCORS(bookings.BookingsHandler))
	// Inventory
	mux.HandleFunc("/api/inventory", withCORS(inventory.InventoryHandler))
	// Users
	mux.HandleFunc("/api/users", withCORS(users.UsersHandler))
	// Auth
	mux.HandleFunc("/api/auth", withCORS(users.AuthHandler))
	// Roles
	mux.HandleFunc("/api/roles", withCORS(roles.RolesHandler))
	// Reports
	mux.HandleFunc("/api/reports", withCORS(reports.ReportsHandler))
	// Customers
	mux.HandleFunc("/api/customers", withCORS(customers.CustomersHandler))
	// Payments
	mux.HandleFunc("/api/payments", withCORS(payments.PaymentsHandler))
	// Receipts
	mux.HandleFunc("/api/receipts", withCORS(receipts.ReceiptsHandler))
	// Discounts
	mux.HandleFunc("/api/discounts", withCORS(discounts.DiscountsHandler))
	// Locations
	mux.HandleFunc("/api/locations", withCORS(locations.LocationsHandler))
	// Offline sync
	mux.HandleFunc("/api/sync", withCORS(sync.SyncHandler))
	// Reservation reminders
	mux.HandleFunc("/api/reminders", withCORS(reminders.RemindersHandler))
	// DB initialization
	mux.HandleFunc("/api/dbinit", withCORS(func(w http.ResponseWriter, r *http.Request) {
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
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on :%s...", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
