package util

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoConn(env map[string]string) (client *mongo.Client, ctx context.Context) {
	// timeout 기반의 Context 생성
	ctx, _ = context.WithTimeout(context.Background(), time.Second*4)

	// Authetication 을 위한 Client Option 구성
	clientOptions := options.Client().ApplyURI(
		env["dbURI"]).SetAuth(
		options.Credential{
			AuthSource: "",
			Username:   env["dbUserName"],
			Password:   env["dbPassword"],
		},
	)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("MongoDB Connection Success")
	return client, ctx
}

func AllData(collection string, mongoDB *mongo.Database, ctx context.Context) string {
	var datas []bson.M
	res, err := mongoDB.Collection(collection).Find(ctx, bson.M{})
	if err = res.All(ctx, &datas); err != nil {
		fmt.Println(err)
	}

	data, _ := json.MarshalIndent(&datas, "", "	")

	return string(data)
}
