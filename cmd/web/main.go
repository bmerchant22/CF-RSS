/*package main

import (
	"encoding/json"
	"fmt"
	"github.com/bmerchant22/project/pkg/cfapi"
	"go.uber.org/zap"
	"log"
)

func main() {
	environment := "development"
	var logger *zap.Logger
	var loggerErr error

	if environment == "development" {
		if logger, loggerErr = zap.NewDevelopment(); loggerErr != nil {
			log.Fatalln(loggerErr)
		}
	} else {
		if logger, loggerErr = zap.NewProduction(); loggerErr != nil {
			log.Fatalln(loggerErr)
		}
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	obj := cfapi.NewCodeforcesClient()
	//obj.RecentActions(1)
	recentActions, err := obj.RecentActions(1)
	if err != nil {
		fmt.Println("error occured")
		return
	}
	//zap.S().Info(recentActions)
	data, err1 := json.MarshalIndent(recentActions, "", " ")
	if err1 != nil {
		fmt.Println("error occurred")
		return
	}
	zap.S().Info(string(data))
}*/

// Paste your connection string URI here
package main

import (
	"context"
	"fmt"
	"github.com/bmerchant22/project/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const uri = "mongodb://localhost:27017/"

func main() {

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Pinged the primary node of the cluster. You successfully connected to MongoDB!")
	//opts := options.Client().SetTimeout(5 * time.Second)
	collection := client.Database("config").Collection("recent_actions")

	filter := bson.M{
		"blogEntry.rating": 69,
		"comment": bson.M{
			"$exists": true,
		},
	}
	resp := collection.FindOne(context.TODO(), filter)
	var respDecode models.RecentAction
	err = resp.Decode(&respDecode)
	if err != nil {
		panic(err)
	}
	fmt.Println(respDecode.Comment)
}
