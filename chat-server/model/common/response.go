package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// ResponseCode 定义响应码和消息的组合
type ResponseCode struct {
	Code int
	Msg  string
}

// 预定义响应状态
var (
	SUCCESS          = ResponseCode{Code: 200, Msg: "操作成功"}
	ERROR            = ResponseCode{Code: 500, Msg: "操作失败"}
	INVALID_PARAMS   = ResponseCode{Code: 400, Msg: "请求参数错误"}
	UNAUTHORIZED     = ResponseCode{Code: 401, Msg: "未授权"}
	FORBIDDEN        = ResponseCode{Code: 403, Msg: "访问被禁止"}
	NOT_FOUND        = ResponseCode{Code: 404, Msg: "资源不存在"}
	SERVER_ERROR     = ResponseCode{Code: 500, Msg: "服务器内部错误"}
	TOO_MANY_REQUEST = ResponseCode{Code: 429, Msg: "请求过于频繁"}
)

// Result 返回统一格式的响应
// 如果不提供data参数，则默认为空map
func Result(c *gin.Context, respCode ResponseCode, data ...interface{}) {
	// 设置默认的空数据
	var responseData interface{} = map[string]interface{}{}

	// 如果提供了data参数，则使用它
	if len(data) > 0 {
		responseData = data[0]
	}

	c.JSON(http.StatusOK, Response{
		Code: respCode.Code,
		Msg:  respCode.Msg,
		Data: responseData,
	})
}
