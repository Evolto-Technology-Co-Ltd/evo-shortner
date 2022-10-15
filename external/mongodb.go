package external

import (
	"context"
	"fmt"
	"time"

	"github.com/ariebrainware/evo-shortner/config"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoConn(document string) *mongo.Collection {
	// load app.env file data to struct
	config, err := config.LoadConfig(".")
	// handle errors
	if err != nil {
		log.Fatalf("can't load environment app.env: %v", err)
	}

	fmt.Println("Mongo Database: ", config.MongoDatabase)

	var client *mongo.Client
	var collection *mongo.Collection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err = mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:27017", config.MongoHost)))
	if err != nil {
		log.Error(err)
		panic("Failed to connect mongo")
	}
	collection = client.Database(config.MongoDatabase).Collection(document)
	err = client.Connect(ctx)
	if err != nil {
		log.Error(err)
		panic("Failed to connect mongo")
	}
	return collection
}
