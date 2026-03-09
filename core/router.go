package core

import "github.com/gin-gonic/gin"

var globalRoutes []func(*gin.Engine)

// Register 注册路由，在各 api_*.go 的 init() 中调用
func Register(fn func(*gin.Engine)) {
	globalRoutes = append(globalRoutes, fn)
}

// InitRoutes 在 main.go 中调用一次，将所有已注册路由挂载到 gin.Engine
func InitRoutes(svr *gin.Engine) {
	for _, fn := range globalRoutes {
		fn(svr)
	}
}

func Handler[T any](fn func(*gin.Context, T) (data any, err error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 参数绑定
		var req T

		if err := c.ShouldBind(&req); err != nil {
			Error(c, err.Error())
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
