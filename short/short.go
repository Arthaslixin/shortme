package short

import (
	"context"
	"doodod.com/doodod/shortme/base"
	"doodod.com/doodod/shortme/conf"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"sync"
	"time"
)

var sequencesMu sync.RWMutex

type shorter struct{}

type ShoutUrlSchema struct {
	Id         primitive.ObjectID `bson:"_id"`
	LongURL    string             `bson:"long_url"`
	ShortURL   string             `bson:"short_url"`
	CreateTime int                `bson:"create_time"`
}

type SequenceSchema struct {
	Id primitive.ObjectID `bson:"_id"`
	TP string             `bson:"tp"`
	ID int                `bson:"id"`
}

func (shorter *shorter) Connect() (ctx context.Context, client *mongo.Client) {
	to := 5 * time.Second
	opts := options.ClientOptions{ConnectTimeout: &to}
	opts.SetDirect(true)
	env := os.Getenv("ENV_TYPE")
	var addr string
	// 你的MongoUri
	addr = "mongodb://"
	opts.ApplyURI(addr)
	client, _ = mongo.Connect(context.TODO(), &opts)

	ctx = context.TODO()
	_ = client.Connect(ctx)
	return ctx, client
}

func (shorter *shorter) NextSequence() (sequence uint64, err error) {
	sequencesMu.RLock()
	defer sequencesMu.RUnlock()
	ctx, client := shorter.Connect()
	collection := client.Database("db").Collection("auto_incr_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var result SequenceSchema

	err = collection.FindOneAndUpdate(ctx, bson.M{"tp": "short_url_sequence"}, bson.M{"$inc": bson.M{"id": 1}}).Decode(&result)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	sequence = uint64(result.ID)
	return sequence, nil
}

func (shorter *shorter) Expand(shortURL string) (longURL string, err error) {
	ctx, client := shorter.Connect()
	collection := client.Database("db").Collection("short_url")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var result ShoutUrlSchema

	err = collection.FindOne(ctx, bson.D{{"short_url", shortURL}}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	longURL = result.LongURL

	return longURL, nil
}

func (shorter *shorter) Short(longURL string) (shortURL string, err error) {
	for {
		var seq uint64
		seq, err = shorter.NextSequence()
		if err != nil {
			log.Printf("get next sequence error. %v", err)
			return "", errors.New("get next sequence error")
		}

		shortURL = base.Int2String(seq)
		if _, exists := conf.Conf.Common.BlackShortURLsMap[shortURL]; exists {
			continue
		} else {
			break
		}
	}

	ctx, client := shorter.Connect()
	collection := client.Database("db").Collection("short_url")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//var result ShoutUrlSchema

	_, err = collection.InsertOne(ctx,
		bson.D{
			{"long_url", longURL},
			{"short_url", shortURL},
			{"create_time", time.Now().Unix()}})

	if err != nil {
		log.Fatal(err)
	}
	return shortURL, nil
}

var Shorter shorter
