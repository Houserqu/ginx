package user

import (
	"ginx/core"
	"ginx/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	core.Register(func(svr *gin.Engine) {
		svr.POST("/api/user/delete_account", middleware.CheckLogin(), core.Handler(DeleteAccount))
	})
}

type DeleteAccountParams struct {
}

// DeleteAccount godoc
// @Summary     delete_account 接口
// @Description delete_account 接口描述
// @Tags        user
// @Accept      json
// @Produce     json
// @Param       request body DeleteAccountParams true "请求参数"
// @Success     200 {object} core.Response "成功响应"
// @Router      /api/user/delete_account [POST]
func DeleteAccount(c *gin.Context, params *DeleteAccountParams) (data any, err error) {

	return
}
