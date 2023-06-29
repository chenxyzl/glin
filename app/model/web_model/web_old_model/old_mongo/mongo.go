package old_mongo

import (
	"context"
	"github.com/chenxyzl/glin/slog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"laiya/config"
)

var mongoClient *mongo.Client

func GetCol(colName string) *mongo.Collection {
	if mongoClient == nil {
		slog.Panicf("mongo client not connect")
	}
	return mongoClient.Database(config.Get().AConfig.OldMongo.DbName).Collection(colName)
}

func Connect() *mongo.Client {
	//
	url := config.Get().AConfig.OldMongo.Url
	slog.Infof("mongo starting...,url:%v", url)
	// Rest of the code will go here
	// Set client options 设置连接参数
	//clientOptions := options.Client().ApplyURI("mongodb://root:xxxxxxxxxxxxxxxxxxxxxxx@172.20.52.158:41134/?connect=direct;authSource=admin")
	clientOptions := options.Client().ApplyURI(url).SetMaxPoolSize(10)

	// Connect to MongoDB 连接数据库
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		slog.Panic(err)
	}

	// Check the connection 测试连接
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		slog.Panic(err)
	}

	//
	mongoClient = client
	slog.Infof("mongo start success!,url:%v", url)
	return mongoClient
}

func Close() {
	slog.Infof("mongo stop...")
	if mongoClient != nil {
		mongoClient.Disconnect(context.TODO())
	}
	slog.Infof("mongo stop success!")
}

func Transaction(f func() (interface{}, error)) (interface{}, error) {
	if mongoClient == nil {
		slog.Panicf("mongo client not connect")
	}
	//
	opts := options.Transaction().SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
	// Start a transaction
	session, err := mongoClient.StartSession()
	if err != nil {
		slog.Error("Error starting a session:", err)
		return nil, err
	}
	defer session.EndSession(context.Background())
	// Define the transaction callback function
	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		return f()
	}
	// Execute the transaction
	return session.WithTransaction(context.Background(), callback, opts)
}
