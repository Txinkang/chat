package v1

import (
	"chat-server/model/common"
	"chat-server/model/request"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type UserApi struct{}

// Test godoc
// @Summary      测试接口
// @Description  测试接口，返回一些示例数据
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200  {object}  common.Response
// @Router       /user/test [get]
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

// Register godoc
// @Summary      用户注册
// @Description  创建新用户账号
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request  body      request.RegisterRequest  true  "用户注册信息"
// @Success      200      {object}  common.Response
// @Router       /api/v1/user/register [post]
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
	tokenPair, err := userService.RegisterUser(req.Username, req.Password, req.Email)
	if err != nil {
		var serviceErr common.ServiceErr
		if errors.As(err, &serviceErr) {
			common.Result(c, serviceErr.GetResponseCode())
		}
		return
	}

	common.Result(c, common.SUCCESS, tokenPair)
}
