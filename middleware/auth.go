// 鉴权中间件

package middleware

import (
	"ginx/core"

	"github.com/gin-gonic/gin"
)

// 检查用户是否登录
func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if core.GetUserId(c) == 0 {
			core.ErrorWithCode(c, 401, "未登录")
			c.Abort()
			return
		}

		// 调用下一个中间件，或者控制器处理函数，具体得看注册了多少个中间件。
		c.Next()
	}
}
