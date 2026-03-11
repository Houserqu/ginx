package core

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

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

		// 若 T 为指针类型（如 *FooParams），需先初始化防止 nil 解引用
		v := reflect.ValueOf(&req).Elem()
		isPtr := v.Kind() == reflect.Ptr
		if isPtr {
			v.Set(reflect.New(v.Type().Elem()))
		}

		// 绑定目标：指针类型直接用 req（本身已是指针），值类型取地址
		var bindTarget any
		if isPtr {
			bindTarget = req
		} else {
			bindTarget = &req
		}

		if c.Request.Method == "GET" {
			if err := c.ShouldBindQuery(bindTarget); err != nil {
				Error(c, err.Error())
				c.Abort()
				return
			}
		} else {
			if err := c.ShouldBindJSON(bindTarget); err != nil {
				Error(c, err.Error())
				c.Abort()
				return
			}
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
