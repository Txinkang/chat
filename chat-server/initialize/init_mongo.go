package initialize

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"log/slog"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

func InitMongo() error {
	slog.Info("初始化mongo")
	mongoConfig := AppConfig.Mongo
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOpts := options.Client().ApplyURI(mongoConfig.URI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return fmt.Errorf("连接mongoDB失败: %v", err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return fmt.Errorf("mongoDB Ping 失败: %v", err)
	}
	MongoClient = client
	MongoDB = client.Database(mongoConfig.DBName)
	slog.Info("MongoDB连接成功")
	return nil
}
