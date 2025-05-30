package initialize

import (
	"chat-server/service"
	"context"
)

func Initialize() error {
	InitConfig()
	//InitMySQL()
	InitMongo()
	InitElasticSearch()
	//InitRedis()

	service.StartMongoToEsSync(context.Background(), MongoDB, EsClient, "messages", "messages")
	return nil
}
