package game

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient *mongo.Client
	GameDB      *mongo.Database
	SavesColl   *mongo.Collection
)

const DefaultMongoURI = "mongodb+srv://maizergomes_db_user:DTpRP0oLnLb1DhoD@magicadvanture.rjlkgbl.mongodb.net/?appName=MagicAdvanture"

func InitDB() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = DefaultMongoURI
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	MongoClient = client
	GameDB = client.Database("magicadventure")
	SavesColl = GameDB.Collection("saves")

	return client, nil
}

func SaveGame(gs *GameState) error {
	gs.LastSeen = time.Now().Unix()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": gs.ID}
	update := bson.M{"$set": gs}

	_, err := SavesColl.UpdateOne(ctx, filter, update, opts)
	return err
}

func LoadSave(slotID int) (*GameState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var gs GameState
	err := SavesColl.FindOne(ctx, bson.M{"_id": slotID}).Decode(&gs)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &gs, err
}

func GetAllSaves() ([]GameState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := SavesColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var saves []GameState
	if err = cursor.All(ctx, &saves); err != nil {
		return nil, err
	}
	return saves, nil
}

func GetOnlinePlayers(currentGS *GameState) ([]GameState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Consider players "online" if seen in the last 2 minutes
	threshold := time.Now().Unix() - 120
	filter := bson.M{
		"last_seen": bson.M{"$gt": threshold},
		"_id":       bson.M{"$ne": currentGS.ID},
	}

	cursor, err := SavesColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var players []GameState
	if err = cursor.All(ctx, &players); err != nil {
		return nil, err
	}
	return players, nil
}
