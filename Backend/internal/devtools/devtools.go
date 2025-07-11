package devtools

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"hospos-backend/internal/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// POST /api/devtools/seed
func SeedTestDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()
	now := time.Now()
	rand.Seed(time.Now().UnixNano())

	// --- Customers ---
	customersColl, _ := db.GetCollection("customers")
	customers := make([]primitive.ObjectID, 0, 20)
	for i := 0; i < 20; i++ {
		id := primitive.NewObjectID()
		customers = append(customers, id)
		cust := bson.M{
			"_id":       id,
			"name":      "Test Customer " + strconv.Itoa(i),
			"phone":     "07000" + strconv.Itoa(i),
			"notes":     "Seeded customer",
			"tags":      []string{"test"},
			"createdAt": now.Add(-time.Duration(rand.Intn(1000)) * time.Hour),
		}
		customersColl.InsertOne(ctx, cust)
	}

	// --- Products ---
	productsColl, _ := db.GetCollection("products")
	products := make([]primitive.ObjectID, 0, 10)
	for i := 0; i < 10; i++ {
		id := primitive.NewObjectID()
		products = append(products, id)
		prod := bson.M{
			"_id":      id,
			"name":     "Product " + strconv.Itoa(i),
			"price":    float64(5 + rand.Intn(50)),
			"category": "Category " + strconv.Itoa(i%3),
		}
		productsColl.InsertOne(ctx, prod)
	}

	// --- Bookings ---
	bookingsColl, _ := db.GetCollection("bookings")
	for i := 0; i < 15; i++ {
		custIdx := rand.Intn(len(customers))
		prodIdx := rand.Intn(len(products))
		bookingID := primitive.NewObjectID()
		bookingProducts := []bson.M{
			{
				"productId": products[prodIdx],
				"name":      "Product " + strconv.Itoa(prodIdx),
				"qty":       1 + rand.Intn(3),
				"price":     float64(5 + rand.Intn(50)),
			},
		}
		booking := bson.M{
			"_id":         bookingID,
			"customerId":  customers[custIdx],
			"tableNumber": strconv.Itoa(1 + (i % 5)),
			"products":    bookingProducts,
			"billTotal":   float64(20 + rand.Intn(80)),
			"status":      []string{"open", "closed", "cancelled"}[rand.Intn(3)],
			"createdAt":   now.Add(-time.Duration(rand.Intn(1000)) * time.Hour),
			"bookingTime": now.Add(time.Duration(rand.Intn(1000)) * time.Hour),
			"notes":       "Seeded booking",
		}
		bookingsColl.InsertOne(ctx, booking)
	}

	// --- Sales ---
	salesColl, _ := db.GetCollection("sales")
	sales := make([]primitive.ObjectID, 0, 30)
	for i := 0; i < 30; i++ {
		prodIdx := rand.Intn(len(products))
		saleID := primitive.NewObjectID()
		sales = append(sales, saleID)
		sale := bson.M{
			"_id":        saleID,
			"product_id": products[prodIdx],
			"quantity":   1 + rand.Intn(5),
			"total":      float64(10 + rand.Intn(90)),
			"vat":        float64(rand.Intn(20)),
		}
		salesColl.InsertOne(ctx, sale)
	}

	// --- Payments ---
	paymentsColl, _ := db.GetCollection("payments")
	for i := 0; i < 30; i++ {
		saleIdx := rand.Intn(len(sales))
		pay := bson.M{
			"_id":     primitive.NewObjectID(),
			"sale_id": sales[saleIdx],
			"amount":  float64(10 + rand.Intn(90)),
			"method":  []string{"cash", "card", "online"}[rand.Intn(3)],
		}
		paymentsColl.InsertOne(ctx, pay)
	}

	// --- Receipts ---
	receiptsColl, _ := db.GetCollection("receipts")
	for i := 0; i < 10; i++ {
		saleIdx := rand.Intn(len(sales))
		receipt := bson.M{
			"_id":     primitive.NewObjectID(),
			"sale_id": sales[saleIdx],
			"detail":  "Seeded receipt detail",
		}
		receiptsColl.InsertOne(ctx, receipt)
	}

	// --- Discounts ---
	discountsColl, _ := db.GetCollection("discounts")
	for i := 0; i < 5; i++ {
		discount := bson.M{
			"_id":     primitive.NewObjectID(),
			"name":    "Discount " + strconv.Itoa(i),
			"percent": float64(5 + rand.Intn(25)),
			"type":    []string{"percent", "fixed"}[rand.Intn(2)],
			"code":    "CODE" + strconv.Itoa(i),
			"active":  rand.Intn(2) == 0,
		}
		discountsColl.InsertOne(ctx, discount)
	}

	// --- Roles ---
	rolesColl, _ := db.GetCollection("roles")
	for _, role := range []string{"admin", "manager", "cashier"} {
		rolesColl.InsertOne(ctx, bson.M{"role": role})
	}

	// --- Users ---
	usersColl, _ := db.GetCollection("users")
	for i := 0; i < 5; i++ {
		usersColl.InsertOne(ctx, bson.M{
			"_id":       primitive.NewObjectID(),
			"name":      "user" + strconv.Itoa(i),
			"pin":       "$2a$10$7EqJtq98hPqEX7fNZaFWoO5r5rQ9g6rQ9g6rQ9g6rQ9g6rQ9g6rQ9G", // bcrypt hash for '1234'
			"role":      []string{"admin", "manager", "cashier"}[i%3],
			"createdAt": now.Add(-time.Duration(rand.Intn(1000)) * time.Hour),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bson.M{"success": true})
}

// POST /api/devtools/clear
func ClearTestDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	for _, collName := range []string{"customers", "products", "sales", "payments", "receipts"} {
		coll, _ := db.GetCollection(collName)
		coll.DeleteMany(ctx, bson.M{})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bson.M{"success": true})
}
