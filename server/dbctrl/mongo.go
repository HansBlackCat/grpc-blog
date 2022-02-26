package dbctrl

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateNewClient get host for connecting to client
func CreateNewClient(host string) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(host))
	if err != nil {
		//zlog.Fatal().Msgf("%v", err)
		return nil, err
	}
	err = client.Connect(context.Background())
	if err != nil {
		//zlog.Fatal().Msgf("%v", err)
		return nil, err
	}
	return client, nil
}

// FetchCollection get mongodb's client, and names and return mongodb's collection
func FetchCollection(client *mongo.Client, dbName, collectionName string) *mongo.Collection {
	return client.Database(dbName).Collection(collectionName)
}
