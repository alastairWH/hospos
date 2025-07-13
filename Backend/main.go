package main

import (
	"hospos-backend/internal/bookings"
	"hospos-backend/internal/business"
	"hospos-backend/internal/customers"
	"hospos-backend/internal/dbinit"
	"hospos-backend/internal/devtools"
	"hospos-backend/internal/discounts"
	"hospos-backend/internal/finance"
	"hospos-backend/internal/inventory"
	"hospos-backend/internal/linking"
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
	"net"
	"net/http"
	"os"
)

// Logging and error handling middleware
func withLoggingAndRecovery(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] %s %s: %v", r.Method, r.URL.Path, err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		log.Printf("[API] %s %s", r.Method, r.URL.Path)
		h(w, r)
	}
}

// CORS middleware
func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
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
	// Heartbeat
	mux.HandleFunc("/api/heartbeat", withLoggingAndRecovery(withCORS(linking.HeartbeatHandler)))
	// Product management
	mux.HandleFunc("/api/products", withLoggingAndRecovery(withCORS(products.ProductsHandler)))
	mux.HandleFunc("/api/products/", withLoggingAndRecovery(withCORS(products.ProductByIDHandler)))
	// Sales
	mux.HandleFunc("/api/sales", withLoggingAndRecovery(withCORS(sales.SalesHandler)))
	// Categories
	mux.HandleFunc("/api/categories", withLoggingAndRecovery(withCORS(products.CategoriesHandler)))
	// Table bookings
	mux.HandleFunc("/api/bookings", withLoggingAndRecovery(withCORS(bookings.BookingsHandler)))
	mux.HandleFunc("/api/bookings/", withLoggingAndRecovery(withCORS(bookings.BookingsHandler)))
	// Inventory
	mux.HandleFunc("/api/inventory", withLoggingAndRecovery(withCORS(inventory.InventoryHandler)))
	// Users
	mux.HandleFunc("/api/users", withLoggingAndRecovery(withCORS(users.UsersHandler)))
	mux.HandleFunc("/api/users/", withLoggingAndRecovery(withCORS(users.UsersHandler)))
	// Auth
	mux.HandleFunc("/api/auth", withLoggingAndRecovery(withCORS(users.AuthHandler)))
	// Roles
	mux.HandleFunc("/api/roles", withLoggingAndRecovery(withCORS(roles.RolesHandler)))
	// Reports
	mux.HandleFunc("/api/reports", withLoggingAndRecovery(withCORS(reports.ReportsHandler)))
	// Customers (list, add, update, delete, get by id)
	mux.HandleFunc("/api/customers", withLoggingAndRecovery(withCORS(customers.CustomersHandler)))
	mux.HandleFunc("/api/customers/", withLoggingAndRecovery(withCORS(customers.CustomersHandler)))
	// Payments
	mux.HandleFunc("/api/payments", withLoggingAndRecovery(withCORS(payments.PaymentsHandler)))
	// Receipts
	mux.HandleFunc("/api/receipts", withLoggingAndRecovery(withCORS(receipts.ReceiptsHandler)))
	// Discounts
	mux.HandleFunc("/api/discounts", withLoggingAndRecovery(withCORS(discounts.DiscountsHandler)))
	mux.HandleFunc("/api/discounts/", withLoggingAndRecovery(withCORS(discounts.DiscountsHandler)))
	// Locations
	mux.HandleFunc("/api/locations", withLoggingAndRecovery(withCORS(locations.LocationsHandler)))
	// Offline sync
	mux.HandleFunc("/api/sync", withLoggingAndRecovery(withCORS(sync.SyncHandler)))
	// Linking (till registration)
	mux.HandleFunc("/api/linking/link", withLoggingAndRecovery(withCORS(linking.LinkHandler)))
	// Reservation reminders
	mux.HandleFunc("/api/reminders", withLoggingAndRecovery(withCORS(reminders.RemindersHandler)))
	// DB initialization
	mux.HandleFunc("/api/dbinit", withLoggingAndRecovery(withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Printf("[ERROR] %s %s: Method not allowed", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		err := dbinit.InitDB()
		if err != nil {
			log.Printf("[ERROR] %s %s: db init failed: %v", r.Method, r.URL.Path, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"db init failed"}`))
			return
		}
		w.Write([]byte(`{"status":"db initialized"}`))
	})))
	// Finance summary
	mux.HandleFunc("/api/finance/summary", withLoggingAndRecovery(withCORS(finance.FinanceSummaryHandler)))
	// Devtools
	mux.HandleFunc("/api/devtools/seed", withLoggingAndRecovery(withCORS(devtools.SeedTestDataHandler)))
	mux.HandleFunc("/api/devtools/clear", withLoggingAndRecovery(withCORS(devtools.ClearTestDataHandler)))
	// Business info (combine GET and POST/PUT in one handler)
	mux.HandleFunc("/api/business", withLoggingAndRecovery(withCORS(business.BusinessInfoHandler)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Log external IP address
	ip, err := getOutboundIP()
	if err != nil {
		log.Printf("Could not determine external IP: %v", err)
	} else {
		log.Printf("External IP address: %s", ip)
	}
	log.Printf("Starting server on :%s...", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// getOutboundIP gets the preferred outbound IP of this machine
func getOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}
