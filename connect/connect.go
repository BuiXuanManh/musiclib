package connect

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client   *mongo.Client
	Database *mongo.Database
}

var Ng MongoInstance

func Connect() error {
	var dbName = os.Getenv("DB_NAME")
	var url = os.Getenv("DB_URL")
	mongoOption := options.Client().ApplyURI(url + dbName)
	client, err := mongo.Connect(context.Background(), mongoOption)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	Ng = MongoInstance{
		Client:   client,
		Database: client.Database(dbName),
	}
	return nil
}
