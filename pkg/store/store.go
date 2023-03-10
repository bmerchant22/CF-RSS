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
	Collection *mongo.Collection
}

func (m *MongoStore) ConnectToDatabaseRecentAction() {

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

	fmt.Println("Pinged the primary node of the cluster. You successfully connected to MongoDB! Collection : recent_actions_demo")

	m.Collection = client.Database("config").Collection("recent_actions_demo")

}

func (m *MongoStore) ConnectToDatabaseUsers() {

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

	fmt.Println("Pinged the primary node of the cluster. You successfully connected to MongoDB! Collection : Users")

	m.Collection = client.Database("config").Collection("users")

}

func (m *MongoStore) StoreRecentActionsInTheDatabase(actions []models.RecentAction) error {
	var toInsertInterface []interface{}
	for _, action := range actions {
		toInsertInterface = append(toInsertInterface, action)
	}

	zap.S().Info("Trying to insert document to mongodb")

	_, err1 := m.Collection.InsertMany(context.TODO(), toInsertInterface)
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
			{"MaxTimeStamp", bson.D{{"$max", "$timeSeconds"}}},
		}},
	}
	cursor, err := m.Collection.Aggregate(context.TODO(), mongo.Pipeline{groupStage})
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
	filter := bson.M{"timeSeconds": bson.M{"$gte": after}}
	cursor, err := m.Collection.Find(context.TODO(), filter)
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

func (m *MongoStore) UserSignup(User models.User) error {
	m.ConnectToDatabaseUsers()

	wrapper := models.User{}

	wrapper = User

	var a, b interface{}
	res := m.Collection.FindOne(context.TODO(), bson.M{"email": wrapper.Email}).Decode(&a)
	res1 := m.Collection.FindOne(context.TODO(), bson.M{"username": wrapper.Username}).Decode(&b)
	if res == mongo.ErrNoDocuments && res1 == mongo.ErrNoDocuments {
		_, err := m.Collection.InsertOne(context.TODO(), User)
		if err != nil {
			zap.S().Errorf("Error while storing User details: %v", err)
			return nil
		}
		zap.S().Infof("User signed up successfully")
		return nil
	}

	zap.S().Errorf("Bad credentials")
	return nil
}

func (m *MongoStore) SubscribeToBlog(username string, BlogID int) error {
	m.ConnectToDatabaseUsers()

	filter := bson.M{
		"username":        username,
		"subscribedBlogs": bson.M{"$not": bson.M{"$all": bson.A{BlogID}}},
	}
	update := bson.M{"$push": bson.M{"subscribedBlogs": BlogID}}

	var result interface{}
	err := m.Collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&result)
	if err == mongo.ErrNoDocuments {
		zap.S().Errorf("Error while adding %v blog to subscribed list of %v user due to %v", BlogID, username, err)
	}

	zap.S().Infof("Successfully added %v blogID to %v user", BlogID, username)
	return err
}

func (m *MongoStore) UnsubscribeFromBlog(username string, BlogID int) error {
	m.ConnectToDatabaseUsers()

	filter := bson.M{"username": username}
	update := bson.M{"$pull": bson.M{"subscribedBlogs": BlogID}}

	var result interface{}
	err := m.Collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&result)
	if err == mongo.ErrNoDocuments {
		zap.S().Errorf("Error while removing %v blog from subscribed list of %v user due to %v", BlogID, username, err)
	}

	zap.S().Infof("Successfully removed %v blogID from %v user", BlogID, username)
	return err
}

func (m *MongoStore) QueryRecentActionsForUser(username string, after int64, limit int64) ([]models.RecentAction, error) {
	m.ConnectToDatabaseUsers()

	filter := bson.M{"username": username}
	result := models.User{}
	if err := m.Collection.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		zap.S().Errorf("Error while finding the username to query recent actions for user")
		return nil, nil
	}

	m.ConnectToDatabaseRecentAction()

	var final []models.RecentAction
	var err error
	for i := 0; i < len(result.SubscribedBlogs); i++ {
		filter = bson.M{
			"blogEntry.id": result.SubscribedBlogs[i],
			"timeSeconds":  bson.M{"$gte": after},
		}
		opts := options.Find().SetLimit(limit)
		cursor, err := m.Collection.Find(context.TODO(), filter, opts)
		if err != nil {
			zap.S().Errorf("Error while finding blogs which user has subscribed")
			return nil, nil
		}
		var res []models.RecentAction

		if err = cursor.All(context.TODO(), &res); err != nil {
			zap.S().Errorf("Error while decoding cursor %vth time: %v", i, err)
		}

		for i := 0; i < len(res); i++ {
			final = append(final, res[i])
		}
	}
	return final, err
}
