package v1

import (
	"chat-server/model/common"
	"chat-server/model/request"
	"fmt"
	"github.com/gin-gonic/gin"
)

type UserApi struct{}

func (a *UserApi) Test(c *gin.Context) {
	var xxx = []int{1, 2, 3}
	var yyy = map[string]interface{}{
		"name": "张三",
		"age":  18,
	}
	var zzz = map[string]interface{}{
		"xxx": xxx,
		"yyy": yyy,
	}
	common.Result(c, common.SUCCESS, zzz)
}

// Register 用户注册
func (a *UserApi) Register(c *gin.Context) {
	var req request.RegisterRequest

	// 校验参数
	if err := c.ShouldBindJSON(&req); err != nil {
		// 打印具体错误信息以便调试
		fmt.Println("绑定错误:", err)
		common.Result(c, common.INVALID_PARAMS)
		return
	}

	//处理业务
	if err := userService.RegisterUser(req.Username, req.Password, req.Email); err != nil {
		common.Result(c, common.ERROR)
		return
	}

}

// Login 用户登录
