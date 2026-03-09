package login

import (
	"ginx/core"
	"ginx/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	core.Register(func(svr *gin.Engine) {
		svr.POST("/api/login/login_by_phone", middleware.CheckLogin(), core.Handler(LoginByPhone))
	})
}

type LoginByPhoneParams struct {
	Phone string `json:"phone" binding:"required" example:"18999999999"` // 手机号
}

// LoginByPhone godoc
// @Summary     login_by_phone 接口
// @Description login_by_phone 接口描述
// @Tags        login
// @Accept      json
// @Produce     json
// @Param       request body LoginByPhoneParams true "请求参数"
// @Success     200 {object} core.Response "成功响应"
// @Router      /api/login/login_by_phone [POST]
func LoginByPhone(c *gin.Context, params *LoginByPhoneParams) (data any, err error) {

	return
}
