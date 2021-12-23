package database

import (
	"context"
	"log"
	"time"

	"github.com/H-Richard/go-graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	return &DB{
		client: client,
	}
}

func (db *DB) Save(input *model.NewArticle) *model.Article {
	collection := db.client.Database("testdb").Collection("article")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
	return &model.Article{
		ID:         res.InsertedID.(primitive.ObjectID).Hex(),
		Name:       input.Name,
		Definition: input.Definition,
		Directory:  input.Directory,
	}
}

func (db *DB) FindByID(ID string) *model.Article {
	ObjectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		log.Fatal(err)
	}
	collection := db.client.Database("testdb").Collection("article")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res := collection.FindOne(ctx, bson.M{"_id": ObjectID})
	article := model.Article{}
	res.Decode(&article)
	return &article
}

func (db *DB) All() []*model.Article {
	collection := db.client.Database("testdb").Collection("article")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var articles []*model.Article
	for cur.Next(ctx) {
		var article *model.Article
		err := cur.Decode(&article)
		if err != nil {
			log.Fatal(err)
		}
		articles = append(articles, article)
	}
	return articles
}
