package login

import (
  "ginx/core"
	"ginx/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	core.Register(func(svr *gin.Engine) {
		svr.GET("/api/login/login_info", middleware.CheckLogin(), core.Handler(LoginInfo))
	})
}

type LoginInfoParams struct {
}

// LoginInfo godoc
// @Summary     login_info 接口
// @Description login_info 接口描述
// @Tags        login
// @Accept      json
// @Produce     json
// @Param       request body LoginInfoParams true "请求参数"
// @Router      /api/login/login_info [GET]
func LoginInfo(c *gin.Context, params *LoginInfoParams) (data any, err error) {

	return
}
