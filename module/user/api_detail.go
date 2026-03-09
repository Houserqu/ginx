package user

import (
	"ginx/core"
	"ginx/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	core.Register(func(svr *gin.Engine) {
		svr.GET("/api/user/detail", middleware.CheckLogin(), core.Handler(Detail))
	})
}

type DetailParams struct {
}

// Detail godoc
// @Summary     detail 接口
// @Description detail 接口描述
// @Tags        user
// @Accept      json
// @Produce     json
// @Param       request query  DetailParams true "请求参数"
// @Success     200 {object} core.Response "成功响应"
// @Security    JWT
// @Router      /api/user/detail [GET]
func Detail(c *gin.Context, params *DetailParams) (data any, err error) {

	return
}
