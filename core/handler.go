// 接口处理器，封装了参数绑定和错误处理

package core

import (
	"github.com/gin-gonic/gin"
)

func Handler[T any](fn func(*gin.Context, T) (data any, err error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 参数绑定
		var req T

		if err := c.ShouldBind(&req); err != nil {
			ErrorWithCode(c, 400, err.Error())
			c.Abort()
			return
		}

		data, err := fn(c, req)
		if err != nil {
			Error(c, err.Error())
			c.Abort()
			return
		}
		Success(c, data)
	}
}
