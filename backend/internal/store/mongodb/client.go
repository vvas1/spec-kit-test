package mongodb

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBName is the default database name.
const DBName = "issue_tracker"

// NewClient connects to MongoDB using MONGODB_URI (default: mongodb://localhost:27017).
func NewClient(ctx context.Context) (*mongo.Client, error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	return mongo.Connect(ctx, options.Client().ApplyURI(uri))
}

// Database returns the database for the given client.
func Database(client *mongo.Client) *mongo.Database {
	return client.Database(DBName)
}
