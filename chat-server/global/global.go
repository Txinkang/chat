package global

import (
	"chat-server/config"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"
	
)

var (
	CHAT_CONFIG  config.AppConfig
	CHAT_MYSQL   *gorm.DB
	CHAT_REDIS   *redis.Client
	CHAT_MONGO   *mongo.Client
	CHAT_MONGODB *mongo.Database
	CHAT_ES      *elasticsearch.Client
	CHAT_ROUTERS *gin.Engine
	CHAT_LOG     *slog.Logger
)
