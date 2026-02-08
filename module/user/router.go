// 路由绑定

package user

import (
	"ginx/core"
	"ginx/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(g *gin.Engine) {
	group := g.Group("/api/user", middleware.CheckLogin())

	group.POST("/login", core.Handler(UserLogin))
	group.GET("/list", UserList)
}
