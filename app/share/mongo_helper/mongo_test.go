package mongo_helper

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestMongo(t *testing.T) {
	// Rest of the code will go here
	// Set client options 设置连接参数
	//clientOptions := options.Client().ApplyURI("mongodb://root:xxxxxxxxxxxxxxxxxxxxxxx@172.20.52.158:41134/?connect=direct;authSource=admin")
	url := "mongodb+srv://ichenzhl:<password>@cluster0.feqwf3z.mongodb.net/?retryWrites=true&w=majority"
	clientOptions := options.Client().ApplyURI(url).SetMaxPoolSize(10)

	// Connect to MongoDB 连接数据库
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		t.Error(err)
	}
	defer client.Disconnect(context.TODO())

	// Check the connection 测试连接
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		t.Error(err)
	}

	fmt.Println("Connected to MongoDB!")

	databases, err := client.ListDatabases(context.TODO(), bson.M{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(databases.TotalSize / 1024 / 1024 / 1024)

	if err != nil {
		t.Error(err)
	}
	fmt.Println("Connection to MongoDB closed.")

	filter := bson.D{{"_id", "aaa"}}
	type SessionModel struct {
		Account string `bson:"_id,omitempty"`
		Uid     uint64 `bson:"uid"`
	}
	var res SessionModel
	err = client.Database("test").Collection("account").FindOne(context.TODO(), filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			res.Account = "aaa"
			res.Uid = 1
			_, err = client.Database("test").Collection("account").InsertOne(context.TODO(), res)
			if err != nil {
				t.Error(err)
			}
			err = client.Database("test").Collection("account").FindOne(context.TODO(), filter).Decode(&res)
			if err != nil {
				t.Error(err)
			}
		} else {
			t.Error(err)
		}

	}
	fmt.Println(res)
}
