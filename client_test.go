package mongoout

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestMongoDBClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	collection := client.Database("aaa").Collection("bbb")
	for i := 0; i < 10; i++ {
		_, err := collection.InsertOne(context.Background(), bson.M{"i": i})
		if err != nil {
			t.Error(err)
		}
		t.Log(i)
	}
}
