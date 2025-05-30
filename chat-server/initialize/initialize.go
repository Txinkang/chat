package initialize

func Initialize() error {
	InitConfig()
	InitLogger()
	InitMySQL()
	InitMongo()
	//InitElasticSearch()
	//InitRedis()

	//service.StartMongoToEsSync(context.Background(), MongoDB, EsClient, "messages", "messages")
	return nil
}
