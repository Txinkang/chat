package v1

import "chat-server/service"

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	UserApi
	ChatApi
}

var (
	chatService   = service.ServiceGroupApp.ChatService
	userService   = service.ServiceGroupApp.UserService
	mongoToEsSync = service.ServiceGroupApp.MongoToEsSync
)
