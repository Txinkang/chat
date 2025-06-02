package v1

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	UserApi
	ChatApi
}
