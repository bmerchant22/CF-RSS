package store

import (
	"context"
	"fmt"
	"github.com/bmerchant22/project/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

const uri = "mongodb://localhost:27017/"

type MongoStore struct {
	RecentActionCollection *mongo.Collection
}

func (m *MongoStore) ConnectToDatabase() {

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}
	/*defer func main() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()*/

	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
		return
	}

	fmt.Println("Pinged the primary node of the cluster. You successfully connected to MongoDB!")

	m.RecentActionCollection = client.Database("config").Collection("recent_actions_demo")

}

func (m *MongoStore) StoreRecentActionsInTheDatabase(actions []models.RecentAction) error {
	var toInsertInterface []interface{}
	for _, action := range actions {
		toInsertInterface = append(toInsertInterface, action)
	}

	zap.S().Info("Trying to insert document to mongodb")

	_, err1 := m.RecentActionCollection.InsertMany(context.TODO(), toInsertInterface)
	if err1 != nil {
		zap.S().Errorf("Error while inserting documents: %v", err1)
		return nil
	}
	zap.S().Info("Insertion successful")
	return nil
}

func (m *MongoStore) GetMaxTimeStamp() (int64, error) {
	//aggregate query
	//cursor.next
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", nil},
			{"MaxTimeStamp", bson.D{{"$max", "$timeseconds"}}},
		}},
	}
	cursor, err := m.RecentActionCollection.Aggregate(context.TODO(), mongo.Pipeline{groupStage})
	if err != nil {
		panic(err)
	}
	wrapper := struct {
		MaxTimeStamp int64
	}{}
	for cursor.Next(context.TODO()) {
		if err = cursor.Decode(&wrapper); err != nil {
			panic(err)
		}
	}
	return wrapper.MaxTimeStamp, err
}

func (m *MongoStore) QueryRecentActions(after int64) ([]models.RecentAction, error) {
	filter := bson.M{"timeseconds": bson.M{"$gte": after}}
	cursor, err := m.RecentActionCollection.Find(context.TODO(), filter)
	if err != nil {
		zap.S().Errorf("Error while quering mongoDB")
	}

	zap.S().Infof("Quering mongodb after %d timeStamp", after)
	zap.S().Infof("Using filter to query : %+v", filter)

	var results []models.RecentAction
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		zap.S().Errorf("Error while unmarshalling")
	}

	for cursor.Next(context.TODO()) {
		if err = cursor.Decode(&results); err != nil {
			panic(err)
		}
	}
	return results, err
}
