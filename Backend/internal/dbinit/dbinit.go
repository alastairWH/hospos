package dbinit

import (
	"context"
	"log"
	"time"

	"hospos-backend/internal/db"

	"go.mongodb.org/mongo-driver/bson"
)

// SeedData holds initial data for collections
// List of all collections to ensure they are created
var AllCollections = []string{
	"users",
	"roles",
	"discounts",
	"bookings",
	"locations",
	"inventory",
	"products",
	"categories",
	"sales",
	"payments",
	"receipts",
	"reminders",
}

var SeedData = map[string][]interface{}{
	"users": {
		bson.M{"name": "admin", "role": "admin"},
	},
	"roles": {
		bson.M{"role": "admin"},
		bson.M{"role": "manager"},
		bson.M{"role": "cashier"},
	},
}

// InitDB seeds the database with main information
func InitDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ensure all collections exist by inserting a dummy doc if empty, then removing it
	for _, collName := range AllCollections {
		coll, err := db.GetCollection(collName)
		if err != nil {
			log.Printf("dbinit: failed to get collection %s: %v", collName, err)
			return err
		}
		count, err := coll.CountDocuments(ctx, bson.M{})
		if err != nil {
			return err
		}
		if count == 0 {
			// Insert a dummy doc to create the collection
			res, err := coll.InsertOne(ctx, bson.M{"_init": true})
			if err != nil {
				log.Printf("dbinit: failed to create %s: %v", collName, err)
				return err
			}
			// Remove the dummy doc
			_, _ = coll.DeleteOne(ctx, bson.M{"_id": res.InsertedID})
			log.Printf("dbinit: ensured collection %s exists", collName)
		}
	}

	// Seed data for collections that need it
	for collName, docs := range SeedData {
		coll, err := db.GetCollection(collName)
		if err != nil {
			log.Printf("dbinit: failed to get collection %s: %v", collName, err)
			return err
		}
		count, err := coll.CountDocuments(ctx, bson.M{})
		if err != nil {
			return err
		}
		if count == 0 && len(docs) > 0 {
			_, err := coll.InsertMany(ctx, docs)
			if err != nil {
				log.Printf("dbinit: failed to seed %s: %v", collName, err)
				return err
			}
			log.Printf("dbinit: seeded %s", collName)
		}
	}
	return nil
}
