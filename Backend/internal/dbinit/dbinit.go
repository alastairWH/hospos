package dbinit

import (
	"context"
	"log"
	"time"

	"hospos-backend/internal/db"

	"go.mongodb.org/mongo-driver/bson"
)

// SeedData holds initial data for collections
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
