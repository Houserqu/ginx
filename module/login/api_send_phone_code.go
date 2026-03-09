package login

import (
	"ginx/core"
	"ginx/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	core.Register(func(svr *gin.Engine) {
		svr.POST("/api/login/send_phone_code", middleware.CheckLogin(), core.Handler(SendPhoneCode))
	})
}

type SendPhoneCodeParams struct {
	Phone string `json:"phone" binding:"required"` // 手机号
	Scene string `json:"scene" binding:"required"` // 场景，例如 "login", "register", "reset_password" 等
}

// SendPhoneCode godoc
// @Summary     send_phone_code 接口
// @Description send_phone_code 接口描述
// @Tags        login
// @Accept      json
// @Produce     json
// @Param       request body SendPhoneCodeParams true "请求参数"
// @Router      /api/login/send_phone_code [POST]
func SendPhoneCode(c *gin.Context, params *SendPhoneCodeParams) (data any, err error) {

	return
}
