package storage

import (
	"testing"
)

func TestMongoDBConnection(t *testing.T) {
	mongodb, err := NewMongoDB("mongodb://localhost:27017", "cbb_analytics")
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongodb.Close()

	t.Log("✅ Successfully connected to MongoDB!")
	t.Log("✅ Database indexes created!")
}
