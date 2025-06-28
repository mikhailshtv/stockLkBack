package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Counter struct {
	ID  string `bson:"_id"`
	Seq int    `bson:"seq"`
}

func getNextSequence(db *mongo.Database, counterName string) (int, error) {
	collection := db.Collection("counters")

	filter := bson.M{"_id": counterName}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var counter Counter
	err := collection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&counter)

	if err != nil {
		return 0, err
	}

	return counter.Seq, nil
}
