// 用户列表接口实现

package user

import (
	"ginx/model"
	"ginx/utils/crud"

	"github.com/gin-gonic/gin"
)

func UserList(c *gin.Context) {
	var req struct {
		Phone string `form:"phone"`
	}

	// 直接使用 curd 封装快速实现用户列表查询
	crud.List[model.User](c, &req, crud.WithPreLoads("Teams"))
}
